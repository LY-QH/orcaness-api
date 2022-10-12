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
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/datatypes"
	"orcaness.com/api/app/anti"
	"orcaness.com/api/util"
)

type Wework struct {
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
	corpid := c.Param("corpid")
	msg_signature := c.Query("msg_signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	corpConfig, err := this.getConfig(corpid)
	if err != nil {
		return "", "", err
	}

	if c.Request.Method == "GET" && echostr != "" {
		strs := []string{corpConfig.ContactsToken, timestamp, nonce, echostr}
		sort.Strings(strs)
		dev_msg_signature := fmt.Sprintf("%02x", sha1.Sum([]byte(strings.Join(strs, ""))))
		if dev_msg_signature == msg_signature {
			decrypted, err := this.decrypt(corpConfig.ContactsAeskey, echostr)
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

		strs := []string{corpConfig.ContactsToken, timestamp, nonce, data.Encrypt}
		sort.Strings(strs)
		dev_msg_signature := fmt.Sprintf("%02x", sha1.Sum([]byte(strings.Join(strs, ""))))
		if data.ToUserName == corpid && data.AgentID == corpConfig.AgentId && dev_msg_signature == msg_signature {
			decrypted, err := this.decrypt(corpConfig.ContactsAeskey, data.Encrypt)
			if err != nil {
				return "", "", err
			}

			detail := NotifyDetailStruct{}
			xml.Unmarshal(decrypted, &detail)

			switch detail.MsgType {
			case "event":
				switch detail.Event {
				case "change_contact":
					return this.changeContact(corpid, decrypted)
				}
			}
		}
	}
	return "", "", errors.New("Signature not match")
}

// Login
func (this *Wework) Login(c *gin.Context) (string, error) {
	corpid := c.Param("corpid")
	accessToken, err := this.getAccessToken(corpid)
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
	corpid := c.Param("corpid")
	corpConfig, err := this.getConfig(corpid)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=%s&agentid=%s&redirect_uri=%s&state=%s", corpid, corpConfig.AgentId, url.QueryEscape("https://"+viper.GetString("server.domain")+"/user/login/wework/"+corpid), util.GenId()), nil
}

// Get access token
func (this *Wework) getAccessToken(corpid string) (string, error) {
	cache := this.getCache("access_token")
	if cache != "" {
		return cache, nil
	}

	corpConfig, err := this.getConfig(corpid)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpid, corpConfig.CorpSecret)
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
func (this *Wework) changeContact(corpid string, body []byte) (event string, data interface{}, err error) {
	var detail struct {
		ChangeType string
		CreateTime int64
		UserID     string // user 事件时存在
		Department []uint // user 事件(create/update)时存在
		Id         string // party 事件时存在
		ParentId   string // party 事件(create/update)时存在
	}

	err = xml.Unmarshal(body, &detail)
	if err != nil {
		return
	}

	switch detail.ChangeType {
	case "create_party", "update_party":
		data, err = this.getDepartment(corpid, detail.Id)
		return detail.ChangeType, data, err
	case "delete_party":
		return detail.ChangeType, detail.Id, err
	case "create_user", "update_user":
		data, err = this.getUser(corpid, detail.UserID, detail.Department[len(detail.Department)-1])
		return detail.ChangeType, data, err
	case "delete_user":
		return detail.ChangeType, detail.UserID, err
	}

	return
}

// 获取部门详情
func (this *Wework) getDepartment(corpid string, id string) (dept anti.Department, err error) {
	accessToken, err := this.getAccessToken(corpid)
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
			Id               int      `json:"id"`
			Name             string   `json:"name"`
			Parentid         int      `json:"parentid"`
			DepartmentLeader []string `json:"department_leader"`
		} `json:"department"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	if result.Errcode != 0 {
		return
	}

	dept = anti.Department{}
	dept.CorpId = corpid
	dept.DeptId = strconv.Itoa(result.Department.Id)
	dept.LeaderUserIds = result.Department.DepartmentLeader
	dept.Name = result.Department.Name
	dept.ParentId = strconv.Itoa(result.Department.Parentid)

	return
}

// 获取部门详情
func (this *Wework) getUser(corpid string, id string, department uint) (user anti.User, err error) {
	accessToken, err := this.getAccessToken(corpid)
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

	for _, u := range result.Userlist {
		if u.Userid == id {
			user = anti.User{}
			user.CorpId = corpid
			user.UserId = id
			user.Name = u.Name
			user.Mobile = u.Mobile
			user.Email = u.Email
			user.Avatar = u.Avatar
			user.DeptIds = u.Department
			user.Gender = u.Gender
			user.Openid = u.OpenUserid
			user.Position = u.Position
			user.Address = u.Address
			user.JoinTime = datatypes.Date(time.Now())
			return
		}
	}

	return
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

type CorpConfigStruct struct {
	CorpId         string `json:"corp_id"`
	AgentId        string `json:"agent_id"`
	CorpSecret     string `json:"corp_secret"`
	ContactsSecret string `json:"contacts_secret"`
	ContactsToken  string `json:"contacts_token"`
	ContactsAeskey string `json:"contacts_aeskey"`
}

func (this *Wework) getConfig(corpid string) (config CorpConfigStruct, err error) {
	result := map[string]interface{}{}
	config = CorpConfigStruct{}

	Db("read").Table("corp").Where("").Take(&result)
	if result != nil && result["wework"] != nil {
		return
	}

	json.Unmarshal([]byte(result["wework"].(string)), &config)
	return
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

func (this *Wework) decrypt(aes_key string, encrypted string) ([]byte, error) {
	aes_key += "="
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
