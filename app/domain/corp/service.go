package corp

import (
	"errors"

	"orcaness.com/api/app/anti"
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
func (this *Service) Detail(id string) (*Entity, error) {
	return this.repository.Get(id)
}

// Create corp
func (this *Service) Create(name string) (id string, err error) {
	entity, err := NewCorp(name)
	if err != nil {
		return "", err
	}

	err = this.repository.Save(entity)
	if err != nil {
		return "", err
	}

	return entity.Id, nil
}

// Remove corp
func (this *Service) Remove(entity Entity) error {
	return this.repository.Remove(&entity)
}

// Add group
func (this *Service) AddGroup(source string, depart anti.Department) error {
	entity, err := this.repository.GetBySource(source, depart.CorpId)
	if err != nil {
		return err
	}

	if depart.ParentId != "" {
		parentGroup, err := this.repository.GetAllGroup(entity.Id, "source = ? and source_id = ?", source, depart.ParentId)
		if err != nil {
			return err
		}

		if len(*parentGroup) == 0 {
			return errors.New("Parent group: " + depart.ParentId + " not exists")
		}

		depart.ParentId = (*parentGroup)[0].Id
	}

	entity.AddGroup(depart.Name, source, depart.DeptId, depart.ParentId)
	this.repository.Save(entity)

	return nil
}

// Update group
func (this *Service) UpdateGroup(source string, depart anti.Department) error {
	entity, err := this.repository.GetBySource(source, depart.CorpId)
	if err != nil {
		return err
	}

	groups, err := this.repository.GetAllGroup(entity.Id, "source = ? and source_id = ?", source, depart.DeptId)
	if err != nil {
		return err
	}

	if len(*groups) == 0 {
		return this.AddGroup(source, depart)
	}

	group := (*groups)[0]

	if depart.ParentId != "" {
		parentGroup, err := this.repository.GetAllGroup(entity.Id, "source = ? and source_id = ?", source, depart.ParentId)
		if err != nil {
			return err
		}

		if len(*parentGroup) == 0 {
			return errors.New("Parent group: " + depart.ParentId + " not exists")
		}

		depart.ParentId = (*parentGroup)[0].Id
	}

	if this.inChildren(entity.Id, source, depart.ParentId, group.Id) {
		return errors.New("Parent id in children")
	}

	entity.UpdateGroup(&group, depart.Name, depart.ParentId)

	this.repository.Save(entity)

	return nil
}

// Remove group
func (this *Service) RemoveGroup(source string, depart anti.Department) error {
	entity, err := this.repository.GetBySource(source, depart.CorpId)
	if err != nil {
		return err
	}

	groups, err := this.repository.GetAllGroup(entity.Id, "source = ? and source_id = ?", source, depart.DeptId)
	if err != nil {
		return err
	}

	group := (*groups)[0]

	entity.RemoveGroup(&group)

	this.repository.Save(entity)

	return nil
}

func (this *Service) inChildren(corpId string, source string, parentId string, currentGroupId string) bool {
	children, err := this.repository.GetAllGroup(corpId, "source = ? and parent_id = ?", source, currentGroupId)
	if err != nil {
		return false
	}

	if len(*children) > 0 {
		for _, child := range *children {
			if child.Id == parentId {
				return true
			} else {
				if this.inChildren(corpId, source, parentId, child.Id) {
					return true
				}
			}
		}
	}

	return false
}
