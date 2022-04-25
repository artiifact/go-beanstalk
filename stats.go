package beanstalk

type Stats struct {
	CurrentJobsUrgent     int     `json:"current-jobs-urgent" yaml:"current-jobs-urgent"`         // is the number of ready jobs with priority < 1024
	CurrentJobsReady      int     `json:"current-jobs-ready" yaml:"current-jobs-ready"`           // is the number of jobs in the ready queue
	CurrentJobsReserved   int     `json:"current-jobs-reserved" yaml:"current-jobs-reserved"`     // is the number of jobs reserved by all clients
	CurrentJobsDelayed    int     `json:"current-jobs-delayed" yaml:"current-jobs-delayed"`       // is the number of delayed jobs
	CurrentJobsBuried     int     `json:"current-jobs-buried" yaml:"current-jobs-buried"`         // is the number of buried jobs
	CmdPut                int     `json:"cmd-put" yaml:"cmd-put"`                                 // is the cumulative number of put commands
	CmdPeek               int     `json:"cmd-peek" yaml:"cmd-peek"`                               // is the cumulative number of peek commands
	CmdPeekReady          int     `json:"cmd-peek-ready" yaml:"cmd-peek-ready"`                   // is the cumulative number of peek-ready commands
	CmdPeekDelayed        int     `json:"cmd-peek-delayed" yaml:"cmd-peek-delayed"`               // is the cumulative number of peek-delayed commands
	CmdPeekBuried         int     `json:"cmd-peek-buried" yaml:"cmd-peek-buried"`                 // is the cumulative number of peek-buried commands
	CmdPeekReserve        int     `json:"cmd-peek-reserve" yaml:"cmd-reserve"`                    // is the cumulative number of reserve commands
	CmdPeekUse            int     `json:"cmd-peek-use" yaml:"cmd-use"`                            // is the cumulative number of use commands
	CmdWatch              int     `json:"cmd-watch" yaml:"cmd-watch"`                             // is the cumulative number of watch commands
	CmdIgnore             int     `json:"cmd-ignore" yaml:"cmd-ignore"`                           // is the cumulative number of ignore commands
	CmdDelete             int     `json:"cmd-delete" yaml:"cmd-delete"`                           // is the cumulative number of delete commands
	CmdRelease            int     `json:"cmd-release" yaml:"cmd-release"`                         // is the cumulative number of release commands
	CmdBury               int     `json:"cmd-bury" yaml:"cmd-bury"`                               // is the cumulative number of bury commands
	CmdKick               int     `json:"cmd-kick" yaml:"cmd-kick"`                               // is the cumulative number of kick commands
	CmdStats              int     `json:"cmd-stats" yaml:"cmd-stats"`                             // is the cumulative number of stats commands
	CmdStatsJob           int     `json:"cmd-stats-job" yaml:"cmd-stats-job"`                     // is the cumulative number of stats-job commands
	CmdStatsTube          int     `json:"cmd-stats-tube" yaml:"cmd-stats-tube"`                   // is the cumulative number of stats-tube commands
	CmdListTubes          int     `json:"cmd-list-tubes" yaml:"cmd-list-tubes"`                   // is the cumulative number of list-tubes commands
	CmdListTubeUsed       int     `json:"cmd-list-tube-used" yaml:"cmd-list-tube-used"`           // is the cumulative number of list-tube-used commands
	CmdListTubesWatched   int     `json:"cmd-list-tubes-watched" yaml:"cmd-list-tubes-watched"`   // is the cumulative number of list-tubes-watched commands
	CmdPauseTube          int     `json:"cmd-pause-tube" yaml:"cmd-pause-tube"`                   // is the cumulative number of pause-tube command
	JobTimeouts           int     `json:"job-timeouts" yaml:"job-timeouts"`                       // is the cumulative count of times a job has timed out
	TotalJobs             int     `json:"total-jobs" yaml:"total-jobs"`                           // is the cumulative count of jobs created
	MaxJobSize            int     `json:"max-job-size" yaml:"max-job-size"`                       // is the maximum number of bytes in a job
	CurrentTubes          int     `json:"current-tubes" yaml:"current-tubes"`                     // is the number of currently-existing tubes
	CurrentConnections    int     `json:"current-connections" yaml:"current-connections"`         // is the number of currently open connections
	CurrentProducers      int     `json:"current-producers" yaml:"current-producers"`             // is the number of open connections that have each issued at least one put command
	CurrentWorkers        int     `json:"current-workers" yaml:"current-workers"`                 // is the number of open connections that have each issued at least one reserve command
	CurrentWaiting        int     `json:"current-waiting" yaml:"current-waiting"`                 // is the number of open connections that have issued a reserve command but not yet received a response
	TotalConnections      int     `json:"total-connections" yaml:"total-connections"`             // is the cumulative count of connections
	PID                   int     `json:"pid" yaml:"pid"`                                         // is the process id of the server
	Version               string  `json:"version" yaml:"version"`                                 // is the version string of the server
	RUsageUTime           float64 `json:"rusage-utime" yaml:"rusage-utime"`                       // is the cumulative user CPU time of this process in seconds and microseconds
	RUsageSTime           float64 `json:"rusage-stime" yaml:"rusage-stime"`                       // is the cumulative system CPU time of this process in seconds and microseconds
	Uptime                int     `json:"uptime" yaml:"uptime"`                                   // is the number of seconds since this server process started running
	BinlogOldestIndex     int     `json:"binlog-oldest-index" yaml:"binlog-oldest-index"`         // is the index of the oldest binlog file needed to store the current jobs
	BinlogCurrentIndex    int     `json:"binlog-current-index" yaml:"binlog-current-index"`       // is the index of the current binlog file being written to. If binlog is not active this value will be 0
	BinlogMaxSize         int     `json:"binlog-max-size" yaml:"binlog-max-size"`                 // is the maximum size in bytes a binlog file is allowed to get before a new binlog file is opened
	BinlogRecordsWritten  int     `json:"binlog-records-written" yaml:"binlog-records-written"`   // is the cumulative number of records written to the binlog
	BinlogRecordsMigrated int     `json:"binlog-records-migrated" yaml:"binlog-records-migrated"` // is the cumulative number of records written as part of compaction
	Draining              bool    `json:"draining" yaml:"draining"`                               // is set to "true" if the server is in drain mode, "false" otherwise
	ID                    string  `json:"id" yaml:"id"`                                           // is a random id string for this server process, generated every time beanstalkd process starts
	Hostname              string  `json:"hostname" yaml:"hostname"`                               // is the hostname of the machine as determined by uname
	OS                    string  `json:"os" yaml:"os"`                                           // is the OS version as determined by uname
	Platform              string  `json:"platform" yaml:"platform"`                               // is the machine architecture as determined by uname
}
