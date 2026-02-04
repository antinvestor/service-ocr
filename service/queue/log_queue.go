package queue

import (
	"context"
	"encoding/json"

	"github.com/antinvestor/service-ocr/service/models"
	"github.com/antinvestor/service-ocr/service/repository"
)

type OCRQueueHandler struct {
	repo repository.OcrRepository
}

func (oq *OCRQueueHandler) Handle(ctx context.Context, _ map[string]string, payload []byte) error {
	ocrLog := &models.OcrLog{}
	err := json.Unmarshal(payload, ocrLog)
	if err != nil {
		return err
	}

	return oq.repo.Save(ctx, ocrLog)
}

func NewOCRQueueHandler(ocrRepo repository.OcrRepository) *OCRQueueHandler {
	return &OCRQueueHandler{repo: ocrRepo}
}
