package common

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
