package models

import (
	"fmt"
	"gorm.io/gorm"
	"im-instance/utils"
	"log"
)

// 私聊 Type 标记为 single
// 群聊 Type 标记为 group_${group_id}，ToID 为 group_id

type Contact struct {
	gorm.Model
	FromID uint
	ToID   uint
	Type   string
	Status string
}

func (table *Contact) TableName() string {
	return "contact"
}

func init() {
	err := utils.DB.AutoMigrate(&Contact{})
	if err != nil {
		log.Fatalf("auto migrate contact table error: %+v", err)
	}
}

func GetContactList(fromId uint) ([]*UserBasic, error) {
	var users []*UserBasic
	db := utils.DB.
		Model(&UserBasic{}).
		Select("user_basic.*").
		Joins("left join contact on user_basic.id = contact.to_id").
		Where("contact.from_id = ? and contact.deleted_at is null", fromId).
		Scan(&users)
	for _, user := range users {
		user.Passwd = ""
	}
	return users, db.Error
}

func CheckContact(fromId, toId uint) bool {
	db := utils.DB.
		Where("from_id = ? AND to_id = ?", fromId, toId).
		Where("status = ?", "agree").
		First(&Contact{})
	return db.RowsAffected == 1
}

func Invite(fromID, toID uint) error {
	if db := utils.DB.Delete(&Contact{}, "from_id = ? AND to_id = ?", fromID, toID); db.Error != nil {
		return db.Error
	} else {
		fmt.Println(db.RowsAffected)
	}
	if db := utils.DB.Delete(&Contact{}, "from_id = ? AND to_id = ?", toID, fromID); db.Error != nil {
		return db.Error
	} else {
		fmt.Println(db.RowsAffected)
	}

	if db := utils.DB.Create(&Contact{
		FromID: fromID,
		ToID:   toID,
		Type:   "single",
		Status: "invite",
	}); db.Error != nil {
		return db.Error
	}
	if db := utils.DB.Create(&Contact{
		FromID: toID,
		ToID:   fromID,
		Type:   "single",
		Status: "invite",
	}); db.Error != nil {
		return db.Error
	}
	return nil
}

func Agree(fromID, toID uint) error {
	db := utils.DB.
		Model(&Contact{}).
		Where("from_id = ? AND to_id = ?", fromID, toID).
		Update("status", "agree")
	if db.Error != nil {
		return db.Error
	}
	db = utils.DB.
		Model(&Contact{}).
		Where("from_id = ? AND to_id = ?", toID, fromID).
		Update("status", "agree")
	return db.Error
}

func Disagree(fromID, toID uint) error {
	db := utils.DB.
		Model(&Contact{}).
		Where("from_id = ? AND to_id = ?", fromID, toID).
		Update("status", "disagree")
	if db.Error != nil {
		return db.Error
	}
	db = utils.DB.
		Model(&Contact{}).
		Where("from_id = ? AND to_id = ?", toID, fromID).
		Update("status", "disagree")
	return db.Error
}

func GetContactStatus(fromID uint, toID uint) (string, error) {
	var contact Contact
	db := utils.DB.
		Where("from_id = ? AND to_id = ?", fromID, toID).
		First(&contact)
	return contact.Status, db.Error
}
