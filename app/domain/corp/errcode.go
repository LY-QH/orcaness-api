package corp

type Errcode struct {
	Code uint   `json:"errcode"`
	Msg  string `json:"errmsg"`
}

var (
	ERR_NAME_LEN_LESS_THAN_MINI_LIMIT = Errcode{
		Code: 6001,
		Msg:  "Name length less then 2 chars",
	}
	ERR_NAME_LEN_GREATER_THAN_MAX_LIMIT = Errcode{
		Code: 6002,
		Msg:  "Name length greater than 100 chars",
	}
	ERR_NAME_CONTAINS_ILLEGAL_CHARS = Errcode{
		Code: 6003,
		Msg:  "Name contains illegal chars",
	}
	ERR_INVALID_SOURCE = Errcode{
		Code: 6004,
		Msg:  "Invalid source, allowed source: [dingtalk, wework, feishu, default]",
	}
	ERR_INVALID_SOURCE_ID = Errcode{
		Code: 6005,
		Msg:  "Invalid source id",
	}
	ERR_INVALID_PARENT_ID = Errcode{
		Code: 6006,
		Msg:  "Invalid parent id",
	}
)
