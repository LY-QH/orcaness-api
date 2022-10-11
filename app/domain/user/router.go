// interface for user domain
package user

import (
	"encoding/json"
	"fmt"

	infra "orcaness.com/api/app/infra"
	util "orcaness.com/api/util"

	corp "orcaness.com/api/app/domain/corp"

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
		group.Any("login/:source", func(c *gin.Context) {
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

		group.GET("qrcode/:source", func(c *gin.Context) {
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
		group.Any("notify/:platform", func(c *gin.Context) {
			platform := c.Param("platform")
			if !util.StringInArray(c.Param("platform"), []string{"dingtalk", "wework", "feishu"}) {
				c.JSON(400, "Invalid notify platform")
				return
			}

			var event string
			var data interface{}
			var err error

			switch platform {
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
				// TODO: service.CreateDepart
			case "update_department":
				// TODO: service.UpdateDepart

			case "create_user":
				// TODO: service.CreateSource(data)
				var userData struct {
					Userid         string `json:"userid"`
					Name           string `json:"name"`
					Department     []uint `json:"department"`
					Position       string `json:"position"`
					Mobile         string `json:"mobile"`
					Gender         string `json:"gender"`
					Avatar         string `json:"avatar"`
					Email          string `json:"email"`
					IsLeaderInDept []uint `json:"is_leader_in_dept"`
					Status         uint   `json:"status"`
					Address        string `json:"address"`
					OpenUserid     string `json:"open_userid"`
				}
				byts, _ := json.Marshal(data)
				json.Unmarshal(byts, &userData)
				id, err := service.Create(userData.Name, userData.Mobile, userData.Email, userData.Address)
				if err != nil {
					return
				}

				entity, err := service.repository.Get(id)
				if err != nil {
					return
				}

				switch userData.Gender {
				case "male":
					entity.SetToMale()
				case "female":
					entity.SetToFemale()
				default:
					entity.HideGender()
				}

				entity.UpdateAvatar(userData.Avatar)

				// TODO: corpName
				// TODO: isSupper
				entity.AddSource("", "wework", userData.OpenUserid, 0)
				corpRespository := corp.NewRepository()
				corpEntity, _ := corpRespository.GetByName("")
				corpGroups, _ := corp.NewRepository().GetAllGroup(corpEntity.Id, "source = ? and source_id = ?", "wework", "")
				entity.FromSources[0].InGroup((*corpGroups)[0].Id, "", false)

				service.repository.Save(entity)

			case "update_user":
				// TODO: service.UpdateSource
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
