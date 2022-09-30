// interface for user domain
package user

import (
	"fmt"

	infra "orcaness.com/api/app/infra"
	util "orcaness.com/api/util"

	"github.com/gin-gonic/gin"
)

var (
	service = NewService()
	wework  = infra.NewWework()
)

func Router(router *gin.Engine) {
	group := router.Group(util.GetPackPath())
	{
		// Login
		group.Any("login/:platform", func(c *gin.Context) {
			platform := c.Param("platform")
			if util.StringInArray(platform, []string{"dingtalk", "wework", "feishu"}) {
				if c.Request.Method != "GET" {
					c.JSON(400, "Invalid login platform")
					return
				}
			} else if platform == "default" {
				if c.Request.Method != "POST" {
					c.JSON(400, "Invalid login platform")
					return
				}
			}

			switch platform {
			case "wework":
				mobile, err := wework.Login(c)
				if err != nil {
					fmt.Println(err)
				}

				token, err := service.LoginFromPlatform(mobile, platform)
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

		group.GET("login-qrcode/:platform", func(c *gin.Context) {
			link, err := wework.LoginQrcode(c)
			if err != nil {
				c.JSON(200, err)
				return
			}

			c.JSON(200, link)
		})

		// Notify
		group.POST("notify/:platform", func(c *gin.Context) {
			platform := c.Param("platform")
			if !util.StringInArray(c.Param("platform"), []string{"dingtalk", "wework", "feishu", "default"}) {
				c.JSON(400, "Invalid notify platform")
				return
			}

			switch platform {
			case "wework":
				// wework.Notify(c)
				c.String(200, "OK")
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
		})

		// Notify
		group.GET("notify/:platform", func(c *gin.Context) {
			platform := c.Param("platform")
			if !util.StringInArray(c.Param("platform"), []string{"dingtalk", "wework", "feishu", "default"}) {
				c.JSON(400, "Invalid notify platform")
				return
			}

			switch platform {
			case "wework":
				ret, err := wework.Notify(c)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("%T\n", ret)
				if fmt.Sprintf("%T", ret) == "[]uint8" {
					c.Data(200, "text/plain; charset=utf-8", ret.([]byte))
					// c.String(200, "aaa")
					return
				}
				c.JSON(200, ret)
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
			id, err := service.Create(c)
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
