package api

import (
	jsonrep "codeProcessor/internal/models/jsonRep"
	"codeProcessor/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommitRouter struct {
	taskServ services.TaskServ
}

func NewCommitRouter(router *gin.RouterGroup, taskServ services.TaskServ) CommitRouter {
	r := CommitRouter{
		taskServ: taskServ,
	}
	router.POST("/commit", r.CommitResult)
	return r
}

// CommitResult @Summary Commit task result
// @Description Update task status and result (used by consumers)
// @Tags Задачи
// @Accept json
// @Produce json
// @Param request body jsonrep.TaskJSON true "Task result data"
// @Success 200 {object} map[string]interface{} "Result committed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input data"
// @Failure 404 {object} map[string]interface{} "Task not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /commit [post]
func (r *CommitRouter) CommitResult(c *gin.Context) {
	var request jsonrep.TaskJSON
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.taskServ.UpdateTask(request.ID, request.Status, request.Result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
