package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"new/handlers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	keyLength := 32 // You can adjust the length as needed
	keyBytes := make([]byte, keyLength)
	_, err := rand.Read(keyBytes)
	if err != nil {
		log.Fatal(err)
	}
	sessionKey := base64.URLEncoding.EncodeToString(keyBytes)
	store := cookie.NewStore([]byte(sessionKey))
	r.Use(sessions.Sessions("my-session", store))

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/login", handlers.ShowLoginForm)
	r.POST("/login", handlers.Login)
	r.GET("/signup", handlers.ShowSignupForm)
	r.POST("/signup", handlers.Signup)
	r.GET("/web", handlers.ShowPage)
	r.GET("/ShowPagination", handlers.ShowPagination)
	r.GET("/api/data", handlers.GetDataHandler)
	r.DELETE("/deleteUser", handlers.DeleteUser)
	r.Run(":8080")
}
