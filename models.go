package beanstalk

type Job struct {
	ID   int
	Data []byte
}

type StatsJob struct {
	// is the job id
	ID int `json:"id" yaml:"id"`
	// is the name of the tube that contains this job
	Tube string `json:"tube" yaml:"tube"`
	// is "ready" or "delayed" or "reserved" or "buried"
	State string `json:"state" yaml:"state"`
	// is the priority value set by the put, release, or bury commands
	Priority int `json:"priority" yaml:"pri"`
	// is the time in seconds since the put command that created this job
	Age int `json:"age" yaml:"age"`
	// is the integer number of seconds to wait before putting this job in the ready queue
	Delay int `json:"delay" yaml:"delay"`
	// time to run -- is the integer number of seconds a worker is allowed to run this job
	TTR int `json:"ttr" yaml:"ttr"`
	// is the number of seconds left until the server puts this job into the ready queue
	TimeLeft int `json:"timeLeft" yaml:"time-left"`
	// is the number of the earliest binlog file containing this job
	File int `json:"file" yaml:"file"`
	// is the number of times this job has been reserved
	Reserves int `json:"reserves" yaml:"reserves"`
	// is the number of times this job has timed out during a reservation
	Timeouts int `json:"timeouts" yaml:"timeouts"`
	// is the number of times a client has released this job from a reservation
	Releases int `json:"releases" yaml:"releases"`
	// is the number of times this job has been buried
	Buries int `json:"buries" yaml:"buries"`
	// is the number of times this job has been kicked
	Kicks int `json:"kicks" yaml:"kicks"`
}

type StatsTube struct {
	// is the tube's name
	Name string `json:"name" yaml:"name"`
	// is the number of ready jobs with priority < 1024
	CurrentJobsUrgent int `json:"currentJobsUrgent" yaml:"current-jobs-urgent"`
	// is the number of jobs in the ready queue in this tube
	CurrentJobsReady int `json:"currentJobsReady" yaml:"current-jobs-ready"`
	// is the number of jobs reserved by all clients in this tube
	CurrentJobsReserved int `json:"currentJobsReserved" yaml:"current-jobs-reserved"`
	// is the number of delayed jobs in this tube
	CurrentJobsDelayed int `json:"currentJobsDelayed" yaml:"current-jobs-delayed"`
	// is the number of buried jobs in this tube
	CurrentJobsBuried int `json:"currentJobsBuried" yaml:"current-jobs-buried"`
	// is the cumulative count of jobs created in this tube in the current beanstalkd process
	TotalJobs int `json:"totalJobs" yaml:"total-jobs"`
	// is the number of open connections that are currently using this tube
	CurrentUsing int `json:"currentUsing" yaml:"current-using"`
	// is the number of open connections that have issued a reserve command while watching this tube but not yet received a response
	CurrentWaiting int `json:"currentWaiting" yaml:"current-waiting"`
	// is the number of open connections that are currently watching this tube
	CurrentWatching int `json:"currentWatching" yaml:"current-watching"`
	// is the number of seconds the tube has been paused for
	Pause int `json:"pause" yaml:"pause"`
	// is the cumulative number of delete commands for this tube
	CmdDelete int `json:"cmdDelete" yaml:"cmd-delete"`
	// is the cumulative number of pause-tube commands for this tube
	CmdPauseTube int `json:"cmdPauseTube" yaml:"cmd-pause-tube"`
	// is the number of seconds until the tube is un-paused
	PauseTimeLeft int `json:"pauseTimeLeft" yaml:"pause-time-left"`
}

type Stats struct {
	// is the number of ready jobs with priority < 1024
	CurrentJobsUrgent int `json:"currentJobsUrgent" yaml:"current-jobs-urgent"`
	// is the number of jobs in the ready queue
	CurrentJobsReady int `json:"currentJobsReady" yaml:"current-jobs-ready"`
	// is the number of jobs reserved by all clients
	CurrentJobsReserved int `json:"currentJobsReserved" yaml:"current-jobs-reserved"`
	// is the number of delayed jobs
	CurrentJobsDelayed int `json:"currentJobsDelayed" yaml:"current-jobs-delayed"`
	// is the number of buried jobs
	CurrentJobsBuried int `json:"currentJobsBuried" yaml:"current-jobs-buried"`
	// is the cumulative number of put commands
	CmdPut int `json:"cmdPut" yaml:"cmd-put"`
	// is the cumulative number of peek commands
	CmdPeek int `json:"cmdPeek" yaml:"cmd-peek"`
	// is the cumulative number of peek-ready commands
	CmdPeekReady int `json:"cmdPeekReady" yaml:"cmd-peek-ready"`
	// is the cumulative number of peek-delayed commands
	CmdPeekDelayed int `json:"cmdPeekDelayed" yaml:"cmd-peek-delayed"`
	// is the cumulative number of peek-buried commands
	CmdPeekBuried int `json:"cmdPeekBuried" yaml:"cmd-peek-buried"`
	// is the cumulative number of reserve commands
	CmdReserve int `json:"cmdReserve" yaml:"cmd-reserve"`
	// is the cumulative number of use commands
	CmdUse int `json:"cmdUse" yaml:"cmd-use"`
	// is the cumulative number of watch commands
	CmdWatch int `json:"cmdWatch" yaml:"cmd-watch"`
	// is the cumulative number of ignore commands
	CmdIgnore int `json:"cmdIgnore" yaml:"cmd-ignore"`
	// is the cumulative number of delete commands
	CmdDelete int `json:"cmdDelete" yaml:"cmd-delete"`
	// is the cumulative number of release commands
	CmdRelease int `json:"cmdRelease" yaml:"cmd-release"`
	// is the cumulative number of bury commands
	CmdBury int `json:"cmdBury" yaml:"cmd-bury"`
	// is the cumulative number of kick commands
	CmdKick int `json:"cmdKick" yaml:"cmd-kick"`
	// is the cumulative number of touch commands
	CmdTouch int `json:"cmdTouch" yaml:"cmd-touch"`
	// is the cumulative number of stats commands
	CmdStats int `json:"cmdStats" yaml:"cmd-stats"`
	// is the cumulative number of stats-job commands
	CmdStatsJob int `json:"cmdStatsJob" yaml:"cmd-stats-job"`
	// is the cumulative number of stats-tube commands
	CmdStatsTube int `json:"cmdStatsTube" yaml:"cmd-stats-tube"`
	// is the cumulative number of list-tubes commands
	CmdListTubes int `json:"cmdListTubes" yaml:"cmd-list-tubes"`
	// is the cumulative number of list-tube-used commands
	CmdListTubeUsed int `json:"cmdListTubeUsed" yaml:"cmd-list-tube-used"`
	// is the cumulative number of list-tubes-watched commands
	CmdListTubesWatched int `json:"cmdListTubesWatched" yaml:"cmd-list-tubes-watched"`
	// is the cumulative number of pause-tube command
	CmdPauseTube int `json:"cmdPauseTube" yaml:"cmd-pause-tube"`
	// is the cumulative count of times a job has timed out
	JobTimeouts int `json:"jobTimeouts" yaml:"job-timeouts"`
	// is the cumulative count of jobs created
	TotalJobs int `json:"totalJobs" yaml:"total-jobs"`
	// is the maximum number of bytes in a job
	MaxJobSize int `json:"maxJobSize" yaml:"max-job-size"`
	// is the number of currently-existing tubes
	CurrentTubes int `json:"currentTubes" yaml:"current-tubes"`
	// is the number of currently open connections
	CurrentConnections int `json:"currentConnections" yaml:"current-connections"`
	// is the number of open connections that have each issued at least one put command
	CurrentProducers int `json:"currentProducers" yaml:"current-producers"`
	// is the number of open connections that have each issued at least one reserve command
	CurrentWorkers int `json:"currentWorkers" yaml:"current-workers"`
	// is the number of open connections that have issued a reserve command but not yet received a response
	CurrentWaiting int `json:"currentWaiting" yaml:"current-waiting"`
	// is the cumulative count of connections
	TotalConnections int `json:"totalConnections" yaml:"total-connections"`
	// is the process id of the server
	PID int `json:"pid" yaml:"pid"`
	// is the version string of the server
	Version string `json:"version" yaml:"version"`
	// is the cumulative user CPU time of this process in seconds and microseconds
	RUsageUTime float64 `json:"rUsageUTime" yaml:"rusage-utime"`
	// is the cumulative system CPU time of this process in seconds and microseconds
	RUsageSTime float64 `json:"rUsageSTime" yaml:"rusage-stime"`
	// is the number of seconds since this server process started running
	Uptime int `json:"uptime" yaml:"uptime"`
	// is the index of the oldest binlog file needed to store the current jobs
	BinlogOldestIndex int `json:"binlogOldestIndex" yaml:"binlog-oldest-index"`
	// is the index of the current binlog file being written to. If binlog is not active this value will be 0
	BinlogCurrentIndex int `json:"binlogCurrentIndex" yaml:"binlog-current-index"`
	// is the maximum size in bytes a binlog file is allowed to get before a new binlog file is opened
	BinlogMaxSize int `json:"binlogMaxSize" yaml:"binlog-max-size"`
	// is the cumulative number of records written to the binlog
	BinlogRecordsWritten int `json:"binlogRecordsWritten" yaml:"binlog-records-written"`
	// is the cumulative number of records written as part of compaction
	BinlogRecordsMigrated int `json:"binlogRecordsMigrated" yaml:"binlog-records-migrated"`
	// is set to "true" if the server is in drain mode, "false" otherwise
	Draining bool `json:"draining" yaml:"draining"`
	// is a random id string for this server process, generated every time beanstalkd process starts
	ID string `json:"id" yaml:"id"`
	// is the hostname of the machine as determined by uname
	Hostname string `json:"hostname" yaml:"hostname"`
	// is the OS version as determined by uname
	OS string `json:"os" yaml:"os"`
	// is the machine architecture as determined by uname
	Platform string `json:"platform" yaml:"platform"`
}
