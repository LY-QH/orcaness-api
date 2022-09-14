package UserDomain

import entity "orcaness.com/api/app/domain/user/entity"

func GetInfoByUID(uid string) entity.UserBase {
	userEntity := entity.UserBase{}
	userEntity.Uid = "70B70602-FC7E-472C-876B-9FEFBE027E5E"
	return userEntity
}
