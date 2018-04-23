package models

import (
	"shorturl/backend/utils"
	"time"
)

type User struct {
	UserId   int       `xorm:"pk notnull autoincr 'u_id'"`
	UserName string    `xorm:"varchar(50) notnull 'u_name'"`
	Email    string    `xorm:"varchar(50) notnull unique 'u_email'"`
	IP       string    `xorm:"varchar(20) notnull 'u_ip'"`
	Imcode   string    `xorm:"varchar(10) notnull unique 'u_imcode'"`
	Allow    bool      `xorm:"default 0 notnull 'u_allow'"`
	Created  time.Time `xorm:"created"`
	Updated  time.Time `xorm:"updated"`
}

func InsertUser(model *User) error {
	_, err := db.Insert(model)
	return err
}

func UpdateUser(model *User) error {
	_, err := db.ID(model.UserId).Update(model)
	return err
}

func DeleteUser(id int) error {
	_, err := db.Delete(&User{UserId: id})
	return err
}

func GetUserById(id int) *User {
	var model User
	has, err := db.ID(id).Get(&model)
	if err != nil || !has {
		return nil
	}
	return &model
}

func GetUserByEmail(email string) *User {
	var model User
	has, err := db.Where("u_email = ? or u_imcode = ?", email, email).Get(&model)
	if err != nil || !has {
		return nil
	}
	return &model
}

func GetUserList(page, size int) utils.Page {
	list := make([]User, 0)
	count, _ := db.Count(&User{})
	err := db.Limit(size, (page-1)*size).Find(&list)
	if err != nil {
		panic(err)
	}
	return utils.PageUtil(int(count), page, size, list)
}
