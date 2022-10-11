package user

import (
	"encoding/json"
	"errors"
	"time"
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
func (this *Service) Create(name string, mobile string, email string, address string) (string, error) {
	entity, err := NewEntity(name, mobile, email, address)
	if err.Code != 0 {
		return "", errors.New(err.Msg)
	}

	return entity.Id, nil
}

// Remove user
func (this *Service) Remove(entity Entity) error {
	return this.repository.Remove(&entity)
}

// Login from platform
func (this *Service) Login(mobile string, platform ...string) (string, error) {
	if len(platform) == 0 {
		platform[0] = "default"
	}

	if len(platform) != 1 {
		return "", errors.New("Platform allowed only: wework, dingtalk, feishu, default or empty")
	}

	entity, err := this.repository.GetByMobile(mobile)
	if err != nil {
		return "", nil
	}

	var token string
	switch platform[0] {
	case "wework":
		token, err = entity.LoginFromWework()
	case "dingtalk":
		token, err = entity.LoginFromDingtalk()
	case "feishu":
		token, err = entity.LoginFromFeishu()
	default:
		token, err = entity.LoginFromDefault()
	}

	if err != nil {
		return "", err
	}

	err = this.repository.Save(entity)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Login authentication
func (this *Service) Verify(token string) (string, error) {
	entity, err := this.repository.GetByToken(token)
	if err != nil {
		return "", err
	}

	return entity.Id, nil
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
