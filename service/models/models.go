package models

import (
	"github.com/pitabwire/frame"
	"gorm.io/datatypes"
)

type OcrLog struct {
	frame.BaseModel
	FileID      string `gorm:"type:varchar(50)"`
	ReferenceID string `gorm:"type:varchar(50)"`
	AccessID    string `gorm:"type:varchar(50)"`
	LanguageID  string `gorm:"type:varchar(20)"`
	State       int32
	Status      int32
	Properties  datatypes.JSONMap
	Text        string
}
