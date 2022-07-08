package master

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/krkrkc/distributed-crontab/common"
	"net/http"
	"strconv"
)

type ApiServer struct {
	Engine gin.Engine
}

var (
	G_apiServer *ApiServer
)

func handleJobSave(c *gin.Context) {
	message := c.PostForm("job")
	//fmt.Println(message)
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

func handleJobDelete(c *gin.Context) {
	jobName := c.PostForm("name")
	oldJob, err := G_jobMgr.DeleteJob(jobName)
	if err != nil {
		resp := common.Response{
			Errno: -1,
			Msg:   err.Error(),
			Data:  nil,
		}
		c.JSON(http.StatusBadRequest, resp)
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

func handleJobList(c *gin.Context) {
	jobs, err := G_jobMgr.ListJob()
	if err != nil {
		resp := common.Response{
			Errno: -1,
			Msg:   err.Error(),
			Data:  nil,
		}
		c.JSON(http.StatusBadRequest, resp)
	} else {
		resp := common.Response{
			Errno: 1,
			Msg:   "success",
			Data:  jobs,
		}
		c.JSON(http.StatusOK, resp)
	}
}

func InitApiServer() {
	engine := gin.Default()
	engine.POST("/job/save", handleJobSave)
	engine.POST("/job/delete/", handleJobDelete)
	engine.GET("/job/list", handleJobList)
	engine.Run(":" + strconv.Itoa(G_config.ApiPort))
}
