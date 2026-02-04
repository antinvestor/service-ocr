package handlers

import (
	"context"

	commonv1 "buf.build/gen/go/antinvestor/common/protocolbuffers/go/common/v1"
	"buf.build/gen/go/antinvestor/ocr/connectrpc/go/ocr/v1/ocrv1connect"
	ocrv1 "buf.build/gen/go/antinvestor/ocr/protocolbuffers/go/ocr/v1"
	"connectrpc.com/connect"
	"github.com/antinvestor/service-ocr/service/business"
)

type OCRServer struct {
	ocrv1connect.UnimplementedOCRServiceHandler
	ocrBusiness business.OCRBusiness
}

func NewOCRServer(ocrBusiness business.OCRBusiness) *OCRServer {
	return &OCRServer{
		ocrBusiness: ocrBusiness,
	}
}

func (os *OCRServer) Recognize(ctx context.Context, request *connect.Request[ocrv1.RecognizeRequest]) (*connect.Response[ocrv1.RecognizeResponse], error) {
	response, err := os.ocrBusiness.Recognize(ctx, request.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(response), nil
}

func (os *OCRServer) Status(ctx context.Context, request *connect.Request[commonv1.StatusRequest]) (*connect.Response[ocrv1.StatusResponse], error) {
	response, err := os.ocrBusiness.CheckProgress(ctx, request.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(response), nil
}
