package queue

import (
	"context"
	"encoding/json"
	"github.com/antinvestor/service-ocr/service/models"
	"github.com/antinvestor/service-ocr/service/repository"
	"github.com/pitabwire/frame"
)

type OCRQueueHandler struct {
	service *frame.Service
	repo    repository.OcrRepository
}

func (oq *OCRQueueHandler) Handle(ctx context.Context, payload []byte) error {

	ocrLog := &models.OcrLog{}
	err := json.Unmarshal(payload, ocrLog)
	if err != nil {
		return err
	}

	return oq.repo.Save(ctx, ocrLog)

}

func NewOCRQueueHandler(service *frame.Service) *OCRQueueHandler {
	ocrRepo := repository.NewOcrRepository(service)
	return &OCRQueueHandler{service, ocrRepo}
}
