package domain

import "time"

<<<<<<< HEAD
type Province struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NameTh string `gorm:"type:varchar(100)" json:"name_th"`
	NameEn string `gorm:"type:varchar(100)" json:"name_en"`
	GeographyId int       `gorm:"not null" json:"geography_id"`
=======
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
>>>>>>> fe689551f1395d734858e14c96b935403361a507
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

<<<<<<< HEAD

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
=======
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
>>>>>>> fe689551f1395d734858e14c96b935403361a507
