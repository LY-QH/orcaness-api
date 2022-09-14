package Entity

import (
  "time"
)

// PlatformWeworkCorp
type PlatformWeworkCorp struct {
  Id         string    `gorm:"column:id;type:smallint;not null;primarykey" json:"id"`
  Name       string    `gorm:"column:name;type:varchar(128);not null" json:"name"`                   // 企业名称
  Corpid     string    `gorm:"column:corpid;type:varchar(100);not null" json:"corpid"`               // 企业id
  Corpsecret string    `gorm:"column:corpsecret;type:varchar(150);not null" json:"corpsecret"`
  CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null" json:"createdAt"`
  UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updatedAt"`
  DeletedAt  time.Time `gorm:"column:deleted_at;type:datetime;index:idx_deletedat" json:"deletedAt"`
}

// Table name
func (m *PlatformWeworkCorp) TableName() string {
  return "platform_wework_corp"
}
