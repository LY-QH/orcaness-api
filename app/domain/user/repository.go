package user

import (
	"time"

	domain "orcaness.com/api/app/domain"
	infra "orcaness.com/api/app/infra"
)

type Repository struct {
	domain.RepositoryBase
}

// Create new repository
func NewRepository() *Repository {
	return &Repository{}
}

// Get one entity by id
func (this *Repository) Get(id string) (*Entity, error) {
	entity := &Entity{}
	infra.Db("read").Where("id = ?", id).First(entity)
	return entity, nil
}

// Get one entity by mobile
func (this *Repository) GetByMobile(mobile string) (*Entity, error) {
	entity := &Entity{}
	infra.Db("read").Where("mobile = ?", mobile).First(entity)
	return entity, nil
}

// Get one entity by token
func (this *Repository) GetByToken(token string) (entity *Entity, err error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tokenObj := Token{}
	infra.Db("read").Where("token = ? and expired_at > ?", token, time.Now().In(loc).Format("2006-01-02 15:04:05")).First(&tokenObj)
	if tokenObj.Id == "" {
		return entity, nil
	}

	entity, err = this.Get(tokenObj.UserId)
	return entity, err
}

// Get one entity by source
func (this *Repository) GetBySource(corpId string, source string, openId string) (entity *Entity, err error) {
	sourceObj := FromSource{}
	infra.Db("read").Where("corp_id = ? and source = ? and open_id = ?", corpId, source, openId).First(&sourceObj)
	if sourceObj.Id == "" {
		return entity, nil
	}

	entity, err = this.Get(sourceObj.UserId)
	return entity, err
}

// Get entity list by query condition
func (this *Repository) GetAll(query ...interface{}) (*[]Entity, error) {
	entities := &[]Entity{}
	db := infra.Db("read")
	if len(query) > 0 {
		var args []interface{}
		for i, q := range query {
			if i > 0 {
				args = append(args, q)
			}
		}
		db = db.Where(query[0], args...)
	}

	db.Find(entities)
	return entities, nil
}

// Get total number of entity by query condition
func (this *Repository) Count(query ...interface{}) (int64, error) {
	db := infra.Db("read").Model(Entity{})
	if len(query) > 0 {
		var args []interface{}
		for i, q := range query {
			if i > 0 {
				args = append(args, q)
			}
		}
		db = db.Where(query[0], args...)
	}

	var count int64

	db.Count(&count)
	return count, nil
}

// Get all token
func (this *Repository) GetAllToken(userId string, query ...interface{}) (*[]Token, error) {
	tokens := &[]Token{}
	db := infra.Db("read")
	if len(query) > 0 {
		var args []interface{}
		for i, q := range query {
			if i > 0 {
				args = append(args, q)
			}
		}
		db = db.Where(query[0], args...)
	}

	db.Where("user_id = ?", userId).Find(tokens)
	return tokens, nil
}

// Get all source
func (this *Repository) GetAllSource(userId string, query ...interface{}) (*[]FromSource, error) {
	sources := &[]FromSource{}
	db := infra.Db("read")
	if len(query) > 0 {
		var args []interface{}
		for i, q := range query {
			if i > 0 {
				args = append(args, q)
			}
		}
		db = db.Where(query[0], args...)
	}

	db.Where("user_id = ?", userId).Find(sources)
	return sources, nil
}

// Get all group
func (this *Repository) GetAllGroup(fromSourceId string, query ...interface{}) (*[]InGroup, error) {
	groups := &[]InGroup{}
	db := infra.Db("read")
	if len(query) > 0 {
		var args []interface{}
		for i, q := range query {
			if i > 0 {
				args = append(args, q)
			}
		}
		db = db.Where(query[0], args...)
	}

	db.Where("user_source_id = ?", fromSourceId).Find(groups)
	return groups, nil
}

// Persistent entity
func (this *Repository) Save(entity *Entity) error {
	events := append([]domain.EventBase{}, entity.Events...)

	// save token
	for _, token := range entity.Tokens {
		deleted := false
		for _, event := range token.Events {
			if event.Action == "Revoked" {
				deleted = true
				infra.Db("write").Delete(&token)
				break
			}
		}

		if !deleted {
			infra.Db("write").Save(&token)
		}

		events = append(events, token.Events...)
	}

	// save source
	for _, source := range entity.FromSources {
		deleted := false
		for _, event := range source.Events {
			if event.Action == "Removed" {
				deleted = true

				// group
				for _, group := range source.InGroups {
					group.PushEvent("Outed")
					events = append(events, group.Events...)
					infra.Db("write").Delete(&group)
				}

				infra.Db("write").Delete(&source)
				break
			}
		}

		if !deleted {
			infra.Db("write").Save(&source)

			// group
			for _, group := range source.InGroups {
				deleted = false
				for _, event := range group.Events {
					events = append(events, event)
					if event.Action == "Outed" {
						deleted = true
						infra.Db("write").Delete(&group)
						break
					}
				}

				if !deleted {
					infra.Db("write").Save(&group)
				}
			}
		}

		events = append(events, source.Events...)
	}

	if len(entity.Events) > 0 {
		infra.Db("write").Save(entity)
	}

	this.PublishEvents(events)
	return nil
}

// Remove entity
func (this *Repository) Remove(entity *Entity) error {
	if _, err := entity.DeletedAt.Value(); err == nil {
		return nil
	}

	events := []domain.EventBase{}

	// revoke all token
	if tokens, err := this.GetAllToken(entity.Id); err != nil && len(*tokens) > 0 {
		for _, token := range *tokens {
			infra.Db("write").Delete(&token)
			token.PushEvent("Removed")
			events = append(events, token.Events...)
		}
	}

	// remove all source
	if sources, err := this.GetAllSource(entity.Id); err != nil && len(*sources) > 0 {
		for _, source := range *sources {
			// remove ingroup
			if groups, err := this.GetAllGroup(source.Id); err != nil && len(*groups) > 0 {
				for _, group := range *groups {
					infra.Db("write").Delete(&group)
					group.PushEvent("Outed")
					events = append(events, group.Events...)
				}
			}

			infra.Db("write").Delete(&source)
			source.PushEvent("Removed")
			events = append(events, source.Events...)
		}
	}

	entity.PushEvent("Removed")
	infra.Db("write").Delete(entity)
	events = append(events, entity.Events...)
	this.PublishEvents(events)
	return nil
}
