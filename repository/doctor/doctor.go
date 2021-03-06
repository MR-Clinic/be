package doctor

import (
	"be/entities"
	"be/utils"
	"errors"

	"github.com/labstack/gommon/log"
	"github.com/lithammer/shortuuid"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Create(req entities.Doctor) (entities.Doctor, error) {

	// check username

	type userNameCheck struct {
		UserName string
	}

	var checkUserName = r.db.Raw("? union all ? ", r.db.Model(&entities.Patient{}).Select("user_name").Where("user_name = ?", req.UserName), r.db.Model(&entities.Doctor{}).Select("user_name").Where("user_name = ?", req.UserName)).Scan(&userNameCheck{})

	if checkUserName.RowsAffected != 0 {
		return entities.Doctor{}, errors.New("user name is already exist")
	}

	// check email

	var checkEmail = r.db.Model(&entities.Doctor{}).Where("email = ?", req.Email).Select("user_name as UserName").Scan(&userNameCheck{})

	if checkEmail.RowsAffected != 0 {
		return entities.Doctor{}, errors.New("email is already exist")
	}

	var uid string
	for {
		uid = shortuuid.New()
		var find = entities.Doctor{}
		var res = r.db.Model(&entities.Doctor{}).Where("doctor_uid = ?", uid).Find(&find)
		if res.RowsAffected == 0 {
			break
		}
	}
	req.Doctor_uid = uid
	req.Doctor_uid_ref = uid
	var err error

	req.Password, err = utils.HashPassword(req.Password)
	if err != nil {
		log.Warn(err)
		return entities.Doctor{}, errors.New("error in hash password")
	}

	req.Type = "doctor"
	if res := r.db.Model(&entities.Doctor{}).Create(&req); res.Error != nil {
		log.Warn(err)
		return entities.Doctor{}, res.Error
	}

	var reqDoctor = req
	// create admin

	req.Type = "admin"
	req.UserName = "admin" + req.UserName
	for {
		uid = shortuuid.New()
		var find = entities.Doctor{}
		var res = r.db.Model(&entities.Doctor{}).Where("doctor_uid = ?", uid).Find(&find)
		if res.RowsAffected == 0 {
			break
		}
	}
	req.Doctor_uid = uid
	if res := r.db.Model(&entities.Doctor{}).Create(&req); res.Error != nil {
		log.Warn(err)
		return entities.Doctor{}, res.Error
	}

	return reqDoctor, nil
}

func (r *Repo) Update(doctor_uid string, req entities.Doctor) (entities.Doctor, error) {

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return entities.Doctor{}, err
	}

	var resInit entities.Doctor

	// check username

	type userNameCheck struct {
		UserName string
	}

	var checkUserName = tx.Raw("? union all ? ", r.db.Model(&entities.Patient{}).Select("user_name").Where("user_name = ?", req.UserName), r.db.Model(&entities.Doctor{}).Select("user_name").Where("user_name = ?", req.UserName)).Scan(&userNameCheck{})

	if checkUserName.RowsAffected != 0  && req.UserName != ""{
		log.Warn(checkUserName.Error)
		tx.Rollback()
		return entities.Doctor{}, errors.New("user name is already exist")
	}

	// check email

	var checkEmail = r.db.Model(&entities.Doctor{}).Where("email = ?", req.Email).Select("user_name as UserName").Scan(&userNameCheck{})

	if checkEmail.RowsAffected != 0 && req.Email != ""{
		log.Warn(checkEmail.Error)
		tx.Rollback()
		return entities.Doctor{}, errors.New("email is already exist")
	}

	// hash password

	if req.Password != "" {
		password, err := utils.HashPassword(req.Password)
		req.Password = password
		if err != nil {
			log.Warn(err)
			return entities.Doctor{}, errors.New("error in hash password")
		}
	}

	if res := tx.Model(&entities.Doctor{}).Where("doctor_uid = ?", doctor_uid).Updates(entities.Doctor{
		UserName: req.UserName,
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Image:    req.Image,
		Address:  req.Address,
		Status:   req.Status,
		OpenDay:  req.OpenDay,
		CloseDay: req.CloseDay,
		Capacity: req.Capacity}); res.Error != nil || res.RowsAffected == 0 {
		switch {
		case res.Error == nil:
			tx.Rollback()
			return entities.Doctor{}, gorm.ErrRecordNotFound
		default:
			tx.Rollback()
			return entities.Doctor{}, res.Error
		}
	}

	return resInit, tx.Commit().Error
}

func (r *Repo) Delete(doctor_uid string) (entities.Doctor, error) {
	var resInit entities.Doctor

	if res := r.db.Model(&entities.Doctor{}).Where("doctor_uid = ?", doctor_uid).Delete(&resInit); res.Error != nil || res.RowsAffected == 0 {
		log.Warn(res.Error)
		return entities.Doctor{}, gorm.ErrRecordNotFound
	}

	return resInit, nil
}

func (r *Repo) GetProfile(doctor_uid, userName, email string) (ProfileResp, error) {

	switch {
	case doctor_uid != "":
		doctor_uid = "doctor_uid = '" + doctor_uid + "'"
	case userName != "":
		doctor_uid = "user_name = '" + userName + "'"
	case email != "":
		doctor_uid = "email = '" + email + "'"
	}

	var profileResp ProfileResp

	var query = "doctor_uid as Doctor_uid, user_name as UserName, email as Email, name as Name, image as Image, address as Address, status as Status, open_day as OpenDay, close_day as CloseDay, capacity as Capacity, doctor_uid_ref as Doctor_uid_ref "

	if res := r.db.Model(&entities.Doctor{}).Where(doctor_uid).Select(query).Find(&profileResp); res.Error != nil || res.RowsAffected == 0 {
		log.Warn(res.Error)
		return ProfileResp{}, gorm.ErrRecordNotFound
	}

	return profileResp, nil
}

func (r *Repo) GetAll() (All, error) {
	var all All

	if res := r.db.Model(&entities.Doctor{}).Where("type = 'doctor'").Find(&all.Doctors); res.Error != nil {
		log.Warn(res.Error)
		return All{}, res.Error
	}

	return all, nil
}
