// interface for user domain
package user

import (
	"fmt"

	util "orcaness.com/api/util"

	"github.com/gin-gonic/gin"
)

var (
	repository = NewRepository()
)

func Router(router *gin.Engine) {
	group := router.Group(util.GetPackPath())
	{
		group.GET("detail/:uid", func(c *gin.Context) {
			userEntity, errcode := NewEntity("andrew", "13510966337", "122238937@qq.com", "")
			if errcode.Code != 0 {
				fmt.Print(errcode)
			}
			c.JSON(200, userEntity)
		})

		group.POST("create", func(c *gin.Context) {
			user, errcode := NewEntity("liyiquan", "13510966337", "122238937@qq.com")
			if errcode.Code != 0 {
				fmt.Println(errcode)
			}

			if err := repository.Save(user); err != nil {
				fmt.Println(err)
			}
		})
	}
}
