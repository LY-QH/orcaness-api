package app

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	git "gopkg.in/src-d/go-git.v4"
	UserDomain "orcaness.com/api/app/domain/user"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func CollectRoute(router *gin.Engine) {
	body404 := gin.H{
		"errcode": 404,
		"errmsg":  "Not Found",
	}

	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Header("Content-Type", "image/x-icon")
		c.String(200, "data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAABiNJREFUWEe1V2tQlGUUfr799sYue2G5uCIpqJhpY4iCOIOpGVqNmCKapEYTaqWW5a3RxsRLkVmUSpaXwDQbI5WUNEwNHYLyMigmXtJ0AVFkgWWBZS/s7te877royu6yjHVm+MO+3/s+55znPOccRjRXw8FHE7DAkF4iDH5MiPBAPpQSHhgAeqMdFQ1WlFVZcFpjgdnq85VgfAEQ3VOItHgZQmQ8nNGYcUZjwc06KxoMdnDgECBhER7IIrqnCMN7i6A32bGjuAUlN8yduuYVQEQgH+kTlKjR25B1ohlVDdZOLyQHuslZzB0pw+NqPtLz9bhS0+bxO48AZsZJkThIgsV7daj08eGHXyFA1k8OwKmbZnx5otktCLcAViUqaR4zCvTgfE+nRy/nj5YhMkSAhbkNsD10XwcAqyYocUtnxbaiFp/C7euh5GgJRkSKseCHBpdPXABMi5GiXwgfqw/pfb23S+fmjZKBBGDzA+loB9BTxeLTZBWmbdPC/h+E3ROynFeD8EmBHpfvEbMdQE5qID48rMd1rW9MJw8M6C4Aidoz/cXoGyyAVMSgyWTH33etOH7ZiO9OGaCpd72PEHPDSw5HiVEA5KI5I2R4J9c1P04vQmQsEgf50RoPU/FBBEnEZ6goESEiZZZd3IJLdyyQiXkgoY7vK4bVzmFHSQuW5zVSYE5bOV6JY5eNKP7H7ACQOSUAX51sxrVaV7QBEh7WTFDileH+EPLJU+7t10tG7D5lgMnKYVQ/MSZG+aFKZ4PRwkHhx4PFxmHk+hqaf2LEoYxJSqTtrAfjN0/D5c4JRvIWR0ic9mSoAOuTVVDLWYQH8eEncA+g1cKh6JoJfkLH7+/n6XC20tJevqFKFlOGSLHheJPL/bvTgpG2sw5M9NpqbuoQKVbmN7qWB9COmIR7YpQE2amBYHkMbHYOFQ02rCvQY2+pAQaze9YODRfSiGQebepAbKINf1W3gZm1U8s1GTkcKGv1WlJEGd94Wga1gkWogo8xmTUetb67gkVKrBQV9VbsK3V/b3xfEUiPYdb8rOPyLxhRdsviFcCbI2VY+6ISUhEPF29bEPPRHReVJHxJeEIMtYKPeoMN+8+1Ug54sh5KFosSFGA2HtdzGwu9NxoioyXvqSnDTW0cEr6ooR2REszKwdjG0VL87YrJZ2GSCBlkpQQ6AGwqbPbYcBgGKFyoRlxvEepabJiZXYfCqyYQciUNliCrsBl9gvkI9mfx583O268TIQGwaZoKzOp8HXfoohHnq9ynoL9agMUJcsyI88eULbU0BftLW5GRpMTRSyYUlBuxZKwcn7khmrdwhAWwePdZOZjXvtVyhMU/nfdMQuJheXoPHCk3wmYHthY5UjYsQoTaZhsMFo5GpSs2IlKEqDAhmKdWV3PTh0mx4oBrGT54GY8BNBlhVEC+P22Auc2OFQcbMbSXCCSUhHDkTFd6yIIxcpRWmMGI52m4fa+HIOnrWq8OfDBegeXPKylXiNdkyFj0o45+ExnCR0y4CHvOGBAXIUJ1o43OiN5sz+xgpObUOaSYTC05JUTLPY9OxNOiJWoMDBXS8iMF9vJ2LU3d4bdCcLDMCOIVefiFjXe9RiNUwYLMHbN31TsA9OsmwNujZZi/x30zcnpCZJnIdmyEiP6rzcZhWZ6ONqSViUqUVlgwdZsWja33G4+7KBA9OVjWSifo9na8ZUYgPj/W5HWAJJexPNAWnBIjxYBQIQ1F0XUTfrloRO5ZQ6c8IAK0LikAM7LrKLZ2AEQ+s1JUmLpVS5n+fxiJ1K60IKTnN+L6vc7rMpJNipIgNkKIZXmeK+JRgC0eK4fOYMc3xffnzQ5D6bLnFHQizjzm2j4f5WHybepwfwwMFWDpPkflOM3tWL50nAJyMUNb9KOmg4R9yTgFVFIeJezDY77HxWRytATTY6U0HVfvei5Pb5HppeLj4yQljpSbsOMP92O+19XMWa+krW4+2dxphTjBRATx6WpGBlASRbJHejKfllPSkGbF+6NPsIBuwGcrzbihtdJ6J4Kk9OPRsY0MGOTvdqMN239v6XTGcClDX0hG9J4QiWzBZI8gWzExsp5X6aw4V2nBhWpLl3jzLztCpq2PPzQiAAAAAElFTkSuQmCC")
	})

	// domain
	UserDomain.Router(router)

	router.GET("/test", func(c *gin.Context) {
		// Shutdown old process
		os.Exit(0)
	})

	// wework domain validate
	router.GET("/WW_verify_:code", func(c *gin.Context) {
		code := c.Param("code")
		if !strings.HasSuffix(code, ".txt") {
			c.JSON(404, body404)
			return
		}

		c.String(200, code[0:len(code)-4])
	})

	router.POST("/build", func(c *gin.Context) {
		signature := c.Request.Header.Get("x-hub-signature-256")
		if len(signature) > 7 && strings.HasPrefix(signature, "sha256=") {
			signature = signature[7:]
			token := viper.GetString("github.build_token")

			genSignature := hmac.New(sha256.New, []byte(token))
			reqBody, _ := ioutil.ReadAll(c.Request.Body)
			genSignature.Write(reqBody)
			newSignature := hex.EncodeToString(genSignature.Sum(nil))
			if signature == newSignature {
				path, _ := os.Getwd()
				// We instance a new repository targeting the given path (the .git folder)
				r, err := git.PlainOpen(path)
				if err != nil {
					fmt.Println(err)
					c.JSON(200, "ok")
					return
				}

				// Get the working directory for the repository
				w, err := r.Worktree()
				if err != nil {
					fmt.Println(err)
					c.JSON(200, "ok")
					return
				}

				// Pull the latest changes from the origin remote and merge into the current branch
				err = w.Pull(&git.PullOptions{RemoteName: "origin"})
				if err != nil {
					fmt.Println(err)
					c.JSON(200, "ok")
					return
				}

				// Build
				cmd := exec.Command(path + "/shell/build-prod.sh")
				err = cmd.Run()
				if err != nil {
					fmt.Println(err)
					c.JSON(200, "ok")
					return
				}

				// Shutdown old process
				os.Exit(0)
			}
		}

		c.JSON(200, "ok")
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, body404)
	})

	router.Use(ginBodyLogMiddleware)
}

func ginBodyLogMiddleware(c *gin.Context) {
	// blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	// c.Writer = blw
	// c.Next()
	// statusCode := c.Writer.Status()
	// if statusCode >= 400 && strings.HasPrefix(c.Request.URL.Path, "/favicon.ico") == false {
	// 	// ok this is an request with error, let's make a record for it
	// 	// now print body (or log in your preferred way)
	// 	// fmt.Println("Response body: " + blw.body.String())
	// 	Util.CollectLog(c, blw.body.String(), statusCode)
	// } else if statusCode == 200 {
	// 	Util.CollectLog(c, blw.body.String(), statusCode)
	// }
}
