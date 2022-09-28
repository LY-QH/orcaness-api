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
			user, errcode := NewEntity("liyiqua", "13510966338", "122238938@qq.com")
			if errcode.Code != 0 {
				fmt.Println(errcode)
			}

			if err := repository.Save(user); err != nil {
				fmt.Println(err)
			}
		})

		group.POST("modify", func(c *gin.Context) {
			user, err := repository.Get("user.ccpvcha7qo0nlcsdpms0")
			if err != nil {
				fmt.Println(err)
			}

			user.SetToMale()

			if err := repository.Save(user); err != nil {
				fmt.Println(err)
			}
		})

		group.GET("list", func(c *gin.Context) {
			user, err := repository.GetAll("id = 'user.ccpvcha7qo0nlcsdpms0'")
			if err != nil {
				fmt.Println(err)
			}

			c.JSON(200, user)
		})

		group.GET("count", func(c *gin.Context) {
			total, err := repository.Count("id = 'user.ccpvcha7qo0nlcsdpms0'")
			if err != nil {
				fmt.Println(err)
			}

			c.JSON(200, total)
		})
	}
}
