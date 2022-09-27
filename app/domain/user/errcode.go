package user

type Errcode struct {
	Code uint   `json:"errcode"`
	Msg  string `json:"errmsg"`
}

var (
	ERR_NAME_LEN_LESS_THAN_MINI_LIMIT = Errcode{
		Code: 5001,
		Msg:  "Name length less then 2 chars",
	}
	ERR_NAME_LEN_GREATER_THAN_MAX_LIMIT = Errcode{
		Code: 5002,
		Msg:  "Name length greater than 20 chars",
	}
	ERR_NAME_CONTAINS_ILLEGAL_CHARS = Errcode{
		Code: 5003,
		Msg:  "Name contains illegal chars, allowed chars: [a-z,0-9,-,_,.]",
	}
	ERR_INVALID_MOBILE_FORMAT = Errcode{
		Code: 5004,
		Msg:  "Invalid mobile format",
	}
	ERR_INVALID_EMAIL_FORMAT = Errcode{
		Code: 5005,
		Msg:  "Invalid email format",
	}
	ERR_INVALID_SOURCE = Errcode{
		Code: 5005,
		Msg:  "Invalid source, allowed source: [dingtalk, wework, feishu, other]",
	}
)
