package master

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/krkrkc/distributed-crontab/common"
	"net/http"
	"strconv"
)

type ApiServer struct {
	Engine gin.Engine
}

func handleJobSave(c *gin.Context) {
	message := c.PostForm("job")
	fmt.Println(message)
	var job common.Job
	if err := json.Unmarshal([]byte(message), &job); err != nil {
		resp := common.Response{
			Errno: -1,
			Msg:   err.Error(),
			Data:  nil,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if oldJob, err := G_jobMgr.SaveJob(&job); err != nil {
		resp := common.Response{
			Errno: -1,
			Msg:   err.Error(),
			Data:  nil,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	} else {
		oldJobMsg, _ := json.Marshal(oldJob)
		resp := common.Response{
			Errno: 1,
			Msg:   "success",
			Data:  string(oldJobMsg),
		}
		c.JSON(http.StatusOK, resp)
	}
}

var (
	G_apiServer *ApiServer
)

func InitApiServer() {
	engine := gin.Default()
	engine.POST("/job/save", handleJobSave)
	engine.Run(":" + strconv.Itoa(G_config.ApiPort))
}
