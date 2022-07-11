package common

const (
	JobSaveDir     = "/cron/jobs/"
	JobKillDir     = "/cron/kill/"
	JobLockDir     = "/cron/lock/"
	JobWorkerDir   = "/cron/workers/"
	JobEventSave   = 1
	JobEventDelete = 2
	JobEventKill   = 3
)
