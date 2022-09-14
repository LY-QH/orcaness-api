package Service

import (
	"log"

	Entity "orcaness.com/api/app/domain/user/entity"
	Repository "orcaness.com/api/app/domain/user/repository"
)

// GetInfo Get user info by uid
func GetInfo(uid string) Entity.UserBase {
	userEntity, err := Repository.QueryUser("70B70602-FC7E-472C-876B-9FEFBE027E5E")
	if err != nil {
		log.Fatal(err)
	}
	return userEntity
}
