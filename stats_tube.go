package beanstalk

type StatsTube struct {
	Name                string `json:"name"`                // is the tube's name
	CurrentJobsUrgent   int    `json:"currentJobsUrgent"`   // is the number of ready jobs with priority < 1024
	CurrentJobsReady    int    `json:"currentJobsReady"`    // is the number of jobs in the ready queue in this tube
	CurrentJobsReserved int    `json:"currentJobsReserved"` // is the number of jobs reserved by all clients in this tube
	CurrentJobsDelayed  int    `json:"currentJobsDelayed"`  // is the number of delayed jobs in this tube
	CurrentJobsBuried   int    `json:"currentJobsBuried"`   // is the number of buried jobs in this tube
	TotalJobs           int    `json:"totalJobs"`           // is the cumulative count of jobs created in this tube in the current beanstalkd process
	CurrentUsing        int    `json:"currentUsing"`        // is the number of open connections that are currently using this tube
	CurrentWaiting      int    `json:"currentWaiting"`      // is the number of open connections that have issued a reserve command while watching this tube but not yet received a response
	CurrentWatching     int    `json:"currentWatching"`     // is the number of open connections that are currently watching this tube
	Pause               int    `json:"pause"`               // is the number of seconds the tube has been paused for
	CmdDelete           int    `json:"cmdDelete"`           // is the cumulative number of delete commands for this tube
	CmdPauseTube        int    `json:"cmdPauseTube"`        // is the cumulative number of pause-tube commands for this tube
	PauseTimeLeft       int    `json:"pauseTimeLeft"`       // is the number of seconds until the tube is un-paused
}
