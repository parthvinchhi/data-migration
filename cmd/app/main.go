package main

import (
	"github.com/parthvinchhi/data-migration/pkg/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	routes.Routes(router)

	router.Run(":8080")
}
