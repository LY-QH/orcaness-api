// interface for user domain
package user

import (
	"fmt"

	"orcaness.com/api/app/anti"
	"orcaness.com/api/app/domain/corp"
	infra "orcaness.com/api/app/infra"
	util "orcaness.com/api/util"

	"github.com/gin-gonic/gin"
)

var (
	service     = NewService()
	corpService = corp.NewService()
	wework      = infra.NewWework()
)

func Router(router *gin.Engine) {
	group := router.Group(util.GetPackPath())
	{
		// Login
		group.Any("login/:source/:corpid", func(c *gin.Context) {
			source := c.Param("source")
			if !util.StringInArray(source, []string{"dingtalk", "wework", "feishu", "default"}) {
				c.JSON(400, "Invalid login source")
				return
			}

			if util.StringInArray(source, []string{"dingtalk", "wework", "feishu"}) {
				if c.Request.Method != "GET" {
					c.JSON(400, "Invalid login source")
					return
				}
			} else if source == "default" {
				if c.Request.Method != "POST" {
					c.JSON(400, "Invalid method")
					return
				}
			}

			switch source {
			case "wework":
				mobile, err := wework.Login(c)
				if err != nil {
					fmt.Println(err)
				}

				token, err := service.Login(mobile, source)
				if err != nil {
					fmt.Println(err)
				}

				c.JSON(200, token)
			case "dingtalk":
				c.JSON(400, "Notify not support")
				return
			case "feishu":
				c.JSON(400, "Login not support")
				return
			default:
				c.JSON(400, "Login not support")
				return
			}
		})

		group.GET("qrcode/:source/:corpid", func(c *gin.Context) {
			source := c.Param("source")
			if !util.StringInArray(source, []string{"dingtalk", "wework", "feishu"}) {
				c.JSON(400, "Invalid login source")
				return
			}

			var link string
			var err error

			switch source {
			case "wework":
				link, err = wework.LoginQrcode(c)
				if err != nil {
					c.JSON(200, err)
					return
				}
			case "dingtalk":
				c.JSON(400, "Notify not support")
				return
			case "feishu":
				c.JSON(400, "Login not support")
				return
			}

			c.JSON(200, link)
		})

		// Notify
		group.Any("notify/:source/:corpid", func(c *gin.Context) {
			source := c.Param("source")
			if !util.StringInArray(source, []string{"dingtalk", "wework", "feishu"}) {
				c.JSON(400, "Invalid notify platform")
				return
			}

			var event string
			var data interface{}
			var err error

			switch source {
			case "wework":
				event, data, err = wework.Notify(c)
				if err != nil {
					fmt.Println(err)
				}
				if event == "echo_str" {
					if fmt.Sprintf("%T", data) == "[]uint8" {
						c.Data(200, "text/plain; charset=utf-8", data.([]byte))
						return
					}
				}
			case "dingtalk":
				c.JSON(400, "Notify not support")
				return
			case "feishu":
				c.JSON(400, "Notify not support")
				return
			default:
				c.JSON(400, "Notify not support")
				return
			}

			switch event {
			case "create_department":
				corpService.AddGroup(source, data.(anti.Department))

			case "update_department":
				corpService.UpdateGroup(source, data.(anti.Department))

			case "delete_department":
				corpService.RemoveGroup(source, data.(anti.Department))

			case "create_user":
				service.CreateFromSource(source, data.(anti.User))

			case "update_user":
				service.UpdateFromSource(source, data.(anti.User))

			case "delete_user":
				service.RemoveFromSource(source, data.(anti.User))
			}
			c.JSON(200, data)

		})

		// Get detail
		group.GET("detail/:id", func(c *gin.Context) {
			vo, err := service.Detail(c.Param("id"))
			if err != nil {
				fmt.Println(err)
			}

			c.JSON(200, vo)
		})

		// Create user
		group.POST("create", func(c *gin.Context) {
			id, err := service.Create(c.PostForm("name"), c.PostForm("mobile"), c.PostForm("email"), c.PostForm("address"))
			if err != nil {
				fmt.Println(err)
			}

			c.JSON(200, id)
		})

		// Modify user
		group.POST("modify", func(c *gin.Context) {
			// service.Modify(c)
		})

		// Get user list
		group.GET("list", func(c *gin.Context) {
			// service.List(c)
		})

		// Get user count
		group.GET("count", func(c *gin.Context) {
			// service.Count(c)
		})
	}
}
