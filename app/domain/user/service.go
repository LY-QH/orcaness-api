package user

import (
	"encoding/json"
	"errors"
	"time"

	"orcaness.com/api/app/anti"
	"orcaness.com/api/app/domain/corp"
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
func (this *Service) Create(name string, mobile string, email string, address string) (id string, err error) {
	entity, errcode := NewEntity(name, mobile, email, address)
	if errcode.Code != 0 {
		return "", errors.New(errcode.Msg)
	}

	err = this.repository.Save(entity)
	if err != nil {
		return "", err
	}

	return entity.Id, nil
}

// Create from source
func (this *Service) CreateFromSource(source string, data anti.User) (id string, err error) {
	id, err = this.Create(data.Name, data.Mobile, data.Email, data.Address)

	if err != nil {
		return "", err
	}

	entity, err := this.repository.Get(id)
	if err != nil {
		return "", err
	}

	// gender
	if data.Gender == "1" {
		entity.SetToMale()
	} else if data.Gender == "2" {
		entity.SetToFemale()
	}

	// avatar
	if data.Avatar != "" {
		entity.UpdateAvatar(data.Avatar)
	}

	// find corp
	corpRespository := corp.NewRepository()
	corpEntitys, err := corpRespository.GetAll("?", data.CorpId)
	if err != nil {
		return
	}

	corpEntity := (*corpEntitys)[0]

	s := NewSource(corpEntity.Id, id, source, data.Openid, data.Super)

	for _, dptId := range data.DeptIds {
		groups, err := corpRespository.GetAllGroup("source = ? and source_id = ?", source, dptId)
		if err != nil {
			return id, err
		}

		group := (*groups)[0]

		g := NewGroup(s.Id, group.Id, data.Position)
		s.InGroup(*g)
	}

	entity.AddSource(*s)
	err = this.repository.Save(entity)

	return
}

// Update from source
func (this *Service) UpdateFromSource(source string, data anti.User) (id string, err error) {
	entity, err := this.repository.Get(id)
	if err != nil {
		return "", err
	}

	// gender
	if data.Gender != entity.Gender {
		if data.Gender == "1" {
			entity.SetToMale()
		} else if data.Gender == "2" {
			entity.SetToFemale()
		} else {
			entity.HideGender()
		}
	}

	// avatar
	if data.Avatar != entity.Avatar {
		entity.UpdateAvatar(data.Avatar)
	}

	// find corp
	corpRespository := corp.NewRepository()
	corpEntitys, err := corpRespository.GetAll("?", data.CorpId)
	if err != nil {
		return
	}

	corpEntity := (*corpEntitys)[0]

	sources, err := service.repository.GetAllSource(entity.Id, "corp_id = ? and source = ? and open_id = ?", corpEntity.Id, source, data.Openid)
	if len(*sources) == 0 {
		s := NewSource(corpEntity.Id, id, source, data.Openid, data.Super)

		for _, dptId := range data.DeptIds {
			groups, err := corpRespository.GetAllGroup("source = ? and source_id = ?", source, dptId)
			if err != nil {
				return id, err
			}

			group := (*groups)[0]

			g := NewGroup(s.Id, group.Id, data.Position)
			s.InGroup(*g)
		}

		entity.AddSource(*s)
	} else {
		s := (*sources)[0]
		groupIds := []string{}
		for _, dptId := range data.DeptIds {
			groups, err := corpRespository.GetAllGroup("source = ? and source_id = ?", source, dptId)
			if err != nil {
				return id, err
			}

			group := (*groups)[0]
			groupIds = append(groupIds, group.Id)

			gs, err := service.repository.GetAllGroup(s.Id, "group_id = ?", group.Id)

			if len(*gs) == 0 {
				g := NewGroup(s.Id, group.Id, data.Position)
				s.InGroup(*g)
			} else {
				g := (*gs)[0]
				if g.Position != data.Position {
					s.UpdateInGroup(g, data.Position)
				}
			}
		}

		// 删除多余数据
		gs, _ := service.repository.GetAllGroup(s.Id, "user_source_id = ? and group_id not in ?", s.Id, groupIds)
		if len(*gs) > 0 {
			for _, g := range *gs {
				s.OutGroup(g)
			}
		}

		entity.UpdateSource(s)
	}

	err = this.repository.Save(entity)

	return
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
