package repository

import (
	"context"

	"github.com/antinvestor/service-ocr/service/models"
	"github.com/pitabwire/frame/datastore/pool"
)

type OcrRepository interface {
	GetByID(ctx context.Context, id string) (*models.OcrLog, error)
	GetByReference(ctx context.Context, reference string) ([]*models.OcrLog, error)
	Save(ctx context.Context, ocrLog *models.OcrLog) error
	Delete(ctx context.Context, id string) error
}

type ocrRepository struct {
	dbPool pool.Pool
}

func (pr *ocrRepository) GetByID(ctx context.Context, id string) (*models.OcrLog, error) {
	ocrLog := &models.OcrLog{}
	err := pr.dbPool.DB(ctx, true).First(ocrLog, "id = ?", id).Error
	return ocrLog, err
}

func (pr *ocrRepository) GetByReference(ctx context.Context, reference string) ([]*models.OcrLog, error) {
	var ocrLog []*models.OcrLog
	err := pr.dbPool.DB(ctx, true).Find(&ocrLog, "reference_id = ?", reference).Error
	return ocrLog, err
}

func (pr *ocrRepository) Save(ctx context.Context, ocrLog *models.OcrLog) error {
	return pr.dbPool.DB(ctx, false).Save(ocrLog).Error
}

func (pr *ocrRepository) Delete(ctx context.Context, id string) error {
	ocrLog, err := pr.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return pr.dbPool.DB(ctx, false).Delete(ocrLog).Error
}

func NewOcrRepository(dbPool pool.Pool) OcrRepository {
	return &ocrRepository{
		dbPool: dbPool,
	}
}
