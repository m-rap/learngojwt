package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt"
)

func CreateJwt() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte("my secret"))

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func CheckAuth(c *gin.Context) {
	tokenCookie, err := c.Cookie("token")

	if err != nil {
		fmt.Printf("error getting cookie token: %s\n", err.Error())
		c.JSON(200, gin.H{
			"isLoggedIn": false,
		})

		return
	}

	fmt.Printf("got token cookie %s\n", tokenCookie)

	token, err := jwt.Parse(tokenCookie, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Printf("wrong signing method\n")
			c.JSON(200, gin.H{
				"isLoggedIn": false,
			})
		}
		return []byte("my secret"), nil
	})

	if err != nil {
		fmt.Printf("error parsing token: %s\n", err.Error())
		c.JSON(200, gin.H{
			"isLoggedIn": false,
		})
		return
	}

	var validStr string
	if token.Valid {
		validStr = "valid"
	} else {
		validStr = "invalid"
	}
	fmt.Printf("token signature is %s\n", validStr)

	c.JSON(200, gin.H{
		"isLoggedIn": token.Valid,
	})
}

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	// username := c.PostForm("username")
	// password := c.PostForm("password")

	var l LoginParams
	c.BindJSON(&l)

	if l.Username != "rian" && l.Password != "rian" {
		fmt.Printf("wrong user/pass %s %s\n", l.Username, l.Password)
		c.JSON(200, gin.H{
			"isLoggedIn": false,
		})
		return
	}

	token, err := CreateJwt()
	if err != nil {
		fmt.Printf("create jwt fail: %s\n", err.Error())
		c.JSON(200, gin.H{
			"isLoggedIn": false,
		})
		return
	}

	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.JSON(200, gin.H{
		"isLoggedIn": true,
	})
	fmt.Println("login success")
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
}

func main() {
	fmt.Printf("Hello, world!\n")

	router := gin.Default()
	// router.GET("/", func(ctx *gin.Context) {
	// 	ctx.String(200, "Hello, world!")
	// })
	router.Static("/app", "../client")
	router.GET("/CheckAuth", CheckAuth)
	router.POST("/Login", Login)
	router.GET("/Logout", Logout)

	router.Run(":8080")
}
