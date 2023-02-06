package checkin

import (
	"time"

	"gorm.io/gorm"
)

type CheckinDetail struct {
	Id         int `json:"id" gorm:"column:id; PRIMARY_KEY"`
	CheckinsId int `json:"checkins_id"`
	ItemId     int `json:"item_id"`
	Qty        int `json:"qty"`
}

type Checkin struct {
	Id         int       `json:"id" gorm:"column:id; PRIMARY_KEY"`
	Code       string    `json:"code"`
	Total      int       `json:"total"`
	SupplierId int       `json:"supplier_id"`
	GudangId   int       `json:"gudang_id"`
	Tanggal    time.Time `json:"taggal"`
	Keterangan string    `json:"keterangan"`
}
type CheckinResponse struct {
	Id         int       `json:"id,omitempty" gorm:"column:id; PRIMARY_KEY"`
	Code       string    `json:"code,omitempty"`
	Total      int       `json:"total,omitempty"`
	Supplier   string    `json:"supplier,omitempty"`
	SupplierId int       `json:"supplier_id,omitempty"`
	GudangId   int       `json:"gudang_id,omitempty"`
	Gudang     string    `json:"gudang,omitempty"`
	Tanggal    time.Time `json:"tanggal"`
	Keterangan string    `json:"keterangan"`
}

type CheckinDetailResponse struct {
	Id     int    `json:"checkins_id" gorm:"column:id; PRIMARY_KEY"`
	ItemId int    `json:"id"`
	Nama   string `json:"nama"`
	Qty    int    `json:"qty"`
}

type CheckinRepositroy interface {
	Create(Data Checkin, Detail []CheckinDetail) error
	FindAll() ([]CheckinResponse, error)
	FindOne(Id string) (CheckinResponse, error)
	FindOneDetail(Id string) ([]CheckinDetailResponse, error)
	Delete(Id string) error
	Update(Id string, Data Checkin, Detail []CheckinDetail) error
}

type checkinRepositroy struct {
	db *gorm.DB
}

func NewCheckinRepository(db *gorm.DB) *checkinRepositroy {
	return &checkinRepositroy{db: db}
}

func (r *checkinRepositroy) Create(Data Checkin, Detail []CheckinDetail) error {

	data := Data

	tx := r.db.Begin()

	err := tx.Table("checkins").Create(&data).Error

	if err != nil {
		tx.Rollback()
	}

	var detail []CheckinDetail

	for _, item := range Detail {
		key := CheckinDetail{
			CheckinsId: data.Id,
			ItemId:     item.ItemId,
			Qty:        item.Qty,
		}

		detail = append(detail, key)
	}

	err = tx.Table("checkins_detail").Create(&detail).Error

	if err != nil {
		tx.Rollback()
	}

	err = tx.Commit().Error

	return err
}

func (r *checkinRepositroy) FindAll() ([]CheckinResponse, error) {
	var data []CheckinResponse

	err := r.db.Table("checkins").Select("checkins.id, checkins.code, checkins.total, tbm_suppliers.nama as supplier, tbm_gudang.nama as gudang, checkins.tanggal, checkins.keterangan").
		Joins("left join tbm_gudang on tbm_gudang.id = checkins.gudang_id").
		Joins("left join tbm_suppliers on tbm_suppliers.id = checkins.supplier_id").Order("checkins.id DESC").
		Find(&data).Error

	return data, err
}

func (r *checkinRepositroy) Delete(Id string) error {
	err := r.db.Exec("DELETE FROM checkins WHERE id = ?", Id).Error

	return err
}

func (r *checkinRepositroy) FindOne(Id string) (CheckinResponse, error) {
	var data CheckinResponse

	err := r.db.Table("checkins").Select("checkins.id, checkins.code, checkins.total, checkins.gudang_id, checkins.supplier_id, tbm_suppliers.nama as supplier, tbm_gudang.nama as gudang, checkins.tanggal, checkins.keterangan").
		Joins("left join tbm_gudang on tbm_gudang.id = checkins.gudang_id").
		Joins("left join tbm_suppliers on tbm_suppliers.id = checkins.supplier_id").Where("checkins.id = ?", Id).
		Find(&data).Error

	return data, err
}

func (r *checkinRepositroy) FindOneDetail(Id string) ([]CheckinDetailResponse, error) {
	var data []CheckinDetailResponse

	err := r.db.Table("checkins_detail").Select("checkins_detail.id, ti.id as item_id, checkins_detail.qty, ti.nama").Joins("left join tbm_items as ti on ti.id = checkins_detail.item_id").Where("checkins_detail.checkins_id = ?", Id).Find(&data).Error

	return data, err
}

func (r *checkinRepositroy) Update(Id string, Data Checkin, Detail []CheckinDetail) error {
	data := Data

	tx := r.db.Begin()

	err := tx.Table("checkins").Where("id = ?", Id).Updates(&data).Error

	if err != nil {
		tx.Rollback()
	}

	err = tx.Exec("DELETE FROM checkins_detail WHERE checkins_id = ?", Id).Error

	if err != nil {
		tx.Rollback()
	}

	err = tx.Table("checkins_detail").Create(&Detail).Error

	if err != nil {
		tx.Rollback()
	}

	err = tx.Commit().Error

	return err

}