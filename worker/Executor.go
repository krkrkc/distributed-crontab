package worker

import (
	"fmt"
	"github.com/krkrkc/distributed-crontab/common"
	"math/rand"
	"os/exec"
	"time"
)

type Excutor struct {
}

var (
	G_executor *Excutor
)

func (executor *Excutor) ExecuteJob(info *common.JobExcutingInfo) {
	go func() {
		result := &common.JobExecuteResult{
			ExecuteInfo: info,
			OutPut:      make([]byte, 0),
		}

		jobLock := G_jobMgr.CreatJobLock(info.Job.Name)

		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		err := jobLock.TryLock()
		defer jobLock.UnLock()
		if err != nil {
			fmt.Println(err)
		} else {
			result.StartTime = time.Now()
			cmd := exec.CommandContext(info.CancleCtx, "C:\\cygwin64\\bin\\bash.exe", "-c", info.Job.Command)
			//cmd := exec.CommandContext(info.CancleCtx, "ping", "baidu.com")
			//cmd := exec.CommandContext(info.CancleCtx, info.Job.Command)
			output, err := cmd.CombinedOutput()
			result.EndTime = time.Now()

			result.OutPut = output
			result.Err = err
		}

		G_schduler.PushJobResult(result)
	}()
}

func InitExcutor() error {
	G_executor = &Excutor{}
	return nil
}
