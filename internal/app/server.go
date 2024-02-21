package app

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lud0m4n/Network/docs"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// Run запускает приложение.
func (app *Application) Run() {
	r := gin.Default()
	// r.Use(func(c *gin.Context) {
	// 	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// 	c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
	// 	c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	// 	c.Next()
	// })
	// Это нужно для автоматического создания папки "docs" в вашем проекте
	docs.SwaggerInfo.Title = "BagTracker RestAPI"
	docs.SwaggerInfo.Description = "API server for BagTracker application"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8081"
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Группа запросов для периода
	r.GET("/", app.Handler.GetPeriods)
	addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
	r.Run(addr)
	log.Println("Server down")
}
