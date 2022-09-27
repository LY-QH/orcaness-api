package user

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

// Persistent entity
func (this *Repository) Save(userEntity *Entity) error {
	return nil

}

// Remove entity
func (this *Repository) Remove(userEntity *Entity) error {
	return nil
}
