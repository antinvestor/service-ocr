package handlers

import (
	"context"
	fapi "github.com/antinvestor/service-files-api"
	ocr "github.com/antinvestor/service-ocr-api"
	"github.com/antinvestor/service-ocr/service/business"
	"github.com/pitabwire/frame"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OCRServer struct {
	Service  *frame.Service
	FilesCli *fapi.FilesClient

	ocr.OCRServiceServer
}

func (os *OCRServer) Recognize(ctx context.Context, request *ocr.OcrRequest) (*ocr.OcrResponse, error) {

	businessOcr := business.NewOcrBusiness(ctx, os.Service, os.FilesCli)
	response, err := businessOcr.Recognize(ctx, request)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return response, nil
}

func (os *OCRServer) Status(ctx context.Context, request *ocr.StatusRequest) (*ocr.OcrResponse, error) {
	businessOcr := business.NewOcrBusiness(ctx, os.Service, os.FilesCli)
	response, err := businessOcr.CheckProgress(ctx, request)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return response, nil
}
