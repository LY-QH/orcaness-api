package corp

import (
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"orcaness.com/api/app/domain"
	"orcaness.com/api/util"
)

type Entity struct {
	Id        string             `gorm:"column:id;type:char(25);not null;primarykey" json:"id"`
	Name      string             `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Logo      string             `gorm:"column:logo;type:varchar(255);not null;default:''" json:"logo"`
	Wework    datatypes.JSON     `gorm:"colum:wework;type:json" json:"wework"`
	Dingtalk  datatypes.JSON     `gorm:"colum:dingtalk;type:json" json:"dingtalk"`
	Feishu    datatypes.JSON     `gorm:"colum:feishu;type:json" json:"feishu"`
	Groups    []Group            `gorm:"-:all" json:"-"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `gorm:"index" json:"-"`
	Events    []domain.EventBase `gorm:"-:all" json:"-"`
}

type Group struct {
	Id        string              `gorm:"column:id;type:char(26);not null;primarykey" json:"id"`
	Name      string              `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Source    string              `gorm:"column:source;type:enum('wework','dingtalk','feishu','default');not null;default:'default'" json:"source"`
	SourceId  string              `gorm:"column:source_id;type:varchar(64);not null;default:''" json:"source_id"`
	ParentId  string              `gorm:"column:parent_id;type:varchar(26);not null;default:''" json:"parent_id"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	DeletedAt gorm.DeletedAt      `gorm:"index" json:"-"`
	PushEvent func(action string) `gorm:"-:all" json:"-"`
}

type WeworkConfig struct {
	CorpId          string `json:"corp_id"`
	AgentId         string `json:"agent_id"`
	CorpSecret      string `json:"corp_secret"`
	Contacts_Secret string `json:"contacts_secret"`
	Contacts_Token  string `json:"contacts_token"`
	Contacts_Aeskey string `json:"contacts_aeskey"`
}

type DingtalkConfig struct {
}

type FeishuConfig struct {
}

// corp

// Get table name
func (this *Entity) TableName() string {
	return "corp"
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

func NewCorp(name string) (entity *Entity, err error) {
	entity = &Entity{Id: util.GenId("corp.")}
	entity.PushEvent("Created")

	errcode := entity.UpdateName(name)
	if errcode.Code != 0 {
		return entity, errors.New(errcode.Msg)
	}

	return entity, nil
}

// Update corp's name
func (this *Entity) UpdateName(name string) (errcode Errcode) {
	reg := regexp.MustCompile(`[~!@#\$%^&\*\|\\，,。、\{\}/><《》?？\[\]]=`)
	if reg.MatchString(name) {
		return ERR_NAME_CONTAINS_ILLEGAL_CHARS
	}

	nameLen := len(name)
	if nameLen < 2 {
		return ERR_NAME_LEN_LESS_THAN_MINI_LIMIT
	}

	if nameLen > 100 {
		return ERR_NAME_LEN_GREATER_THAN_MAX_LIMIT
	}

	if this.Name != name {
		this.Name = name
		this.PushEvent("Modified name to: " + name)
	}

	return
}

// Update logo
func (this *Entity) UpdateLogo(logo string) (errcode Errcode) {
	if this.Logo != "" && this.Logo != logo {
		this.Logo = logo
		this.PushEvent("Modified logo to: " + logo)
	}

	return
}

// Update wework configure
func (this *Entity) UpdateWework(wework WeworkConfig) (errcode Errcode) {
	if wework.CorpId != "" && wework.AgentId != "" && wework.CorpSecret != "" && wework.Contacts_Token != "" && wework.Contacts_Secret != "" && wework.Contacts_Aeskey != "" {
		bts, _ := json.Marshal(wework)
		if string(bts) != datatypes.JSON(this.Wework).String() {
			json.Unmarshal(bts, &this.Wework)
			this.PushEvent("Modified wework configure to: " + string(bts))
		}
	}
	return
}

// Update dingtalk configure
func (this *Entity) UpdateDingtalk(dingtalk DingtalkConfig) (errcode Errcode) {
	return
}

// Update feishu configure
func (this *Entity) UpdateFeishu(feishu FeishuConfig) (errcode Errcode) {
	return
}

// Add group
func (this *Entity) AddGroup(group Group) (errcode Errcode) {
	return
}

// Modify group
func (this *Entity) ModifyGroup(group Group) (errcode Errcode) {
	return
}

// // Update corp's source
// func (this *Entity) UpdateSource(source string) (err Errcode) {
// 	if !util.StringInArray(source, []string{"dingtalk", "wework", "feishu", "default"}) {
// 		return ERR_INVALID_SOURCE
// 	}

// 	if this.Source != source {
// 		this.PushEvent("Source updated to: " + source)
// 		this.Source = source
// 	}

// 	return
// }
