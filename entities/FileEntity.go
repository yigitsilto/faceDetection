package entities

import "time"

type FileEntity struct {
	Id        uint                  `gorm:"primarykey" json:"id"`
	Path      string                `json:"path"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	Values    []FileFeedValueEntity `gorm:"foreignKey:FileId" json:"values"`
}

type FileFeedValueEntity struct {
	Id        uint      `gorm:"primarykey" json:"id"`
	FileId    uint      `json:"file_id"` // Foreign key
	Value     string    `json:"value"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
