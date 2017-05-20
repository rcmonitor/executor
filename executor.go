package executor

import (
	"github.com/romana/rlog"
)

type TAggregator struct {
	Count   int
	Channel chan *TWorkerDescriptor
}

type TFramedExecutor struct {
	SW         TSLServiceWorker
	Aggregator *TAggregator
	chJobs     []*TChannelContainer
}

func (psJob *TFramedExecutor) MRun() (err error) {

	psJob.Aggregator = &TAggregator{Channel: make(chan *TWorkerDescriptor)}

	//iterating over services that should be started
	for _, sSW := range psJob.SW {
		if !sSW.Service.MIsRunning() {
			rlog.Debugf("Service '%s' is not running. Starting",
				sSW.Service.MGetName())
			err = sSW.Service.MStart()
			if err != nil {
				return err
			}
		}
		rlog.Debugf("Service '%s' started; Running the job", sSW.Service.MGetName())

		sSW.Container = &TChannelContainer{sSW.Service.MGetName(), make(chan bool)}
		go sSW.runner()
		go psJob.listener(sSW.Container)

	}

	psJob.mListenAggregator()

	rlog.Trace(6, "Job well done")

	return nil
}

func (psJob *TFramedExecutor) mListenAggregator() {
	for pswDescriptor := range psJob.Aggregator.Channel {
		rlog.Tracef(5, "Channels to keep track of: %d", psJob.Aggregator.Count)

		switch pswDescriptor.Status {
		case StatusRunning:
			rlog.Debugf("Worker for '%s' is running", pswDescriptor.Name)
		case StatusDone:
			rlog.Warnf("Worker for '%s' stopped working", pswDescriptor.Name)
		default:
			rlog.Errorf("%d is unknow status for worker '%s'",
				pswDescriptor.Status, pswDescriptor.Name)
		}

		if psJob.Aggregator.Count == 0 {
			close(psJob.Aggregator.Channel)
		}

		rlog.Trace(7, "Aggregator channel closed")
	}
}

func (psJob *TFramedExecutor) listener(psContainer *TChannelContainer) {

	psJob.Aggregator.Count++
	rlog.Tracef(5, "Amount of workers is increased to: %d by %s",
		psJob.Aggregator.Count, psContainer.Name)

	for _ = range psContainer.Channel {
		psJob.Aggregator.Channel <- &TWorkerDescriptor{psContainer.Name, StatusRunning}
	}

	psJob.Aggregator.Count--
	rlog.Tracef(5, "Amount of workers is decreased to: %d by %s",
		psJob.Aggregator.Count, psContainer.Name)

	psJob.Aggregator.Channel <- &TWorkerDescriptor{psContainer.Name, StatusDone}
}
