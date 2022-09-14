package WeworkDomain

import (
  "time"

  "gorm.io/datatypes"
)

// WeworkUser
type WeworkUser struct {
  Id               uint           `gorm:"column:id;type:int unsigned;not null;primarykey" json:"id"`
  Uid              string         `gorm:"column:uid;type:varchar(36);not null" json:"uid"`                                       // user uid
  Userid           string         `gorm:"column:userid;type:varchar(100);not null" json:"userid"`                                // wework userid
  Corpid           int            `gorm:"column:corpid;type:int" json:"corpid"`
  Name             string         `gorm:"column:name;type:varchar(100);not null" json:"name"`
  Alias            string         `gorm:"column:alias;type:varchar(100)" json:"alias"`
  Mobile           string         `gorm:"column:mobile;type:varchar(15)" json:"mobile"`
  Department       datatypes.JSON `gorm:"column:department;type:json;not null" json:"department"`                                // department id(int) array
  Order            datatypes.JSON `gorm:"column:order;type:json" json:"order"`                                                   // 部门排序，与department对应，值大者靠前
  Position         string         `gorm:"column:position;type:varchar(100)" json:"position"`                                     // 职位名称
  Gender           string         `gorm:"column:gender;type:enum('0','1','2');not null;default:'0'" json:"gender"`               // 性别，0-未知，1-男，2-女
  Email            string         `gorm:"column:email;type:varchar(64)" json:"email"`
  BizMail          string         `gorm:"column:biz_mail;type:varchar(64)" json:"bizMail"`                                       // 企业邮箱
  Telephone        string         `gorm:"column:telephone;type:varchar(32)" json:"telephone"`
  IsLeaderInDept   datatypes.JSON `gorm:"column:is_leader_in_dept;type:json" json:"isLeaderInDept"`                              // 与department对应，0-非负责人，1-负责人
  DirectLeader     datatypes.JSON `gorm:"column:direct_leader;type:json" json:"directLeader"`                                    // 直属上级userid
  Avatar           string         `gorm:"column:avatar;type:varchar(150)" json:"avatar"`                                         // 头像url
  ThumbAvatar      string         `gorm:"column:thumb_avatar;type:varchar(150)" json:"thumbAvatar"`                              // 头像缩略图
  Address          string         `gorm:"column:address;type:varchar(128)" json:"address"`                                       // 住址
  MainDepartment   uint8          `gorm:"column:main_department;type:tinyint unsigned;not null;default:0" json:"mainDepartment"` // 0-非主部门，1-主部门
  ToInvite         uint8          `gorm:"column:to_invite;type:tinyint unsigned;not null;default:1" json:"toInvite"`             // 是否邀请新同事，企微接口为boolean，默认true
  ExternalPosition string         `gorm:"column:external_position;type:varchar(100)" json:"externalPosition"`                    // 对外职务
  Extattr          datatypes.JSON `gorm:"column:extattr;type:json" json:"extattr"`                                               // 自定义属性
  ExternalProfile  datatypes.JSON `gorm:"column:external_profile;type:json" json:"externalProfile"`                              // 对外属性
  Status           uint8          `gorm:"column:status;type:tinyint unsigned;not null;default:4" json:"status"`                  // 状态，1-已激活，2-已禁用，4-未激活，5-退出企业
  QrCode           string         `gorm:"column:qr_code;type:varchar(150)" json:"qrCode"`                                        // 员工二维码
  CreatedAt        time.Time      `gorm:"column:created_at;type:datetime;not null" json:"createdAt"`
  UpdatedAt        time.Time      `gorm:"column:updated_at;type:datetime;not null" json:"updatedAt"`
  DeletedAt        time.Time      `gorm:"column:deleted_at;type:datetime;index:idx_deletedat" json:"deletedAt"`
}

// Table name
func (m *WeworkUser) TableName() string {
  return "wework_user"
}
