package beanstalk

type StatsTube struct {
	Name                string `json:"name" yaml:"name"`                                   // is the tube's name
	CurrentJobsUrgent   int    `json:"current-jobs-urgent" yaml:"current-jobs-urgent"`     // is the number of ready jobs with priority < 1024
	CurrentJobsReady    int    `json:"current-jobs-ready" yaml:"current-jobs-ready"`       // is the number of jobs in the ready queue in this tube
	CurrentJobsReserved int    `json:"current-jobs-reserved" yaml:"current-jobs-reserved"` // is the number of jobs reserved by all clients in this tube
	CurrentJobsDelayed  int    `json:"current-jobs-delayed" yaml:"current-jobs-delayed"`   // is the number of delayed jobs in this tube
	CurrentJobsBuried   int    `json:"current-jobs-buried" yaml:"current-jobs-buried"`     // is the number of buried jobs in this tube
	TotalJobs           int    `json:"total-jobs" yaml:"total-jobs"`                       // is the cumulative count of jobs created in this tube in the current beanstalkd process
	CurrentUsing        int    `json:"current-using" yaml:"current-using"`                 // is the number of open connections that are currently using this tube
	CurrentWaiting      int    `json:"current-waiting" yaml:"current-waiting"`             // is the number of open connections that have issued a reserve command while watching this tube but not yet received a response
	CurrentWatching     int    `json:"current-watching" yaml:"current-watching"`           // is the number of open connections that are currently watching this tube
	Pause               int    `json:"pause" yaml:"pause"`                                 // is the number of seconds the tube has been paused for
	CmdDelete           int    `json:"cmd-delete" yaml:"cmd-delete"`                       // is the cumulative number of delete commands for this tube
	CmdPauseTube        int    `json:"cmd-pause-tube" yaml:"cmd-pause-tube"`               // is the cumulative number of pause-tube commands for this tube
	PauseTimeLeft       int    `json:"pause-time-left" yaml:"pause-time-left"`             // is the number of seconds until the tube is un-paused
}
