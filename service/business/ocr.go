package business

import (
	"context"
	"fmt"
	"github.com/antinvestor/apis/common"
	fapi "github.com/antinvestor/service-files-api"
	ocr "github.com/antinvestor/service-ocr-api"
	"github.com/antinvestor/service-ocr/config"
	"github.com/antinvestor/service-ocr/service/models"
	"github.com/antinvestor/service-ocr/service/repository"
	"github.com/pitabwire/frame"
	"os"
)

type OCRBusiness interface {
	Recognize(ctx context.Context, request *ocr.OcrRequest) (*ocr.OcrResponse, error)
	CheckProgress(ctx context.Context, request *ocr.StatusRequest) (*ocr.OcrResponse, error)
	ToApi(ocrLogList []*models.OcrLog) *ocr.OcrResponse
}

type Recognizer interface {
	Recognize(ctx context.Context, image *os.File) (string, error)
}

func NewOcrBusiness(ctx context.Context, service *frame.Service, filesCli *fapi.FilesClient) OCRBusiness {
	return &ocrBusiness{
		service:    service,
		filesCli:   filesCli,
		ocrRepo:    repository.NewOcrRepository(service),
		recognizer: &googleCloudVision{},
	}
}

type ocrBusiness struct {
	service    *frame.Service
	filesCli   *fapi.FilesClient
	ocrRepo    repository.OcrRepository
	recognizer Recognizer
}

func (ob *ocrBusiness) ToApi(ocrLogList []*models.OcrLog) *ocr.OcrResponse {

	response := &ocr.OcrResponse{}

	filesResultList := make([]*ocr.OCRFile, 0)

	for index, ocrLog := range ocrLogList {

		if index == 1 {
			response.ReferenceId = ocrLog.ReferenceID
			response.State = common.STATE(ocrLog.State)
		} else {
			if response.State != common.STATE(ocrLog.State) {
				response.State = common.STATE(ocrLog.State)
			}
		}

		ocrFile := &ocr.OCRFile{
			FileId:     ocrLog.FileID,
			Language:   ocrLog.LanguageID,
			Text:       ocrLog.Text,
			Properties: frame.DBPropertiesToMap(ocrLog.Properties),
		}

		filesResultList = append(filesResultList, ocrFile)
	}

	response.Result = filesResultList

	return response
}

func (ob *ocrBusiness) CheckProgress(ctx context.Context, request *ocr.StatusRequest) (*ocr.OcrResponse, error) {

	ocrLogList, err := ob.ocrRepo.GetByReference(ctx, request.GetReferenceId())
	if err != nil {
		return nil, err
	}

	return ob.ToApi(ocrLogList), nil
}

func (ob *ocrBusiness) Recognize(ctx context.Context, request *ocr.OcrRequest) (*ocr.OcrResponse, error) {

	authClaims := frame.ClaimsFromContext(ctx)

	accessId := authClaims.AccessID

	ocrLogList := make([]*models.OcrLog, 0)

	for _, fileId := range request.GetFileId() {

		newOcrLog := &models.OcrLog{
			ReferenceID: request.GetReferenceId(),
			FileID:      fileId,
			AccessID:    accessId,
			LanguageID:  request.GetLanguageId(),
			State:       int32(common.STATE_ACTIVE),
			Status:      int32(common.STATUS_QUEUED),
			Properties:  frame.DBPropertiesFromMap(request.GetProperties()),
		}

		if request.GetAsync() {

			newOcrLog.GenID()

			err := ob.service.Publish(ctx, config.QueueOcrSyncName, newOcrLog)
			if err != nil {
				return nil, err
			}

		} else {
			err := ob.ocrRepo.Save(ctx, newOcrLog)
			if err != nil {
				return nil, err
			}

			newOcrLog, err = ob.recognize(ctx, newOcrLog)
			if err != nil {
				return nil, err
			}

		}
		ocrLogList = append(ocrLogList, newOcrLog)

	}

	return ob.ToApi(ocrLogList), nil
}

func (ob *ocrBusiness) recognize(ctx context.Context, ocrLog *models.OcrLog) (*models.OcrLog, error) {

	if common.STATUS_QUEUED == common.STATUS(ocrLog.Status) {
		ocrLog.Status = int32(common.STATUS_IN_PROCESS)
		err := ob.ocrRepo.Save(ctx, ocrLog)
		if err != nil {
			return nil, err
		}
	}

	file, _, err := ob.filesCli.DefaultApi.FindFileById(ctx, ocrLog.FileID).Execute()

	ocrLog.Text, err = ob.recognizer.Recognize(ctx, file)
	if err != nil {
		ocrLog.Status = int32(common.STATUS_FAILED)
		fmt.Printf(" recognize -- there was an error recognizing text : %v", err)
	} else {
		ocrLog.Status = int32(common.STATUS_SUCCESSFUL)
	}

	ocrLog.State = int32(common.STATE_INACTIVE)
	err = ob.ocrRepo.Save(ctx, ocrLog)
	if err != nil {
		return nil, err
	}

	return ocrLog, nil
}
