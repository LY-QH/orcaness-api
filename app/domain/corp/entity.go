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
	Id        string             `gorm:"column:id;type:char(26);not null;primarykey" json:"id"`
	Name      string             `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Source    string             `gorm:"column:source;type:enum('wework','dingtalk','feishu','default');not null;default:'default'" json:"source"`
	SourceId  string             `gorm:"column:source_id;type:varchar(64);not null;default:''" json:"source_id"`
	ParentId  string             `gorm:"column:parent_id;type:varchar(26);not null;default:''" json:"parent_id"`
	AdminIds  []string           `gorm:"column:admin_ids;type:json" json:"admin_ids"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `gorm:"index" json:"-"`
	Events    []domain.EventBase `gorm:"-:all" json:"-"`
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

func (this *Group) TableName() string {
	return "corp_group"
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

// Push event
func (this *Group) PushEvent(action string) {
	this.Events = append(this.Events, domain.EventBase{
		Id:         util.GenId("evt."),
		ResourceId: this.Id,
		Action:     action,
		Time:       time.Now(),
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
func (this *Entity) AddGroup(name string, source string, sourceId string, parentId string) (group *Group, errcode Errcode) {
	group = &Group{Id: util.GenId("group.")}
	group.PushEvent("Created")

	errcode = group.ModifyName(name)
	if errcode.Code != 0 {
		return nil, errcode
	}

	if !util.StringInArray(source, []string{"dingtalk", "wework", "feishu", "default"}) {
		return nil, ERR_INVALID_SOURCE
	}
	group.Source = source

	if source != "default" {
		if sourceId == "" {
			return nil, ERR_INVALID_SOURCE_ID
		}

		group.SourceId = sourceId
		this.PushEvent("Updated source to: " + sourceId)
	}

	errcode = group.ModifyParent(parentId)

	return
}

// Update group
func (this *Entity) UpdateGroup(group *Group, name string, parentId string) (errcode Errcode) {

	errcode = group.ModifyName(name)
	if errcode.Code != 0 {
		return
	}

	errcode = group.ModifyParent(parentId)

	return
}

// Remove group
func (this *Entity) RemoveGroup(group *Group) error {
	group.PushEvent("Removed")
	return nil
}

// Modify group's name
func (this *Group) ModifyName(name string) (errcode Errcode) {
	this.PushEvent("Updated name to: " + name)
	return
}

// Modify group's parent
func (this *Group) ModifyParent(parentId string) (errcode Errcode) {
	// TODO: convert parentId
	if parentId == this.Id {
		return ERR_INVALID_PARENT_ID
	}

	// TODO: is child

	this.ParentId = parentId
	this.PushEvent("Updated parent id to: " + parentId)

	return
}
