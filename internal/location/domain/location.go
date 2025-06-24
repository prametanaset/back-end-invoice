package domain

import "time"

type Province struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NameTh string `gorm:"type:varchar(100)" json:"name_th"`
	NameEn string `gorm:"type:varchar(100)" json:"name_en"`
	GeographyId int       `gorm:"not null" json:"geography_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}


type District struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NameTh string 			`gorm:"type:varchar(100)" json:"name_th"`
	NameEn string 			`gorm:"type:varchar(100)" json:"name_en"`
	ProvinceId  int       `gorm:"not null" json:"province_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}


type SubDistrict struct{
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NameTh 		string 			`gorm:"type:varchar(100)" json:"name_th"`
	NameEn 		string 			`gorm:"type:varchar(100)" json:"name_en"`
	DistrictId  int       `gorm:"not null" json:"district_id"`
	ZipCode     int       `gorm:"not null" json:"zip_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
