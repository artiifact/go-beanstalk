package beanstalk

type Stats struct {
	CurrentJobsUrgent     int     `json:"currentJobsUrgent"`     // is the number of ready jobs with priority < 1024
	CurrentJobsReady      int     `json:"currentJobsReady"`      // is the number of jobs in the ready queue
	CurrentJobsReserved   int     `json:"currentJobsReserved"`   // is the number of jobs reserved by all clients
	CurrentJobsDelayed    int     `json:"currentJobsDelayed"`    // is the number of delayed jobs
	CurrentJobsBuried     int     `json:"currentJobsBuried"`     // is the number of buried jobs
	CmdPut                int     `json:"cmdPut"`                // is the cumulative number of put commands
	CmdPeek               int     `json:"cmdPeek"`               // is the cumulative number of peek commands
	CmdPeekReady          int     `json:"cmdPeekReady"`          // is the cumulative number of peek-ready commands
	CmdPeekDelayed        int     `json:"cmdPeekDelayed"`        // is the cumulative number of peek-delayed commands
	CmdPeekBuried         int     `json:"cmdPeekBuried"`         // is the cumulative number of peek-buried commands
	CmdPeekReserve        int     `json:"cmdPeekReserve"`        // is the cumulative number of reserve commands
	CmdPeekUse            int     `json:"cmdPeekUse"`            // is the cumulative number of use commands
	CmdWatch              int     `json:"cmdWatch"`              // is the cumulative number of watch commands
	CmdIgnore             int     `json:"cmdIgnore"`             // is the cumulative number of ignore commands
	CmdDelete             int     `json:"cmdDelete"`             // is the cumulative number of delete commands
	CmdRelease            int     `json:"cmdRelease"`            // is the cumulative number of release commands
	CmdBury               int     `json:"cmdBury"`               // is the cumulative number of bury commands
	CmdKick               int     `json:"cmdKick"`               // is the cumulative number of kick commands
	CmdStats              int     `json:"cmdStats"`              // is the cumulative number of stats commands
	CmdStatsJob           int     `json:"cmdStatsJob"`           // is the cumulative number of stats-job commands
	CmdStatsTube          int     `json:"cmdStatsTube"`          // is the cumulative number of stats-tube commands
	CmdListTubes          int     `json:"cmdListTubes"`          // is the cumulative number of list-tubes commands
	CmdListTubeUsed       int     `json:"cmdListTubeUsed"`       // is the cumulative number of list-tube-used commands
	CmdListTubesWatched   int     `json:"cmdListTubesWatched"`   // is the cumulative number of list-tubes-watched commands
	CmdPauseTube          int     `json:"cmdPauseTube"`          // is the cumulative number of pause-tube command
	JobTimeouts           int     `json:"jobTimeouts"`           // is the cumulative count of times a job has timed out
	TotalJobs             int     `json:"totalJobs"`             // is the cumulative count of jobs created
	MaxJobSize            int     `json:"maxJobSize"`            // is the maximum number of bytes in a job
	CurrentTubes          int     `json:"currentTubes"`          // is the number of currently-existing tubes
	CurrentConnections    int     `json:"currentConnections"`    // is the number of currently open connections
	CurrentProducers      int     `json:"currentProducers"`      // is the number of open connections that have each issued at least one put command
	CurrentWorkers        int     `json:"currentWorkers"`        // is the number of open connections that have each issued at least one reserve command
	CurrentWaiting        int     `json:"currentWaiting"`        // is the number of open connections that have issued a reserve command but not yet received a response
	TotalConnections      int     `json:"totalConnections"`      // is the cumulative count of connections
	PID                   int     `json:"pid"`                   // is the process id of the server
	Version               string  `json:"version"`               // is the version string of the server
	RUsageUTime           float64 `json:"rUsageUTime"`           // is the cumulative user CPU time of this process in seconds and microseconds
	RUsageSTime           float64 `json:"rUsageSTime"`           // is the cumulative system CPU time of this process in seconds and microseconds
	Uptime                int     `json:"uptime"`                // is the number of seconds since this server process started running
	BinlogOldestIndex     int     `json:"binlogOldestIndex"`     // is the index of the oldest binlog file needed to store the current jobs
	BinlogCurrentIndex    int     `json:"binlogCurrentIndex"`    // is the index of the current binlog file being written to. If binlog is not active this value will be 0
	BinlogMaxSize         int     `json:"binlogMaxSize"`         // is the maximum size in bytes a binlog file is allowed to get before a new binlog file is opened
	BinlogRecordsWritten  int     `json:"binlogRecordsWritten"`  // is the cumulative number of records written to the binlog
	BinlogRecordsMigrated int     `json:"binlogRecordsMigrated"` // is the cumulative number of records written as part of compaction
	Draining              bool    `json:"draining"`              // is set to "true" if the server is in drain mode, "false" otherwise
	ID                    string  `json:"id"`                    // is a random id string for this server process, generated every time beanstalkd process starts
	Hostname              string  `json:"hostname"`              // is the hostname of the machine as determined by uname
	OS                    string  `json:"os"`                    // is the OS version as determined by uname
	Platform              string  `json:"platform"`              // is the machine architecture as determined by uname
}
