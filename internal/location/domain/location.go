package domain

import "time"

// Geography represents thai_geographies table
// Each geography has ID and Name
//
// gorm model to map to thai_geographies table

// We'll include CreatedAt, UpdatedAt? Not necessary. Since table has only id and name

type Geography struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name"`
}

type Province struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	NameTH      string    `gorm:"column:name_th" json:"name_th"`
	NameEN      string    `gorm:"column:name_en" json:"name_en"`
	GeographyID int       `gorm:"column:geography_id" json:"geography_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Amphure struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	NameTH     string    `gorm:"column:name_th" json:"name_th"`
	NameEN     string    `gorm:"column:name_en" json:"name_en"`
	ProvinceID int       `gorm:"column:province_id" json:"province_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Tambon struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	ZipCode   int       `gorm:"column:zip_code" json:"zip_code"`
	NameTH    string    `gorm:"column:name_th" json:"name_th"`
	NameEN    string    `gorm:"column:name_en" json:"name_en"`
	AmphureID int       `gorm:"column:amphure_id" json:"amphure_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
