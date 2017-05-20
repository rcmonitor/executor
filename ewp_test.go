package executor_test

import (
	"github.com/rcmonitor/executor/testdata"
	"testing"
)

func TestEncryptWorkerProvider(t *testing.T) {
	psEncryptWorker := &testdata.TEncryptWorker{1}
	psEncryptWorker.MExecute()
}

func BenchmarkEncryptExecute(b *testing.B) {
	psEncryptWorker := &testdata.TEncryptWorker{1}
	for n := 0; n < b.N; n++ {
		psEncryptWorker.MExecute()
	}
}
