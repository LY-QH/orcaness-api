package corp

import "orcaness.com/api/app/anti"

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

	group := (*groups)[0]

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
