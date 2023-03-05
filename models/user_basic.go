package models

import (
	"gorm.io/gorm"
	"im-instance/utils"
	"log"
)

type UserBasic struct {
	gorm.Model
	Username  string `json:"username" gorm:"unique"`
	Passwd    string `json:"passwd"`
	Email     string `json:"email" gorm:"unique"`
	Signature string `json:"signature"`
	Avatar    string `json:"avatar"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func init() {
	err := utils.DB.AutoMigrate(&UserBasic{})
	if err != nil {
		log.Fatalf("UserBasic table migrate failed, err: %+v", err)
	}
}

func GetUserList() ([]*UserBasic, error) {
	var data []*UserBasic
	db := utils.DB.Find(&data)
	return data, db.Error
}

func GetUserByID(id uint) (*UserBasic, error) {
	var user UserBasic
	db := utils.DB.Where("id = ?", id).First(&user)
	return &user, db.Error
}

func CreateUser(user *UserBasic) error {
	user.Avatar = "0.png"
	return utils.DB.Create(user).Error
}

func CheckUserByUsername(user *UserBasic) bool {
	db := utils.DB.Where("username = ? AND passwd = ?", user.Username, user.Passwd).First(&user)
	return db.RowsAffected == 1
}

func CheckUserByEmail(user *UserBasic) bool {
	db := utils.DB.Where("email = ? AND passwd = ?", user.Email, user.Passwd).First(&user)
	return db.RowsAffected == 1
}

func CheckUsername(username string) bool {
	db := utils.DB.Where("username = ?", username).First(&UserBasic{})
	return db.RowsAffected == 1
}

func CheckEmail(email string) bool {
	db := utils.DB.Where("email = ?", email).First(&UserBasic{})
	return db.RowsAffected == 1
}

func UploadAvatar(id uint, avatar string) error {
	return utils.DB.Model(&UserBasic{}).Where("id = ?", id).Update("avatar", avatar).Error
}

func GetAvatar(id uint) (string, error) {
	var user UserBasic
	db := utils.DB.Where("id = ?", id).First(&user)
	return user.Avatar, db.Error
}

func UpdateUsername(id uint, username string) error {
	return utils.DB.Model(&UserBasic{}).Where("id = ?", id).Update("username", username).Error
}

func UpdateSignature(id uint, signature string) error {
	return utils.DB.Model(&UserBasic{}).Where("id = ?", id).Update("signature", signature).Error
}

func FindUserByUsername(username string) (*UserBasic, error) {
	var user UserBasic
	db := utils.DB.Where("username = ?", username).First(&user)
	return &user, db.Error
}
