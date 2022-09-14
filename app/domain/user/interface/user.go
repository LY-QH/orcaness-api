// interface for user domain
package Interface

import (
	Service "orcaness.com/api/app/domain/user/service"
	util "orcaness.com/api/util"

	"github.com/gin-gonic/gin"
)

func Interface(router *gin.Engine) {
	group := router.Group(util.GetPackPath())
	{
		group.GET("info/:uid", func(c *gin.Context) {
			info := Service.GetInfo(c.Param("uid"))
			c.JSON(200, info)
		})
	}
}
