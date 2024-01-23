package main

import (
	"fmt"
	"log"
	_ "project/cmd/api/docs"
	"project/config"
	"project/di"
)

//	@title			lapify eCommerce API
//	@version		1.0
//	@description	API for ecommerce website
// @securityDefinitions.apiKey	JWT
//	@in							cookie
//	@name						Authorise
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@host			www.zogfestiv.store
//	@BasePath		/
// @schemes	http

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("error loading files using viper")
	}

	err = di.InitializeAPI(config)
	if err != nil {
		fmt.Println("error at initial setup")
	}
}
