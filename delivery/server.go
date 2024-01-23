package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"project/delivery/handlers"
	"project/delivery/routes"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files" 
)

func NewServerHttp(admin *handlers.AdminHandler, userHandler *handlers.UserHandler, orderHandler *handlers.OrderHandler) error {
	
	// engin:=gin.New()
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.UserRouter(router, userHandler)
	routes.AdminRouter(router, admin)
	routes.OrderRouter(router, orderHandler)

	router.LoadHTMLGlob("template/*.html")
	fmt.Println("Templates loaded from:", "template/*.html")
	fmt.Println("Starting server on port 8080...")

	err1 := http.ListenAndServe(":8080", router)
	if err1 != nil {
		log.Fatal(err1)
		return errors.New("can't run gin server")
	}

	return err1
}
