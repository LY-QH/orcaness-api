package user

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"time"

	"gorm.io/gorm"
	domain "orcaness.com/api/app/domain"
	util "orcaness.com/api/util"
)

// Entity
type Entity struct {
	Id          string             `gorm:"column:id;type:char(25);not null;primarykey" json:"id"`
	Name        string             `gorm:"column:name;type:varchar(20);not null" json:"name"`                                                    // 用户名
	Avatar      string             `gorm:"column:avatar;type:varchar(255);not null;default:''"                                    json:"avatar"` // 头像
	Gender      string             `gorm:"column:gender;type:enum('0','1','2');not null;default:'0'" json:"gender"`                              // 性别，0-未知，1-男，2-女
	Mobile      string             `gorm:"column:mobile;type:varchar(11);not null" json:"mobile"`                                                // 手机号码
	Email       string             `gorm:"column:email;type:varchar(50);not null" json:"email"`                                                  // 邮件地址
	Address     string             `gorm:"column:address;type:varchar(255);not null;default:''" json:"address"`                                  // 邮件地址
	FromSources []FromSource       `gorm:"-:all" json:"-"`                                                                                       // 来源
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `gorm:"index" json:"-"`
	Events      []domain.EventBase `gorm:"-:all" json:"-"`
	Tokens      []Token            `gorm:"-:all" json:"-"`
}

// Token
type Token struct {
	Id        string             `gorm:"column:id;type:char(23);not null;primarykey" json:"id"`
	UserId    string             `gorm:"column:user_id;type:char(25);not null" json:"userid"`
	Token     string             `gorm:"column:avatar;type:char(64);not null"                                    json:"token"`                     // Token
	Source    string             `gorm:"column:source;type:enum('wework','dingtalk','feishu','default');not null;default:'default'" json:"source"` // 平台
	CreatedAt time.Time          `json:"created_at"`
	ExpiredAt time.Time          `gorm:"column:expired_at;type:datetime;not null" json:"expired_at"` // 过期时间
	Events    []domain.EventBase `gorm:"-:all" json:"-"`
}

// From Source
type FromSource struct {
	Id        string             `gorm:"column:id;type:char(23);not null;primarykey" json:"id"`
	CorpId    string             `gorm:"column:corp_id;type:char(25);not null" json:"corp_id"`
	UserId    string             `gorm:"column:user_id;type:char(25);not null" json:"userid"`
	Source    string             `gorm:"column:source;type:enum('wework','dingtalk','feishu','default');not null;default:'default'" json:"source"`
	OpenId    string             `gorm:"column:open_id;type:varchar(100);not null;default:''" json:"open_id"`
	IsSuper   uint8              `gorm:"column:is_super;type:tinyint;not null;default:0" json:"is_super"`
	InGroups  []InGroup          `gorm:"-:all" json:"-"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `gorm:"index" json:"-"`
	Events    []domain.EventBase `gorm:"-:all" json:"-"`
}

// InGroup
type InGroup struct {
	Id        string             `gorm:"column:id;type:char(23);not null;primarykey" json:"id"`
	SourceId  string             `gorm:"column:source_id;type:char(23);not null" json:"source_id"`
	GroupId   string             `gorm:"column:group_id;type:char(26);not null;default:''" json:"group_id"`
	Position  string             `gorm:"column:position;type:varchar(20);not null;default:'member'" json:"position"`
	IsAdmin   uint8              `gorm:"column:is_admin;type:tinyint;not null;default:0" json:"is_admin"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `gorm:"index" json:"-"`
	Events    []domain.EventBase `gorm:"-:all" json:"-"`
}

// Table name
func (this *Entity) TableName() string {
	return "user"
}

func (this *Token) TableName() string {
	return "user_token"
}

func (this *FromSource) TableName() string {
	return "user_from_source"
}

func (this *InGroup) TableName() string {
	return "user_in_group"
}

// Create a new user
func NewEntity(name string, mobile string, email string, address ...string) (this *Entity, err Errcode) {

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

	// address
	if len(address) == 1 && address[0] != "" {
		this.UpdateAddress(address[0])
	}

	// set Default value
	this.Gender = "0"

	return this, err
}

func NewToken(userId string, source string) *Token {
	token := &Token{
		Id:        util.GenId("tk."),
		UserId:    userId,
		Source:    source,
		ExpiredAt: time.Now().Add(90 * 24 * time.Hour),
	}

	token.Token = fmt.Sprintf("%x", sha256.Sum256([]byte(token.Id+"*"+util.GenId())))

	return token
}

func NewSource(corpId string, userId string, source string, openId string, isSuper uint8) *FromSource {
	return &FromSource{
		CorpId:  corpId,
		UserId:  userId,
		Source:  source,
		OpenId:  openId,
		IsSuper: isSuper,
		Id:      util.GenId("us."),
	}
}

func NewGroup(sourceId string, groupId string, position string, isAdmin uint8) *InGroup {
	return &InGroup{
		Id:       util.GenId("ig."),
		SourceId: sourceId,
		GroupId:  groupId,
		Position: position,
		IsAdmin:  isAdmin,
	}
}

// Push event
func (this *Entity) PushEvent(action string) {
	this.Events = append(this.Events, domain.EventBase{
		Id:         util.GenId("evt."),
		ResourceId: this.Id,
		Action:     action,
		Time:       time.Now(),
	})
}

func (this *Token) PushEvent(action string) {
	this.Events = append(this.Events, domain.EventBase{
		Id:         util.GenId("evt."),
		ResourceId: this.Id,
		Action:     action,
		Time:       time.Now(),
	})
}

func (this *FromSource) PushEvent(action string) {
	this.Events = append(this.Events, domain.EventBase{
		Id:         util.GenId("evt."),
		ResourceId: this.Id,
		Action:     action,
		Time:       time.Now(),
	})
}

func (this *InGroup) PushEvent(action string) {
	this.Events = append(this.Events, domain.EventBase{
		Id:         util.GenId("evt."),
		ResourceId: this.Id,
		Action:     action,
		Time:       time.Now(),
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
		this.PushEvent("Name updated to: " + name)
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
		this.PushEvent("Mobile updated to: " + mobile)
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
		this.PushEvent("Email updated to: " + email)
	}

	this.Email = email
	return
}

// Update user's address
func (this *Entity) UpdateAddress(address string) (err Errcode) {
	if address == "" {
		return
	}

	if this.Address != address {
		this.PushEvent("Address updated to: " + address)
	}

	this.Address = address
	return
}

// Set user's gender to male
func (this *Entity) SetToMale() (err Errcode) {
	if this.Gender != "1" {
		this.PushEvent("Gender updated to: male")
	}

	this.Gender = "1"
	return
}

// Set user's gender to female
func (this *Entity) SetToFemale() (err Errcode) {
	if this.Gender != "2" {
		this.PushEvent("Gender updated to: female")
	}

	this.Gender = "2"
	return
}

// Hide user's gender
func (this *Entity) HideGender() (err Errcode) {
	if this.Gender != "0" {
		this.PushEvent("Gender been hidden")
	}

	this.Gender = "0"
	return
}

func (this *Entity) AddSource(corpId string, source string, openId string, isSuper uint8) error {
	for _, fromSource := range this.FromSources {
		if fromSource.CorpId == corpId && fromSource.Source == source && fromSource.OpenId == openId {
			return errors.New("Source duplicate")
		}
	}

	this.FromSources = append(this.FromSources, *NewSource(corpId, this.Id, source, openId, isSuper))

	return nil
}

func (this *Entity) UpdateAvatar(avatar string) error {
	if avatar != "" {
		this.Avatar = avatar
		this.PushEvent("Update avatar to: " + avatar)
	}

	return nil
}

func (this *Entity) LoginFromWework() (string, error) {
	return this.loginFromSource("wework")
}

func (this *Entity) LoginFromDingtalk() (string, error) {
	return this.loginFromSource("dingtalk")
}

func (this *Entity) LoginFromFeishu() (string, error) {
	return this.loginFromSource("feishu")
}

func (this *Entity) LoginFromDefault() (string, error) {
	return this.loginFromSource("default")
}

func (this *Entity) RevokeToken(token string) {
	for _, tk := range this.Tokens {
		if tk.Token == token {
			tk.PushEvent("Revoked")
			break
		}
	}
}

func (this *Entity) RemoveSource(source string) {
	for _, src := range this.FromSources {
		if src.Source == source {
			src.PushEvent("Removed")
			break
		}
	}
}

func (this *FromSource) InGroup(groupId string, position string, isAdmin bool) {
	var isAdminInt uint8
	if isAdmin {
		isAdminInt = 1
	}
	group := NewGroup(this.Id, groupId, position, isAdminInt)
	group.PushEvent("Created")
	this.InGroups = append(this.InGroups, *group)
}

func (this *FromSource) OutGroup(groupId string) {
	for _, group := range this.InGroups {
		if group.Id == groupId {
			group.PushEvent("Outed")
			break
		}
	}
}

func (this *Entity) loginFromSource(source string) (string, error) {
	newToken := *NewToken(this.Id, source)
	this.Tokens = append(this.Tokens, newToken)
	this.PushEvent("Logged in " + source)
	return newToken.Token, nil
	return "", nil
}
