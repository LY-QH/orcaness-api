package corp

import (
	"errors"

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

// Get one entity by name
func (this *Repository) GetByName(name string) (*Entity, error) {
	entity := &Entity{}
	infra.Db("read").Where("name = ?", name).First(entity)
	return entity, nil
}

// Get by source
func (this *Repository) GetBySource(source string, id string) (entity *Entity, err error) {
	entity = &Entity{}

	switch source {
	case "wework":
		infra.Db("read").Raw("select * from "+entity.TableName()+" where wework->'$.corp_id' = ?", id).Scan(&entity)
	case "dingtalk":
		infra.Db("read").Raw("select * from "+entity.TableName()+" where dingtalk->'$.corp_id' = ?", id).Scan(&entity)
	case "feishu":
		infra.Db("read").Raw("select * from "+entity.TableName()+" where feishu->'$.corp_id' = ?", id).Scan(&entity)
	default:
		err = errors.New("Invalid source")
	}

	return
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

// Get groups
func (this *Repository) GetGroup(id string, groupId string) (*Group, error) {
	group := &Group{}
	infra.Db("read").Where("id = ? and corp_id = ?", groupId, id).First(group)
	return group, nil
}

// Get Group list by query condition
func (this *Repository) GetAllGroup(id string, query ...interface{}) (*[]Group, error) {
	groups := &[]Group{}
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

	db.Where("corp_id = ?", id).Find(groups)
	return groups, nil
}

// Get total number of group by query condition
func (this *Repository) CountGroup(id string, query ...interface{}) (int64, error) {
	db := infra.Db("read").Model(Group{})
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

	db.Where("corp_id = ?", id).Count(&count)
	return count, nil
}

// Persistent entity
func (this *Repository) Save(entity *Entity) error {
	events := append([]domain.EventBase{}, entity.Events...)

	// save group
	for _, group := range entity.Groups {
		deleted := false
		for _, event := range group.Events {
			if event.Action == "Removed" {
				deleted = true
				infra.Db("write").Delete(&group)
				break
			}
		}

		if !deleted {
			infra.Db("write").Save(&group)
		}

		events = append(events, group.Events...)
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

	// remove groups
	if groups, err := this.GetAllGroup(entity.Id); err != nil && len(*groups) > 0 {
		for _, group := range *groups {
			// this.RemoveGroup(entity, &group)
			infra.Db("write").Delete(&group)
			group.PushEvent("Removed")
			events = append(events, group.Events...)
		}
	}

	entity.PushEvent("Removed")
	infra.Db("write").Delete(entity)
	events = append(events, entity.Events...)
	this.PublishEvents(events)
	return nil
}
