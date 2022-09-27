package user

import (
	"net/mail"
	"regexp"
	"time"

	util "orcaness.com/api/util"
)

// Entity
type Entity struct {
	Id        string    `gorm:"column:id;type:char(22);not null;primarykey" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(20);not null" json:"name"`                                                    // 用户名
	Gender    string    `gorm:"column:gender;type:enum('0','1','2');not null;default:'0'" json:"gender"`                              // 性别，0-未知，1-男，2-女
	Mobile    string    `gorm:"column:mobile;type:varchar(11);not null" json:"mobile"`                                                // 手机号码
	Email     string    `gorm:"column:email;type:varchar(50);not null" json:"email"`                                                  // 邮件地址
	Source    string    `gorm:"column:source;type:enum('wework','dingtalk','feishu','other');not null;default:'other'" json:"source"` // 来源
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"createdAt"`                                            // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updatedAt"`                                            // 更新时间
	DeletedAt time.Time `gorm:"column:deleted_at;type:datetime;index:idx_deletedat" json:"deletedAt"`                                 // 删除时间
}

// Table name
func (this *Entity) TableName() string {
	return "user"
}

// Create a new user
func NewEntity(name string, mobile string, email string, source ...string) (this *Entity, err Errcode) {

	this = &Entity{Id: util.GenId("user.")}

	// name
	if err = this.UpdateName(name); err.Code != 0 {
		return nil, err
	}

	// mobile
	if err = this.UpdateMobile(mobile); err.Code != 0 {
		return nil, err
	}

	// email
	if err = this.UpdateEmail(email); err.Code != 0 {
		return nil, err
	}

	// source
	if len(source) == 0 {
		source = append(source, "other")
	}
	if err = this.UpdateSource(source[0]); err.Code != 0 {
		return nil, err
	}

	// set Default value
	this.Gender = "0"

	return this, err
}

// Update user's name
func (this *Entity) UpdateName(name string) (err Errcode) {
	nameLen := len(name)
	if nameLen < 2 {
		return ERR_NAME_LEN_LESS_THAN_MINI_LIMIT
	}
	if nameLen > 20 {
		return ERR_NAME_LEN_GREATER_THAN_MAX_LIMIT
	}

	reg := regexp.MustCompile(`^[a-z0-9_\-\.]+$`)
	if !reg.MatchString(name) {
		return ERR_NAME_CONTAINS_ILLEGAL_CHARS
	}

	this.Name = name
	return
}

// Update user's mobile
func (this *Entity) UpdateMobile(mobile string) (err Errcode) {
	reg := regexp.MustCompile(`1[3456789][0-9]{9}`)
	if !reg.MatchString(mobile) {
		return ERR_INVALID_MOBILE_FORMAT
	}

	this.Mobile = mobile
	return
}

// Update user's email
func (this *Entity) UpdateEmail(email string) (err Errcode) {
	if _, err := mail.ParseAddress(email); err != nil {
		return ERR_INVALID_EMAIL_FORMAT
	}

	this.Email = email
	return
}

// Update user's source
func (this *Entity) UpdateSource(source string) (err Errcode) {
	if source != "dingtalk" && source != "wework" && source != "feishu" && source != "other" {
		return ERR_INVALID_SOURCE
	}

	this.Source = source

	return
}

// Set user's gender to male
func (this *Entity) SetToMale() (err Errcode) {
	this.Gender = "1"
	return
}

// Set user's gender to female
func (this *Entity) SetToFemale() (err Errcode) {
	this.Gender = "2"
	return
}
