package user

import (
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
	entity := &Entity{Id: id}
	infra.Db().First(entity)
	return entity, nil
}

// Get entity list by query condition
func (this *Repository) GetAll(query ...interface{}) (*[]Entity, error) {
	entities := &[]Entity{}
	db := infra.Db()
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
	db := infra.Db().Model(Entity{})
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
func (this *Repository) Save(userEntity *Entity) error {
	if len(userEntity.Events) == 0 {
		return nil
	}

	infra.Db().Save(userEntity)
	this.PublishEvents(userEntity.Events)
	return nil

}

// Remove entity
func (this *Repository) Remove(userEntity *Entity) error {
	if _, err := userEntity.DeletedAt.Value(); err == nil {
		return nil
	}

	userEntity.PushEvent("Removed")
	this.PublishEvents(userEntity.Events)
	return nil
}
