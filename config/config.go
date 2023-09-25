package config

import "github.com/pitabwire/frame"

type OcrConfig struct {
	frame.ConfigurationDefault
	FilesServiceURI string `default:"127.0.0.1:7020" envconfig:"FILES_SERVICE_URI"`

	QueueOcrSync     string `default:"mem://ocr_model_sync" envconfig:"QUEUE_OCR_SYNC"`
	QueueOcrSyncName string `default:"ocr_model_sync" envconfig:"QUEUE_OCR_SYNC_NAME"`
}
