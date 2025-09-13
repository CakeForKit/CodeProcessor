// @title Code Processor API
// @version 1.0
// @description API для обработки кода
// @host localhost:8000
// @BasePath /
package main

import (
	"codeProcessor/internal/api"
	"codeProcessor/internal/cnfg"
	"codeProcessor/internal/middleware"
	"codeProcessor/internal/services"
	"codeProcessor/internal/services/auth"
	"codeProcessor/internal/services/hasher"
	"codeProcessor/internal/storage"
	"fmt"

	_ "codeProcessor/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// Config
	appCnfg, err := cnfg.LoadAppConfig("./configs/", "app", "yaml")
	if err != nil {
		panic(fmt.Errorf("LoadAppConfig: %v", err))
	}

	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", appCnfg.Port))
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	storage.RegisterAllSessionStorages()

	taskStorage, err := storage.NewTaskStorage()
	if err != nil {
		panic(fmt.Errorf("NewTaskStorage: %v", err))
	}
	taskServ, err := services.NewTaskServ(taskStorage)
	if err != nil {
		panic(fmt.Errorf("NewTaskServ: %v", err))
	}

	sessionStorage := storage.SessionStorages["mem"]
	hashServ, err := hasher.NewHasher()
	if err != nil {
		panic(fmt.Errorf("NewHasher: %v", err))
	}
	userStorage, err := storage.NewUserStorage()
	if err != nil {
		panic(fmt.Errorf("NewUserStorage: %v", err))
	}
	authServ, err := auth.NewAuthUserServ(sessionStorage, hashServ, userStorage)
	if err != nil {
		panic(fmt.Errorf("NewAuthServ: %v", err))
	}

	apiGroup := engine.Group("/")
	userGroup := apiGroup.Group("/")
	userGroup.Use(middleware.AuthMiddleware(authServ))

	authRouter := api.NewAuthRouter(apiGroup, authServ)
	_ = authRouter

	tasksRouter := api.NewTasksRouter(userGroup, taskServ)
	_ = tasksRouter

	engine.Run(fmt.Sprintf(":%d", appCnfg.Port))
}
