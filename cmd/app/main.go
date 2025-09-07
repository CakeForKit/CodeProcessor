// @title Code Processor API
// @version 1.0
// @description API для обработки кода
// @host localhost:8000
// @BasePath /
package main

import (
	"codeProcessor/internal/api"
	"codeProcessor/internal/cnfg"
	"codeProcessor/internal/services"
	"codeProcessor/internal/storage"
	"fmt"

	_ "codeProcessor/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// ctx := context.Background()
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	// // Настройка CORS
	// engine.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"}, // Можно указать конкретные домены вместо "*"
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
	// 	ExposeHeaders:    []string{"Content-Length", "Content-Type"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))
	// engine.OPTIONS("/*any", func(c *gin.Context) {
	// 	c.AbortWithStatus(http.StatusNoContent)
	// })

	// Config
	appCnfg, err := cnfg.LoadAppConfig("./configs/", "app", "yaml")
	if err != nil {
		panic(fmt.Errorf("LoadAppConfig: %v", err))
	}

	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", appCnfg.Port))
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	apiGroup := engine.Group("/")

	taskStorage, err := storage.NewTaskStorage()
	if err != nil {
		panic(fmt.Errorf("NewTaskStorage: %v", err))
	}
	taskServ, err := services.NewTaskServ(taskStorage)
	if err != nil {
		panic(fmt.Errorf("NewTaskServ: %v", err))
	}

	tasksRouter := api.NewTasksRouter(apiGroup, taskServ)
	_ = tasksRouter

	engine.Run(fmt.Sprintf(":%d", appCnfg.Port))
}
