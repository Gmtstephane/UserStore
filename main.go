package main

import (
	userStoreAPi "UserStore/pkg/api/v1"
	config "UserStore/pkg/config"
	"io/ioutil"
	"log"

	_ "github.com/jackc/pgconn" // swagger generator need types in main
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

func main() {

	yamlFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	c := config.GenerateConfigFromFile(yamlFile)
	UserStoreRouter := userStoreAPi.NewStore(c)
	UserStoreRouter.Serve()
}
