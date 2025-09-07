package api

import (
	"codeProcessor/internal/models"
	"codeProcessor/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TasksRouter struct {
	taskServ services.TaskServ
}

func NewTasksRouter(router *gin.RouterGroup, taskServ services.TaskServ) TasksRouter {
	r := TasksRouter{
		taskServ: taskServ,
	}
	// gr := router.Group("task")
	router.POST("/task", r.UploadTask)
	router.GET("/status/:task_id", r.GetTaskStatus)
	router.GET("/result/:task_id", r.GetTaskResult)
	return r
}

// UploadTask godoc
// @Summary Загрузка таски на обработку
// @Description Загрузка таски на обработку, генерация uuid64 и выдача его пользователю.
// @Tags Задачи
// @Accept json
// @Produce json
// // @Param request body jsonrep.TaskAdd true "Данные задачи"
// @Success 201 "Задача успешно создана"
// @Failure 400 "Неверный запрос - ошибка валидации"
// @Router /task [post]
func (r *TasksRouter) UploadTask(c *gin.Context) {

	// var req jsonrep.TaskAdd
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	task, err := models.NewTask(uuid.New(), models.StatusInProcess, "q", "q", "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r.taskServ.AddTask(task)

	c.JSON(http.StatusCreated, gin.H{
		"task_id": task.ID(),
	})
}

// GetTaskStatus godoc
// @Summary Получение статуса таски.
// @Description Возвращает текущий статус задачи по её идентификатору
// @Tags Задачи
// @Produce json
// @Param task_id path string true "UUID идентификатор задачи" Format(uuid)
// @Success 200 "Успешное получение статуса"
// @Failure 400 "Неверный формат UUID"
// @Failure 404 "Задача не найдена"
// @Router /status/{task_id} [get]
func (r *TasksRouter) GetTaskStatus(c *gin.Context) {

	taskID, err := uuid.Parse(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID format"})
		return
	}

	stat, err := r.taskServ.GetTaskStatus(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no task with this ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": stat,
	})
}

// GetTaskResult godoc
// @Summary Получение результата задачи
// @Description Возвращает результат выполнения задачи. Если задача ещё не выполнена, возвращает ошибку.
// @Tags Задачи
// @Produce json
// @Param task_id path string true "UUID идентификатор задачи" Format(uuid)
// @Success 200 "Успешное получение результата"
// @Failure 400 "Неверный формат UUID или задача не готова"
// @Failure 404 "Задача не найдена"
// @Router /result/{task_id} [get]
func (r *TasksRouter) GetTaskResult(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID format"})
		return
	}

	stat, err := r.taskServ.GetTaskResult(taskID)
	if err != nil {
		// fmt.Print(err)
		if errors.Is(err, services.ErrTaskNotRead) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "task not ready"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "no task with this ID"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": stat,
	})
}
