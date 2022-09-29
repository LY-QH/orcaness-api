package infra

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

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
func (this *Wework) Notify(c *gin.Context) {

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
