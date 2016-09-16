package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/naokij/imgdrop/utils"
	"regexp"
	"time"
)

type User struct {
	Id       int
	Username string    `orm:"size(30);unique"`
	Password string    `orm:"size(128)"`
	Email    string    `orm:"size(80);unique"`
	Avatar   string    `orm:"size(32)"`
	Salt     string    `orm:"size(6)"`
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"auto_now;type(datetime)"`
}

const (
	activeCodeLife = 180
	resetPasswordCodeLife
	UsernameRegex = `^[a-zA-Z0-9]+$`
)

func (m *User) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *User) ValidUsername() (err error) {
	reg := regexp.MustCompile(UsernameRegex)
	if !reg.MatchString(m.Username) {
		err = errors.New("只能包含英文、数字和汉字")
	} else {
		if !(utils.HZStringLength(m.Username) >= 3 && utils.HZStringLength(m.Username) <= 16) {
			err = errors.New("长度3-16（汉字长度按2计算）")
		}
	}
	return err
}

func (m *User) SetPassword(password string) error {
	m.Salt = utils.GetRandomString(6)
	m.Password = utils.EncodeMd5(utils.EncodeMd5(password) + m.Salt)
	return nil
}

func (m *User) VerifyPassword(password string) bool {
	if m.Password == utils.EncodeMd5(utils.EncodeMd5(password)+m.Salt) {
		return true
	}
	return false
}

func Users() orm.QuerySeter {
	return orm.NewOrm().QueryTable("user")
}

func init() {
	orm.RegisterModel(new(User))
}
