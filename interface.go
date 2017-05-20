package executor

import "time"

type IServiceProvider interface {
	MShouldRun() bool
	MIsRunning() bool
	MStart() error
	MStop() error
	//seemed to be unwise
	//MGetFrame() time.Duration
	MGetPeriod() (*time.Duration, error)
	MGetName() string
}

type IWorker interface {
	MExecute() error
	//seemed not to be needed
	//or hard to be implemented
	//MGetSleepTime() time.Duration
}

type IErrorTimeFrame interface {
	MIsTiming() bool
	MGetFrameLeft() *time.Duration
}
