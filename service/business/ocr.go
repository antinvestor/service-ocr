package business

import (
	"context"
	"fmt"
	"os"

	commonv1 "buf.build/gen/go/antinvestor/common/protocolbuffers/go/common/v1"
	"buf.build/gen/go/antinvestor/files/connectrpc/go/files/v1/filesv1connect"
	filesv1 "buf.build/gen/go/antinvestor/files/protocolbuffers/go/files/v1"
	ocrv1 "buf.build/gen/go/antinvestor/ocr/protocolbuffers/go/ocr/v1"
	"connectrpc.com/connect"
	"github.com/antinvestor/service-ocr/config"
	"github.com/antinvestor/service-ocr/service/models"
	"github.com/antinvestor/service-ocr/service/repository"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/queue"
	"github.com/pitabwire/frame/security"
	"google.golang.org/protobuf/types/known/structpb"
)

type OCRBusiness interface {
	Recognize(ctx context.Context, request *ocrv1.RecognizeRequest) (*ocrv1.RecognizeResponse, error)
	CheckProgress(ctx context.Context, request *commonv1.StatusRequest) (*ocrv1.StatusResponse, error)
	ToAPI(ocrLogList []*models.OcrLog) *ocrv1.RecognizeResponse
}

type Recognizer interface {
	Recognize(ctx context.Context, image *os.File) (string, error)
}

func NewOcrBusiness(_ context.Context, svc *frame.Service, filesCli filesv1connect.FilesServiceClient,
	ocrRepo repository.OcrRepository, queueMan queue.Manager) OCRBusiness {
	return &ocrBusiness{
		svc:        svc,
		filesCli:   filesCli,
		ocrRepo:    ocrRepo,
		queueMan:   queueMan,
		recognizer: &googleCloudVision{},
	}
}

type ocrBusiness struct {
	svc        *frame.Service
	filesCli   filesv1connect.FilesServiceClient
	ocrRepo    repository.OcrRepository
	queueMan   queue.Manager
	recognizer Recognizer
}

func jsonMapToStruct(m map[string]interface{}) *structpb.Struct {
	if m == nil {
		return nil
	}
	s, err := structpb.NewStruct(m)
	if err != nil {
		return nil
	}
	return s
}

func structToJSONMap(s *structpb.Struct) map[string]interface{} {
	if s == nil {
		return nil
	}
	return s.AsMap()
}

func (ob *ocrBusiness) ToAPI(ocrLogList []*models.OcrLog) *ocrv1.RecognizeResponse {
	response := &ocrv1.RecognizeResponse{}

	filesResultList := make([]*ocrv1.OCRFile, 0)

	for _, ocrLog := range ocrLogList {
		response.ReferenceId = ocrLog.ReferenceID

		ocrFile := &ocrv1.OCRFile{
			FileId:     ocrLog.FileID,
			Language:   ocrLog.LanguageID,
			Text:       ocrLog.Text,
			Status:     commonv1.STATUS(ocrLog.Status),
			Properties: jsonMapToStruct(ocrLog.Properties),
		}

		filesResultList = append(filesResultList, ocrFile)
	}

	response.Result = filesResultList

	return response
}

func (ob *ocrBusiness) CheckProgress(ctx context.Context, request *commonv1.StatusRequest) (*ocrv1.StatusResponse, error) {
	ocrLogList, err := ob.ocrRepo.GetByReference(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	return &ocrv1.StatusResponse{
		Data: ob.ToAPI(ocrLogList),
	}, nil
}

func (ob *ocrBusiness) Recognize(ctx context.Context, request *ocrv1.RecognizeRequest) (*ocrv1.RecognizeResponse, error) {
	authClaims := security.ClaimsFromContext(ctx)

	ocrConfig := ob.svc.Config().(*config.OcrConfig)

	accessID := ""
	if authClaims != nil {
		accessID = authClaims.AccessID
	}

	ocrLogList := make([]*models.OcrLog, 0)

	for _, fileID := range request.GetFileId() {
		properties := structToJSONMap(request.GetProperties())

		newOcrLog := &models.OcrLog{
			ReferenceID: request.GetReferenceId(),
			FileID:      fileID,
			AccessID:    accessID,
			LanguageID:  request.GetLanguageId(),
			State:       int32(commonv1.STATE_ACTIVE),
			Status:      int32(commonv1.STATUS_QUEUED),
			Properties:  properties,
		}
		if request.GetAsync() {
			newOcrLog.GenID(ctx)

			err := ob.queueMan.Publish(ctx, ocrConfig.QueueOcrSyncName, newOcrLog)
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

	return ob.ToAPI(ocrLogList), nil
}

func (ob *ocrBusiness) recognize(ctx context.Context, ocrLog *models.OcrLog) (*models.OcrLog, error) {
	if commonv1.STATUS_QUEUED == commonv1.STATUS(ocrLog.Status) {
		ocrLog.Status = int32(commonv1.STATUS_IN_PROCESS)
		err := ob.ocrRepo.Save(ctx, ocrLog)
		if err != nil {
			return nil, err
		}
	}

	// Fetch file content from Files service
	resp, err := ob.filesCli.GetContent(ctx, connect.NewRequest(&filesv1.GetContentRequest{
		MediaId: ocrLog.FileID,
	}))
	if err != nil {
		return nil, err
	}

	// Write content to temp file for recognizer
	tmpFile, err := os.CreateTemp("", "ocr-*.bin")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(resp.Msg.GetContent())
	if err != nil {
		tmpFile.Close()
		return nil, err
	}
	tmpFile.Close()

	// Reopen for reading
	readFile, err := os.Open(tmpFile.Name())
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	ocrLog.Text, err = ob.recognizer.Recognize(ctx, readFile)
	if err != nil {
		ocrLog.Status = int32(commonv1.STATUS_FAILED)
		fmt.Printf(" recognize -- there was an error recognizing text : %v", err)
	} else {
		ocrLog.Status = int32(commonv1.STATUS_SUCCESSFUL)
	}

	ocrLog.State = int32(commonv1.STATE_INACTIVE)
	err = ob.ocrRepo.Save(ctx, ocrLog)
	if err != nil {
		return nil, err
	}

	return ocrLog, nil
}
