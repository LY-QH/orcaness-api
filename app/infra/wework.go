package infra

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
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

// Create new wework object
func NewWework() *Wework {
	return &Wework{}
}

// Notify
func (this *Wework) Notify(c *gin.Context) (interface{}, error) {
	msg_signature := c.Query("msg_signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	if c.Request.Method == "GET" && echostr != "" {
		strs := []string{viper.GetString("wework.contacts_token"), timestamp, nonce, echostr}
		sort.Strings(strs)
		dev_msg_signature := fmt.Sprintf("%02x", sha1.Sum([]byte(strings.Join(strs, ""))))
		if dev_msg_signature == msg_signature {
			aes_key := viper.GetString("wework.contacts_aeskey") + "="
			aes_key_new, err := base64.StdEncoding.DecodeString(aes_key)
			if err != nil {
				fmt.Println(err)
				return "", err
			}

			echostr_new, err := base64.StdEncoding.DecodeString(echostr)
			if err != nil {
				fmt.Println(err)
				return "", err
			}

			decrypted, err := aesDecrypt(echostr_new, aes_key_new)
			if err != nil {
				return "", err
			}

			// content := decrypted[16:]
			// fmt.Print("decrypted:", binary.BigEndian.Uint32(content[0:4]))
			// msg_len := binary.BigEndian.Uint32(content[0:4])
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }
			// msg := content[4 : msg_len+4]
			// receveid := content[msg_len+4:]
			// fmt.Printf("msg:%s\nreceiveid:%s\n", msg, receveid)
			// fmt.Printf("len: %d", len(content)-8)
			return decrypted[20 : len(decrypted)-8], nil
		}
	}
	return "", nil
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
