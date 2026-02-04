package config

import "github.com/pitabwire/frame/config"

type OcrConfig struct {
	config.ConfigurationDefault
	FilesServiceURI string `envDefault:"127.0.0.1:7020" env:"FILES_SERVICE_URI"`

	QueueOcrSync     string `envDefault:"mem://ocr_model_sync" env:"QUEUE_OCR_SYNC"`
	QueueOcrSyncName string `envDefault:"ocr_model_sync" env:"QUEUE_OCR_SYNC_NAME"`
}
