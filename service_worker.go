package executor

import (
	"github.com/romana/rlog"
	"time"
)

const (
	StatusRunning = iota
	StatusDone    = iota
)

type TWorkerDescriptor struct {
	Name   string
	Status int
}

type TSLServiceWorker []*TServiceWorker

type TChannelContainer struct {
	Name    string
	Channel chan bool
}

type TServiceWorker struct {
	Service   IServiceProvider
	Worker    IWorker
	Container *TChannelContainer
}

func (psSW *TServiceWorker) runner() {
	psSW.Container.Channel <- true

	//@TODO implement logger whith separate destinations for each worker
	//if !psSW.mSetServiceLog() {
	//	psSW.mTerminate()
	//	return
	//}

	rlog.Infof("started job for '%s'", psSW.Container.Name)
	for psSW.Service.MShouldRun() {

		psSW.mHandleExecutionError(psSW.Worker.MExecute())

		pstPeriod, err := psSW.Service.MGetPeriod()
		if err != nil {
			rlog.Error(err)
			psSW.mTerminate()
		}
		time.Sleep(*pstPeriod)
		rlog.Debugf("Service '%s' slept for %s, which is a periodic rest",
			psSW.Container.Name, pstPeriod.String())
	}

	rlog.Infof("Service %s shouldn't run anymore and is going to quit", psSW.Container.Name)
	psSW.mTerminate()
}

func (psSW *TServiceWorker) mTerminate() {
	err := psSW.Service.MStop()
	if err != nil {
		rlog.Errorf("Failed to gracefully terminate '%s'; exiting \n",
			psSW.Container.Name)
	}
	close(psSW.Container.Channel)
}

func (psSW *TServiceWorker) mHandleExecutionError(err error) {
	if err == nil {
		return
	}

	if sFramedError, ok := err.(IErrorTimeFrame); ok {
		if sFramedError.MIsTiming() {
			rlog.Infof("Not fitting within frame for '%s'; Sleeping for '%s' \n",
				psSW.Container.Name, sFramedError.MGetFrameLeft().String())
			ptFrameLeft := *sFramedError.MGetFrameLeft()
			time.Sleep(ptFrameLeft)
			rlog.Debugf("Service '%s' slept for %s, because otherwise it won`t fit in frame",
				psSW.Container.Name, ptFrameLeft.String())

			//	we shouldn`t hit it at all
		} else {
			rlog.Debugf("Unknown error with '%s': asserts to framed but not timing \n",
				psSW.Container.Name)

			psSW.mTerminate()

			return
		}

		//erroneous termination
	} else {
		rlog.Error(err)

		psSW.mTerminate()

		return
	}
}

/*
func (psJob *TFramedExecutor) mSetServiceLog() bool {

	if psJob.logSet {
		rlog.Debugf("New log file for service %s is already set; nothing to do",
			psJob.Service.MGetName())
		return true
	}

	var err error

	if psJob.LogConf != "" {
		rlog.SetConfFile(psJob.LogConf)
		os.Stderr.Write([]byte("rlog service configuration file set \n"))

		psJob.logSet = true
	}

	if psJob.LogDir != "" {
		psJob.fLog, err = os.Create(psJob.LogDir + psJob.Service.MGetName() + ".log")
		if err != nil {
			os.Stderr.Write([]byte("Error when trying to create log file \n"))
			return false
		}
		os.Stderr.Write([]byte("log file for service created \n"))

		rlog.SetOutput(psJob.fLog)
		rlog.Debug("Service log fiie obtained")

		psJob.logSet = true
	}


	//if psJob.fLog != nil {
	//	rlog.SetOutput(psLogFile)
	//}

	return true
}
*/
