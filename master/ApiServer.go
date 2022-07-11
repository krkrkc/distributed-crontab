package master

import (
	"github.com/gin-contrib/static"
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

func handleKillJob(c *gin.Context) {
	jobName := c.PostForm("name")
	err := G_jobMgr.KillJob(jobName)
	if err != nil {
		resp := &common.Response{
			Errno: -1,
			Msg:   err.Error(),
			Data:  nil,
		}
		c.JSON(http.StatusBadRequest, resp)
	} else {
		resp := &common.Response{
			Errno: 1,
			Msg:   "success",
			Data:  nil,
		}
		c.JSON(http.StatusOK, resp)
	}
}

func handleWorkerList(c *gin.Context) {
	workerArr, err := G_workerMgr.ListWorkers()
	if err != nil {
		resp := &common.Response{
			Errno: -1,
			Msg:   err.Error(),
			Data:  nil,
		}
		c.JSON(http.StatusBadRequest, resp)
	} else {
		resp := &common.Response{
			Errno: 1,
			Msg:   "success",
			Data:  workerArr,
		}
		c.JSON(http.StatusOK, resp)
	}
}

func InitApiServer() {
	engine := gin.Default()
	engine.Use(static.Serve("/", static.LocalFile(G_config.WebRoot, true)))
	engine.POST("/job/save", handleJobSave)
	engine.POST("/job/delete", handleJobDelete)
	engine.GET("/job/list", handleJobList)
	engine.POST("/job/kill", handleKillJob)
	engine.GET("/worker/list", handleWorkerList)
	//engine.Static("/static", G_config.WebRoot)

	engine.Run(":" + strconv.Itoa(G_config.ApiPort))
}
