package user

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

type Service struct {
	repository Repository
}

// Create a new service
func NewService() *Service {
	return &Service{
		repository: *NewRepository(),
	}
}

// Get detail by id
func (this *Service) Detail(id string) (vo VODetail, err error) {
	entity, err := this.repository.Get(id)
	if err != nil {
		return vo, err
	}

	bts, _ := json.Marshal(entity)

	json.Unmarshal(bts, &vo)
	vo.Gender = this.convertGender(vo.Gender)
	vo.CreatedAt = this.convertDatetime(vo.CreatedAt)

	return vo, err
}

// Create user
func (this *Service) Create(c *gin.Context) (string, error) {
	// entity := NewEntity(c.PostForm("name"), c.PostForm())
	// err := c.ShouldBindJSON(&entity)
	// if err != nil {
	// 	return "", err
	// }

	// entity.

	return "", nil
}

// Login from platform
func (this *Service) LoginFromPlatform(mobile string, platform string) (string, error) {
	entity, err := this.repository.GetByMobile(mobile)
	if err != nil {
		return "", err
	}

	err = entity.LoginPlatform(platform)
	if err != nil {
		return "", err
	}

	// save token
	this.repository.SaveToken(entity)

	return entity.Token.getToken(), nil
}

// Convert gender to humanized word
func (this *Service) convertGender(gender string) string {
	if gender == "1" {
		gender = "male"
	} else if gender == "2" {
		gender = "female"
	} else {
		gender = "hidden"
	}

	return gender
}

// Convert time to humanized word
func (this *Service) convertDatetime(datetime string) string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation(time.RFC3339, datetime, loc)

	return t.Format("2006-01-02 15:04:05")
}
