package main

import (
	"Go-000/Week02/service"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.GET("user/:id", func(c *gin.Context) {

		name, err := service.GetUserInfor(c.Param("id"))
		if err != nil {
			log.Printf("error occurred: %v\n", errors.Unwrap(err))
			log.Printf("stack trace:\n%+v", err)
		}
		c.String(http.StatusOK, "hello %s", name)
	})
	router.Run(":8080")
}
