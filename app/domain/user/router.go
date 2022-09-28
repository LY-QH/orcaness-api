// interface for user domain
package user

import (
	"fmt"

	util "orcaness.com/api/util"

	"github.com/gin-gonic/gin"
)

var (
	service = NewService()
)

func Router(router *gin.Engine) {
	group := router.Group(util.GetPackPath())
	{
		group.GET("detail/:id", func(c *gin.Context) {
			vo, err := service.Detail(c.Param("id"))
			if err != nil {
				fmt.Println(err)
			}

			c.JSON(200, vo)
		})

		group.POST("create", func(c *gin.Context) {
			id, err := service.Create(c)
			if err != nil {
				fmt.Println(err)
			}

			c.JSON(200, id)
		})

		group.POST("modify", func(c *gin.Context) {
			// service.Modify(c)
		})

		group.GET("list", func(c *gin.Context) {
			// service.List(c)
		})

		group.GET("count", func(c *gin.Context) {
			// service.Count(c)
		})
	}
}
