package common

import (
	"encoding/json"
	"strings"
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
