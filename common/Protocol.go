package common

import (
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

type Response struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
	Data  any    `json:"data"`
}

type JobEvent struct {
	EventType int
	Job       *Job
}

type JobSchdulePlan struct {
	Job      *Job
	Expr     *cronexpr.Expression
	NextTime time.Time
}

type JobExcutingInfo struct {
	Job      *Job
	PlanTime time.Time
	RealTime time.Time
}

type JobExecuteResult struct {
	ExecuteInfo *JobExcutingInfo
	OutPut      []byte
	Err         error
	StartTime   time.Time
	EndTime     time.Time
}

func UnpackJob(value []byte) (*Job, error) {
	var job Job
	err := json.Unmarshal(value, &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JobSaveDir)
}

func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

func BuildSchdulePlan(job *Job) (*JobSchdulePlan, error) {
	expr, err := cronexpr.Parse(job.CronExpr)
	if err != nil {
		return nil, err
	}

	jobSchdulePlan := &JobSchdulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}
	return jobSchdulePlan, nil
}

func BuildExcuteInfo(jobSchdulePlan *JobSchdulePlan) *JobExcutingInfo {
	return &JobExcutingInfo{
		Job:      jobSchdulePlan.Job,
		PlanTime: jobSchdulePlan.NextTime,
		RealTime: time.Now(),
	}
}
