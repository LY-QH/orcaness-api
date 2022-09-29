package user

import (
	"crypto/sha256"
	"net/mail"
	"regexp"
	"time"

	"gorm.io/gorm"
	domain "orcaness.com/api/app/domain"
	util "orcaness.com/api/util"
)

// Entity
type Entity struct {
	Id        string             `gorm:"column:id;type:char(25);not null;primarykey" json:"id"`
	Name      string             `gorm:"column:name;type:varchar(20);not null" json:"name"`                                                        // 用户名
	Avatar    string             `gorm:"column:avatar;type:varchar(255);not null;default:''"                                    json:"avatar"`     // 头像
	Gender    string             `gorm:"column:gender;type:enum('0','1','2');not null;default:'0'" json:"gender"`                                  // 性别，0-未知，1-男，2-女
	Mobile    string             `gorm:"column:mobile;type:varchar(11);not null" json:"mobile"`                                                    // 手机号码
	Email     string             `gorm:"column:email;type:varchar(50);not null" json:"email"`                                                      // 邮件地址
	Source    string             `gorm:"column:source;type:enum('wework','dingtalk','feishu','default');not null;default:'default'" json:"source"` // 来源
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `gorm:"index" json:"-"`
	Events    []domain.EventBase `gorm:"-:all" json:"-"`
	Token     Token              `gorm:"-:all" json:"-"`
}

// Table name
func (this *Entity) TableName() string {
	return "user"
}

// Create a new user
func NewEntity(name string, mobile string, email string, source ...string) (this *Entity, err Errcode) {

	this = &Entity{Id: util.GenId("user.")}
	this.PushEvent("Created")

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
		source = append(source, "default")
	}
	if err = this.UpdateSource(source[0]); err.Code != 0 {
		return nil, err
	}

	// set Default value
	this.Gender = "0"

	return this, err
}

// Push event
func (this *Entity) PushEvent(action string) {
	this.Events = append(this.Events, domain.EventBase{
		Id:       util.GenId("evt."),
		SourceId: this.Id,
		Action:   action,
		Time:     time.Now(),
	})
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

	if this.Name != name {
		this.PushEvent("Name updated to " + name)
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

	if this.Mobile != mobile {
		this.PushEvent("Mobile updated to " + mobile)
	}

	this.Mobile = mobile
	return
}

// Update user's email
func (this *Entity) UpdateEmail(email string) (err Errcode) {
	if _, err := mail.ParseAddress(email); err != nil {
		return ERR_INVALID_EMAIL_FORMAT
	}

	if this.Email != email {
		this.PushEvent("Email updated to " + email)
	}

	this.Email = email
	return
}

// Update user's source
func (this *Entity) UpdateSource(source string) (err Errcode) {
	if !util.StringInArray(source, []string{"dingtalk", "wework", "feishu", "default"}) {
		return ERR_INVALID_SOURCE
	}

	if this.Source != source {
		this.PushEvent("Source updated to " + source)
	}

	this.Source = source
	return
}

// Set user's gender to male
func (this *Entity) SetToMale() (err Errcode) {
	if this.Gender != "1" {
		this.PushEvent("Gender updated to male")
	}

	this.Gender = "1"
	return
}

// Set user's gender to female
func (this *Entity) SetToFemale() (err Errcode) {
	if this.Gender != "2" {
		this.PushEvent("Gender updated to female")
	}

	this.Gender = "2"
	return
}

// Hide user's gender
func (this *Entity) HideGender() (err Errcode) {
	if this.Gender != "0" {
		this.PushEvent("Hide gender")
	}

	this.Gender = "0"
	return
}

// Login platform
func (this *Entity) LoginPlatform(platform string) error {
	this.genToken(platform)
	this.PushEvent("Logged in " + platform)
	return nil
}

// Generate token
func (this *Entity) genToken(platform string) {
	token := NewToken()
	token.UserId = this.Id
	token.Platform = platform
	token.Token = sha256.Sum256([]byte(this.Id + "*" + util.GenId()))
	this.Token = token
}

//////////////// token ////////////////
// Token
type Token struct {
	Id        string    `gorm:"column:id;type:char(23);not null;primarykey" json:"id"`
	UserId    string    `gorm:"column:user_id;type:char(25);not null" json:"user_id"`                                                     // 用户 id
	Token     string    `gorm:"column:avatar;type:char(64);not null"                                    json:"token"`                     // Token
	Platform  string    `gorm:"column:source;type:enum('wework','dingtalk','feishu','default');not null;default:'default'" json:"source"` // 平台
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `gorm:"column:expired_at;type:datetime;not null" json:"expired_at"` // 过期时间
}

// Table name
func (this *Token) TableName() string {
	return "user_token"
}

func NewToken() *Token {
	return &Token{
		Id: util.GenId("tk."),
		ExpiredAt: time.Now().Add(90 * 24 * time.Hour),
	}
}

func (this *Token) getToken() string {
	return this.Token
}
