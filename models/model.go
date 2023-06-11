package models

import "time"

type Brand struct {
	ID        uint32 `gorm:"primary_key"`
	BrandCode uint16
	Name      string `gorm:"type:varchar(255)"`
}

type Model struct {
	ID      uint32 `gorm:"primary_key"`
	BrandId uint32
	Brand   Brand `gorm:"foreignKey:BrandId"`
	Name    string
}

type CascoValue struct {
	ID            uint32 `gorm:"primary_key"`
	ModelId       uint32
	Model         Model `gorm:"foreignKey:ModelId"`
	Casco         float64
	ModelYear     uint16
	LastUpdatedAt time.Time
}
