package UserDomain

import (
  "time"

  "gorm.io/gorm"
)

// UserBaseEntity
type UserBaseEntity struct {
  gorm.Model
  Id        uint      `gorm:"column:id;type:int unsigned;not null" json:"id"`
  Uid       string    `gorm:"column:uid;type:varchar(36);not null" json:"uid"`                                                      // uuid
  Name      string    `gorm:"column:name;type:varchar(20);not null" json:"name"`                                                    // 用户名
  Gender    string    `gorm:"column:gender;type:enum('0','1','2');not null;default:'0'" json:"gender"`                              // 性别，0-未知，1-男，2-女
  Mobile    string    `gorm:"column:mobile;type:varchar(11);not null" json:"mobile"`                                                // 手机号码
  Email     string    `gorm:"column:email;type:varchar(50);not null" json:"email"`                                                  // 邮件地址
  Source    string    `gorm:"column:source;type:enum('wework','dingtalk','feishu','other');not null;default:'other'" json:"source"` // 来源
  CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"createdAt"`                                            // 创建时间
  UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updatedAt"`                                            // 更新时间
  DeletedAt time.Time `gorm:"column:deleted_at;type:datetime" json:"deletedAt"`                                                     // 删除时间
}

// Table name
func (m *UserBaseEntity) TableName() string {
  return "user_base_entity"
}
