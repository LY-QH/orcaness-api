package Entity

import (
  "time"
)

// UserBase
type userbase struct {
  Id        uint      `gorm:"column:id;type:int unsigned;not null;primarykey" json:"id"`
  Uid       string    `gorm:"column:uid;type:varchar(36);not null" json:"uid"`                                                      // uuid
  Name      string    `gorm:"column:name;type:varchar(20);not null" json:"name"`                                                    // 用户名
  Gender    string    `gorm:"column:gender;type:enum('0','1','2');not null;default:'0'" json:"gender"`                              // 性别，0-未知，1-男，2-女
  Mobile    string    `gorm:"column:mobile;type:varchar(11);not null" json:"mobile"`                                                // 手机号码
  Email     string    `gorm:"column:email;type:varchar(50);not null" json:"email"`                                                  // 邮件地址
  Source    string    `gorm:"column:source;type:enum('wework','dingtalk','feishu','other');not null;default:'other'" json:"source"` // 来源
  CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"createdAt"`                                            // 创建时间
  UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updatedAt"`                                            // 更新时间
  DeletedAt time.Time `gorm:"column:deleted_at;type:datetime;index:idx_deletedat" json:"deletedAt"`                                 // 删除时间
}

type UserBase struct {

}

var newUserBase userbase

// Table name
func (m *UserBase) TableName() string {
  return "user_base"
}

// Create
func (m *UserBase) Create() UserBase {
  newUserBase.Id = 0
  return *m
}

// Get/Set Uid
func (m *UserBase) Uid(v ...string) string {
  if len(v) == 1 {
    newUserBase.Uid = v[0]
  }

  return newUserBase.Uid
}

// Get/Set Name
func (m *UserBase) Name(v ...string) string {
  if len(v) == 1 {
    newUserBase.Name = v[0]
  }

  return newUserBase.Name
}

// Get/Set Gender
func (m *UserBase) Gender(v ...string) string {
  if len(v) == 1 {
    newUserBase.Gender = v[0]
  }

  return newUserBase.Gender
}

// Get/Set Mobile
func (m *UserBase) Mobile(v ...string) string {
  if len(v) == 1 {
    newUserBase.Mobile = v[0]
  }

  return newUserBase.Mobile
}

// Get/Set Email
func (m *UserBase) Email(v ...string) string {
  if len(v) == 1 {
    newUserBase.Email = v[0]
  }

  return newUserBase.Email
}

// Get/Set Source
func (m *UserBase) Source(v ...string) string {
  if len(v) == 1 {
    newUserBase.Source = v[0]
  }

  return newUserBase.Source
}
