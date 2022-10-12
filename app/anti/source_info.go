package anti

import "gorm.io/datatypes"

// 开放平台用户信息
type User struct {
	CorpId   string         `json:"corp_id"`     // 平台企业ID
	UserId   string         `json:"user_id"`     // 平台用户ID
	Name     string         `json:"name"`        // 用户名
	DeptIds  []uint         `json:"dpt_ids"`     // 所处部门ID列表
	Position string         `json:"position"`    // 职位
	Mobile   string         `json:"mobile"`      // 手机号码
	Gender   string         `json:"gender"`      // 性别,0-保密,1-男,2-女
	Avatar   string         `json:"avatar"`      // 头像原图地址
	Email    string         `json:"email"`       // 邮件地址
	Status   uint           `json:"status"`      // 状态,frozen-冻结,activated-激活,resigned-离职
	Address  string         `json:"address"`     // 地址,"城市,具体地址"
	Openid   string         `json:"open_userid"` // 应用open id
	JoinTime datatypes.Date `json:"join_time"`   // 入职时间
	Super    bool           `json:"super"`       // 超级管理员
}

type Department struct {
	CorpId        string   `json:"corp_id"`         // 平台企业ID
	DeptId        string   `json:"dept_id"`         // 平台部门ID
	Name          string   `json:"name"`            // 部门名称
	ParentId      string   `json:"parent_id"`       // 上级部门ID
	LeaderUserIds []string `json:"leader_user_ids"` // 部门主管ID
}
