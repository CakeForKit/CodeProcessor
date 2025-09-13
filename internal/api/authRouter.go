package api

import (
	jsonrep "codeProcessor/internal/models/jsonRep"
	"codeProcessor/internal/services/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
)

type AuthRouter struct {
	authServ auth.AuthUserServ
}

func NewAuthRouter(router *gin.RouterGroup, authServ auth.AuthUserServ) AuthRouter {
	r := AuthRouter{
		authServ: authServ,
	}
	router.POST("/register", r.Register)
	router.POST("/login", r.Login)
	return r
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создание нового пользователя в системе
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body jsonrep.UserAuth true "Данные для регистрации"
// @Success 201 "Пользователь успешно зарегистрирован"
// @Failure 400 "Неверный запрос - ошибка валидации"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /register [post]
func (r *AuthRouter) Register(c *gin.Context) {
	var req jsonrep.UserAuth
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.authServ.RegisterUser(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// Login godoc
// @Summary Вход пользователя
// @Description Аутентификация пользователя и получение токена
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body jsonrep.UserAuth true "Данные для входа"
// @Success 200 {object} map[string]string "Успешный вход"
// @Failure 400 "Неверный запрос - ошибка валидации"
// @Failure 401 "Неавторизованный доступ - неверный токен или credentials"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /login [post]
func (r *AuthRouter) Login(c *gin.Context) {
	var req jsonrep.UserAuth
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := r.authServ.LoginUser(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token.String(),
	})
}
