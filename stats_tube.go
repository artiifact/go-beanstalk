package beanstalk

type StatsTube struct {
	Name                string `yaml:"name"`                  // is the tube's name
	CurrentJobsUrgent   int    `yaml:"current-jobs-urgent"`   // is the number of ready jobs with priority < 1024
	CurrentJobsReady    int    `yaml:"current-jobs-ready"`    // is the number of jobs in the ready queue in this tube
	CurrentJobsReserved int    `yaml:"current-jobs-reserved"` // is the number of jobs reserved by all clients in this tube
	CurrentJobsDelayed  int    `yaml:"current-jobs-delayed"`  // is the number of delayed jobs in this tube
	CurrentJobsBuried   int    `yaml:"current-jobs-buried"`   // is the number of buried jobs in this tube
	TotalJobs           int    `yaml:"total-jobs"`            // is the cumulative count of jobs created in this tube in the current beanstalkd process
	CurrentUsing        int    `yaml:"current-using"`         // is the number of open connections that are currently using this tube
	CurrentWaiting      int    `yaml:"current-waiting"`       // is the number of open connections that have issued a reserve command while watching this tube but not yet received a response
	CurrentWatching     int    `yaml:"current-watching"`      // is the number of open connections that are currently watching this tube
	Pause               int    `yaml:"pause"`                 // is the number of seconds the tube has been paused for
	CmdDelete           int    `yaml:"cmd-delete"`            // is the cumulative number of delete commands for this tube
	CmdPauseTube        int    `yaml:"cmd-pause-tube"`        // is the cumulative number of pause-tube commands for this tube
	PauseTimeLeft       int    `yaml:"pause-time-left"`       // is the number of seconds until the tube is un-paused
}
