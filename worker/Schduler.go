package worker

import (
	"fmt"
	"github.com/krkrkc/distributed-crontab/common"
	"time"
)

type Schduler struct {
	JobEventChan     chan *common.JobEvent
	JobPlanTable     map[string]*common.JobSchdulePlan
	JobExcutingTable map[string]*common.JobExcutingInfo
	JobResultChan    chan *common.JobExecuteResult
}

var (
	G_schduler *Schduler
)

func InitSchduler() {
	G_schduler = &Schduler{
		JobEventChan:     make(chan *common.JobEvent, 1000),
		JobPlanTable:     make(map[string]*common.JobSchdulePlan),
		JobExcutingTable: make(map[string]*common.JobExcutingInfo),
		JobResultChan:    make(chan *common.JobExecuteResult, 1000),
	}
	go G_schduler.schduleLoop()
}

func (schduler *Schduler) schduleLoop() {
	schduleAfter := schduler.TrySchdule()
	schduleTimer := time.NewTimer(schduleAfter)
	for {
		select {
		case jobEvent := <-schduler.JobEventChan:
			schduler.handleJobEvent(jobEvent)
		case <-schduleTimer.C:
		case jobResult := <-schduler.JobResultChan:
			schduler.handleJobResult(jobResult)
		}
		schduleAfter = schduler.TrySchdule()
		schduleTimer.Reset(schduleAfter)
	}
}

func (schduler *Schduler) PushJobEvent(jobEvent *common.JobEvent) {
	schduler.JobEventChan <- jobEvent
}

func (schduler *Schduler) handleJobEvent(jobEvent *common.JobEvent) {
	switch jobEvent.EventType {
	case common.JobEventSave:
		jobSchdulePlan, err := common.BuildSchdulePlan(jobEvent.Job)
		if err != nil {
			return
		}
		schduler.JobPlanTable[jobEvent.Job.Name] = jobSchdulePlan
	case common.JobEventDelete:
		_, exist := schduler.JobPlanTable[jobEvent.Job.Name]
		if exist {
			delete(schduler.JobPlanTable, jobEvent.Job.Name)
		}
	}
}

func (schduler Schduler) TrySchdule() time.Duration {
	var schduleAfter time.Duration
	if len(schduler.JobPlanTable) == 0 {
		schduleAfter = 1 * time.Second
		return schduleAfter
	}
	now := time.Now()
	var nearTime *time.Time
	for _, jobPlan := range schduler.JobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			jobPlan.NextTime = jobPlan.Expr.Next(now)
			schduler.TryStartJob(jobPlan)
		}

		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}

	schduleAfter = (*nearTime).Sub(now)
	return schduleAfter
}

func (schduler *Schduler) TryStartJob(jobPlan *common.JobSchdulePlan) {
	_, exist := schduler.JobExcutingTable[jobPlan.Job.Name]
	if exist {
		fmt.Println("正在执行:", jobPlan.Job.Name)
		return
	}

	jobExcuteInfo := common.BuildExcuteInfo(jobPlan)
	schduler.JobExcutingTable[jobPlan.Job.Name] = jobExcuteInfo
	fmt.Println("start execut:", jobExcuteInfo.Job.Name)
	G_executor.ExecuteJob(jobExcuteInfo)
}

func (schduler *Schduler) PushJobResult(jobResult *common.JobExecuteResult) {
	schduler.JobResultChan <- jobResult
}

func (schduler *Schduler) handleJobResult(jobResult *common.JobExecuteResult) {
	delete(schduler.JobExcutingTable, jobResult.ExecuteInfo.Job.Name)
	fmt.Println("finish:", jobResult.ExecuteInfo.Job.Name, "start:", jobResult.StartTime, "endtime:", jobResult.EndTime, "output:", string(jobResult.OutPut))
}
