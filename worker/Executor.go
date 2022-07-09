package worker

import (
	"context"
	"github.com/krkrkc/distributed-crontab/common"
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

		result.StartTime = time.Now()
		cmd := exec.CommandContext(context.TODO(), "C:\\cygwin64\\bin\\bash.exe", "-c", info.Job.Command)
		output, err := cmd.CombinedOutput()
		result.EndTime = time.Now()

		result.OutPut = output
		result.Err = err

		G_schduler.PushJobResult(result)
	}()
}

func InitExcutor() error {
	G_executor = &Excutor{}
	return nil
}
