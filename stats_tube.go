package beanstalk

type StatsTube struct {
	Name                string `json:"name" yaml:"name"`                                 // is the tube's name
	CurrentJobsUrgent   int    `json:"currentJobsUrgent" yaml:"current-jobs-urgent"`     // is the number of ready jobs with priority < 1024
	CurrentJobsReady    int    `json:"currentJobsReady" yaml:"current-jobs-ready"`       // is the number of jobs in the ready queue in this tube
	CurrentJobsReserved int    `json:"currentJobsReserved" yaml:"current-jobs-reserved"` // is the number of jobs reserved by all clients in this tube
	CurrentJobsDelayed  int    `json:"currentJobsDelayed" yaml:"current-jobs-delayed"`   // is the number of delayed jobs in this tube
	CurrentJobsBuried   int    `json:"currentJobsBuried" yaml:"current-jobs-buried"`     // is the number of buried jobs in this tube
	TotalJobs           int    `json:"totalJobs" yaml:"total-jobs"`                      // is the cumulative count of jobs created in this tube in the current beanstalkd process
	CurrentUsing        int    `json:"currentUsing" yaml:"current-using"`                // is the number of open connections that are currently using this tube
	CurrentWaiting      int    `json:"currentWaiting" yaml:"current-waiting"`            // is the number of open connections that have issued a reserve command while watching this tube but not yet received a response
	CurrentWatching     int    `json:"currentWatching" yaml:"current-watching"`          // is the number of open connections that are currently watching this tube
	Pause               int    `json:"pause" yaml:"pause"`                               // is the number of seconds the tube has been paused for
	CmdDelete           int    `json:"cmdDelete" yaml:"cmd-delete"`                      // is the cumulative number of delete commands for this tube
	CmdPauseTube        int    `json:"cmdPauseTube" yaml:"cmd-pause-tube"`               // is the cumulative number of pause-tube commands for this tube
	PauseTimeLeft       int    `json:"pauseTimeLeft" yaml:"pause-time-left"`             // is the number of seconds until the tube is un-paused
}
