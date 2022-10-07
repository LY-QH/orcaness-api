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

// Persistent entity
func (this *Repository) Save(entity *Entity) error {
	if len(entity.Events) == 0 {
		return nil
	}

	infra.Db("write").Save(entity)
	this.PublishEvents(entity.Events)
	return nil

}

// Remove entity
func (this *Repository) Remove(entity *Entity) error {
	if _, err := entity.DeletedAt.Value(); err == nil {
		return nil
	}

	entity.PushEvent("Removed")
	infra.Db("write").Delete(entity)
	this.PublishEvents(entity.Events)
	return nil
}

// Save token
func (this *Repository) SaveToken(entity *Entity) error {
	infra.Db("write").Save(entity.Token)
	this.PublishEvents(entity.Events)
	return nil
}
