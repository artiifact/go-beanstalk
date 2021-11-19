package beanstalk

type Stats struct {
	CurrentJobsUrgent     int     `yaml:"current-jobs-urgent"`     // is the number of ready jobs with priority < 1024
	CurrentJobsReady      int     `yaml:"current-jobs-ready"`      // is the number of jobs in the ready queue
	CurrentJobsReserved   int     `yaml:"current-jobs-reserved"`   // is the number of jobs reserved by all clients
	CurrentJobsDelayed    int     `yaml:"current-jobs-delayed"`    // is the number of delayed jobs
	CurrentJobsBuried     int     `yaml:"current-jobs-buried"`     // is the number of buried jobs
	CmdPut                int     `yaml:"cmd-put"`                 // is the cumulative number of put commands
	CmdPeek               int     `yaml:"cmd-peek"`                // is the cumulative number of peek commands
	CmdPeekReady          int     `yaml:"cmd-peek-ready"`          // is the cumulative number of peek-ready commands
	CmdPeekDelayed        int     `yaml:"cmd-peek-delayed"`        // is the cumulative number of peek-delayed commands
	CmdPeekBuried         int     `yaml:"cmd-peek-buried"`         // is the cumulative number of peek-buried commands
	CmdPeekReserve        int     `yaml:"cmd-reserve"`             // is the cumulative number of reserve commands
	CmdPeekUse            int     `yaml:"cmd-use"`                 // is the cumulative number of use commands
	CmdWatch              int     `yaml:"cmd-watch"`               // is the cumulative number of watch commands
	CmdIgnore             int     `yaml:"cmd-ignore"`              // is the cumulative number of ignore commands
	CmdDelete             int     `yaml:"cmd-delete"`              // is the cumulative number of delete commands
	CmdRelease            int     `yaml:"cmd-release"`             // is the cumulative number of release commands
	CmdBury               int     `yaml:"cmd-bury"`                // is the cumulative number of bury commands
	CmdKick               int     `yaml:"cmd-kick"`                // is the cumulative number of kick commands
	CmdStats              int     `yaml:"cmd-stats"`               // is the cumulative number of stats commands
	CmdStatsJob           int     `yaml:"cmd-stats-job"`           // is the cumulative number of stats-job commands
	CmdStatsTube          int     `yaml:"cmd-stats-tube"`          // is the cumulative number of stats-tube commands
	CmdListTubes          int     `yaml:"cmd-list-tubes"`          // is the cumulative number of list-tubes commands
	CmdListTubeUsed       int     `yaml:"cmd-list-tube-used"`      // is the cumulative number of list-tube-used commands
	CmdListTubesWatched   int     `yaml:"cmd-list-tubes-watched"`  // is the cumulative number of list-tubes-watched commands
	CmdPauseTUne          int     `yaml:"cmd-pause-tube"`          // is the cumulative number of pause-tube command
	JobTimeouts           int     `yaml:"job-timeouts"`            // is the cumulative count of times a job has timed out
	TotalJobs             int     `yaml:"total-jobs"`              // is the cumulative count of jobs created
	MaxJobSize            int     `yaml:"max-job-size"`            // is the maximum number of bytes in a job
	CurrentTubes          int     `yaml:"current-tubes"`           // is the number of currently-existing tubes
	CurrentConnections    int     `yaml:"current-connections"`     // is the number of currently open connections
	CurrentProducers      int     `yaml:"current-producers"`       // is the number of open connections that have each issued at least one put command
	CurrentWorkers        int     `yaml:"current-workers"`         // is the number of open connections that have each issued at least one reserve command
	CurrentWaiting        int     `yaml:"current-waiting"`         // is the number of open connections that have issued a reserve command but not yet received a response
	TotalConnections      int     `yaml:"total-connections"`       // is the cumulative count of connections
	PID                   int     `yaml:"pid"`                     // is the process id of the server
	Version               string  `yaml:"version"`                 // is the version string of the server
	RUsageUTime           float64 `yaml:"rusage-utime"`            // is the cumulative user CPU time of this process in seconds and microseconds
	RUsageSTime           float64 `yaml:"rusage-stime"`            // is the cumulative system CPU time of this process in seconds and microseconds
	Uptime                int     `yaml:"uptime"`                  // is the number of seconds since this server process started running
	BinlogOldestIndex     int     `yaml:"binlog-oldest-index"`     // is the index of the oldest binlog file needed to store the current jobs
	BinlogCurrentIndex    int     `yaml:"binlog-current-index"`    // is the index of the current binlog file being written to. If binlog is not active this value will be 0
	BinlogMaxSize         int     `yaml:"binlog-max-size"`         // is the maximum size in bytes a binlog file is allowed to get before a new binlog file is opened
	BinlogRecordsWritten  int     `yaml:"binlog-records-written"`  // is the cumulative number of records written to the binlog
	BinlogRecordsMigrated int     `yaml:"binlog-records-migrated"` // is the cumulative number of records written as part of compaction
	Draining              bool    `yaml:"draining"`                // is set to "true" if the server is in drain mode, "false" otherwise
	ID                    string  `yaml:"id"`                      // is a random id string for this server process, generated every time beanstalkd process starts
	Hostname              string  `yaml:"hostname"`                // is the hostname of the machine as determined by uname
	OS                    string  `yaml:"os"`                      // is the OS version as determined by uname
	Platform              string  `yaml:"platform"`                // is the machine architecture as determined by uname
}
