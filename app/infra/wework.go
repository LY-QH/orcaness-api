package infra

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"orcaness.com/api/util"
)

type Wework struct {
	Departments []Department `json:"-"`
}

type Department struct {
	gorm.Model
	Name     string   `gorm:"column:name;type:varchar(128);not null" json:"name"`
	ParentId uint8    `gorm:"column:parent_id;type:smallint;not null;default:0" json:"parent_id"`
	Members  []Member `gorm:"-:all" json:"-"`
}

type Member struct {
}

type NotifyStruct struct {
	ToUserName string
	AgentID    string
	Encrypt    string
}

type NotifyDetailStruct struct {
	ToUserName   string // 企业微信CorpID
	FromUserName string // 此事件该值固定为sys，表示该消息由系统生成
	CreateTime   string // 消息创建时间 （整型）
	MsgType      string // 消息的类型，此时固定为event
	Event        string // 事件的类型，此时固定为change_contact
}

// Create new wework object
func NewWework() *Wework {
	return &Wework{}
}

// Notify
func (this *Wework) Notify(c *gin.Context) (string, interface{}, error) {
	msg_signature := c.Query("msg_signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	if c.Request.Method == "GET" && echostr != "" {
		strs := []string{viper.GetString("wework.contacts_token"), timestamp, nonce, echostr}
		sort.Strings(strs)
		dev_msg_signature := fmt.Sprintf("%02x", sha1.Sum([]byte(strings.Join(strs, ""))))
		if dev_msg_signature == msg_signature {
			decrypted, err := this.decrypt(echostr)
			if err != nil {
				return "echo_str", "", err
			}

			content := decrypted[16:]
			// fmt.Print("decrypted:", binary.BigEndian.Uint32(content[0:4]))
			msg_len := binary.BigEndian.Uint32(content[0:4])
			msg := content[4 : msg_len+4]
			// receveid := content[msg_len+4:]
			return "echo_str", msg, nil
		}
	} else {
		if c.Request.Method != "POST" {
			return "", "", errors.New("Invalid request method")
		}

		data := NotifyStruct{}
		xml.NewDecoder(c.Request.Body).Decode(&data)

		strs := []string{viper.GetString("wework.contacts_token"), timestamp, nonce, data.Encrypt}
		sort.Strings(strs)
		dev_msg_signature := fmt.Sprintf("%02x", sha1.Sum([]byte(strings.Join(strs, ""))))
		if data.ToUserName == viper.GetString("wework.corp_id") && data.AgentID == viper.GetString("wework.agent_id") && dev_msg_signature == msg_signature {
			decrypted, err := this.decrypt(data.Encrypt)
			if err != nil {
				return "", "", err
			}

			detail := NotifyDetailStruct{}
			xml.Unmarshal(decrypted, &detail)

			switch detail.MsgType {
			case "event":
				switch detail.Event {
				case "change_contact":
					this.changeContact(decrypted)
				}
			}
		}
	}
	return "", "", errors.New("Signature not match")
}

// Login
func (this *Wework) Login(c *gin.Context) (string, error) {
	accessToken, err := this.getAccessToken()
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/auth/getuserinfo?access_token=%s&code=%s", accessToken, c.Query("code"))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	var result struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
		Userid  string `json:"userid"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if result.Errcode != 0 {
		return "", errors.New(result.Errmsg)
	}

	return result.Userid, nil
}

// Login qrcode
func (this *Wework) LoginQrcode(c *gin.Context) (string, error) {
	return fmt.Sprintf("https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=%s&agentid=%s&redirect_uri=%s&state=%s", viper.GetString("wework.corp_id"), viper.GetString("wework.agent_id"), url.QueryEscape("https://"+viper.GetString("server.domain")+"/user/login/wework"), util.GenId()), nil
}

// Get access token
func (this *Wework) getAccessToken() (string, error) {
	cache := this.getCache("access_token")
	if cache != "" {
		return cache, nil
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", viper.GetString("wework.corp_id"), viper.GetString("wework.corp_secret"))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	var result struct {
		Errcode     int    `json:"errcode"`
		Errmsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if result.Errcode != 0 {
		return "", errors.New(result.Errmsg)
	}

	this.setCache("access_key", result.AccessToken, result.ExpiresIn-60)

	return result.AccessToken, nil
}

// ChangeContact
func (this *Wework) changeContact(data []byte) error {
	var detail struct {
		ChangeType string
		CreateTime int64
		UserID     string // user 事件时存在
		Department []uint // user 事件(create/update)时存在
		Id         string // party 事件时存在
		ParentId   string // party 事件(create/update)时存在
	}

	err := xml.Unmarshal(data, &detail)
	if err != nil {
		return err
	}

	switch detail.ChangeType {
	case "create_party", "update_party":
		// 获取详情
		this.getDepartment(detail.Id)
	case "delete_party":

	case "create_user", "update_user":
		// 获取详情
		this.getUser(detail.UserID, detail.Department[len(detail.Department)-1])
	case "delete_user":
	}

	return nil
}

// 获取部门详情
func (this *Wework) getDepartment(id string) {
	accessToken, err := this.getAccessToken()
	if err != nil {
		return
	}
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/department/get?access_token=%s&id=%s", accessToken, id)

	resp, err := http.Get(url)
	if err != nil {
		return
	}

	var result struct {
		Errcode    int    `json:"errcode"`
		Errmsg     string `json:"errmsg"`
		Department struct {
			Name     string `json:"name"`
			Parentid uint   `json:"parentid"`
		} `json:"department"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	if result.Errcode != 0 {
		return
	}
}

// 获取部门详情
func (this *Wework) getUser(id string, department uint) {
	accessToken, err := this.getAccessToken()
	if err != nil {
		return
	}
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/list?access_token=%s&department_id=%d", accessToken, department)

	resp, err := http.Get(url)
	if err != nil {
		return
	}

	var result struct {
		Errcode  int    `json:"errcode"`
		Errmsg   string `json:"errmsg"`
		Userlist []struct {
			Userid         string `json:"userid"`
			Name           string `json:"name"`
			Department     []uint `json:"department"`
			Position       string `json:"position"`
			Mobile         string `json:"mobile"`
			Gender         string `json:"gender"`
			Email          string `json:"email"`
			IsLeaderInDept []uint `json:"is_leader_in_dept"`
			Avatar         string `json:"avatar"`
			Status         uint   `json:"status"`
			Address        string `json:"address"`
			OpenUserid     string `json:"open_userid"`
		} `json:"userlist"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	if result.Errcode != 0 {
		return
	}

	for _, user := range result.Userlist {
		if user.Userid == id {
			user.Gender = strings.Replace(user.Gender, "1", "male", -1)
			user.Gender = strings.Replace(user.Gender, "2", "female", -1)
			return
		}
	}
}

// Get from cache
func (this *Wework) getCache(key string) string {
	return ""
}

// Set cache
func (this *Wework) setCache(key string, value string, ttl int) {

}

// key
func (this *Wework) key(key string) string {
	return key + "@wework." + viper.GetString("wework.corp_id")
}

func aesDecrypt(cipherData []byte, aesKey []byte) ([]byte, error) {
	//PKCS#7
	if len(cipherData)%len(aesKey) != 0 {
		return nil, errors.New("crypto/cipher: ciphertext size is not multiple of aes key length")
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	iv := aesKey[0:aes.BlockSize]

	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainData := make([]byte, len(cipherData))
	blockMode.CryptBlocks(plainData, cipherData)
	return plainData, nil
}

func (this *Wework) decrypt(encrypted string) ([]byte, error) {
	aes_key := viper.GetString("wework.contacts_aeskey") + "="
	aes_key_new, err := base64.StdEncoding.DecodeString(aes_key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	echostr_new, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	decrypted, err := aesDecrypt(echostr_new, aes_key_new)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
