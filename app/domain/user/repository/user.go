package Repository

import (
	Entity "orcaness.com/api/app/domain/user/entity"
)

type User struct {
}

// Query one user by uid
func (r *User) QueryUser(uid string) (entity Entity.UserBase, err error) {
	return entity, err
}

// Query user list by entity condition
func (r *User) QueryUserList(entityWhere Entity.UserBase, page uint, size uint) (entityList []Entity.UserBase, total uint, err error) {
	return entityList, total, err
}
