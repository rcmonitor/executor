package executor_test

import (
	"bitbucket.org/rcmonitor/proxy-parse/parser"
	"bitbucket.org/rcmonitor/proxy-provider/provider"
	"github.com/rcmonitor/executor"
	"github.com/rcmonitor/executor/testdata"
	"github.com/romana/rlog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

func init() {
	rlog.SetConfFile("conf/log.conf")
}

type TTestSuiteExecute struct {
	suite.Suite
	Executor *executor.TFramedExecutor
}

func (suite *TTestSuiteExecute) TestFramedExecutorTestdata() {
	suite.Executor = &executor.TFramedExecutor{}

	suite.mSetTestWorkers()

	err := suite.Executor.MRun()

	suite.Require().Nil(err)

}

func (suite *TTestSuiteExecute) mSetTestWorkers() {
	pswPrime := &testdata.TPrimalWorker{
		Passes: 10,
		ExecutesInFrame: 1,
	}

	psdbmsPrime, err := provider.FGetServiceByTitleCached("test_prime_service")
	suite.Require().Nil(err)
	suite.Require().NotNil(psdbmsPrime)

	suite.Executor.SW = append(suite.Executor.SW, &executor.TServiceWorker{
		Service: psdbmsPrime,
		Worker:  pswPrime,
	})

	pswEncrypt := &testdata.TEncryptWorker{1000}

	psdbmsEncrypt, err := provider.FGetServiceByTitleCached("test_encrypt_service")
	suite.Require().Nil(err)
	suite.Require().NotNil(psdbmsEncrypt)

	suite.Executor.SW = append(suite.Executor.SW, &executor.TServiceWorker{
		Service: psdbmsEncrypt,
		Worker:  pswEncrypt,
	})

}

/*

func ExcludedTestTFramedJob_MRun_Primal(t *testing.T) {
	sWorker := &TPrimalWorker{10}
	sFactory := provider.TTitleFactory{&provider.TSDBMService{}}
	cService, err := sFactory.MGet("test_service")
	require.Nil(t, err)
	require.NotNil(t, cService)
	psdbmService, ok := cService.(*provider.TSDBMService)
	require.True(t, ok)
	require.NotNil(t, psdbmService)

	sJob := executor.TFramedExecutor{
		Service: psdbmService,
		Worker: sWorker,
	}
	sJob.MRun()

	time.Sleep(10 * time.Second)

	bRunning := psdbmService.MIsRunning()
	require.True(t, bRunning)

	time.Sleep(10 * time.Second)
	err = psdbmService.MStop()
	require.Nil(t, err)

	time.Sleep(10 * time.Second)
	bRunning = psdbmService.MIsRunning()
	require.False(t, bRunning)
}

*/

/*

func Excluded_TestTFramedJob_MRun_ParseDownload(t *testing.T) {
	psWorker := &parser.TWorkerParse{&parser.TDownload{}, parser.TParser{}}

	psdbmService, err := provider.FGetServiceByTitleCached("parser_download")

	require.Nil(t, err)
	require.NotNil(t, psdbmService)

	sJob := executor.TFramedExecutor{
		Service: psdbmService,
		Worker:  psWorker,
		//LogConf: "conf/service_log.conf",
		//LogDir: "log/",
	}
	sJob.MRun()

	//fCheckRunStop(t, psdbmService)
}

func Excluded_TestTFramedJob_MRun_ParseProxy(t *testing.T) {
	psWorker := &parser.TWorkerParse{&parser.TProxy{}, parser.TParser{}}

	psdbmService, err := provider.FGetServiceByTitleCached("parser_proxy")

	require.Nil(t, err)
	require.NotNil(t, psdbmService)

	sJob := executor.TFramedExecutor{
		Service: psdbmService,
		Worker: psWorker,
	}

	sJob.MRun()

	//fCheckRunStop(t, psdbmService)
}

*/

func ExcludedTestTFramedExecutor_MRun(t *testing.T) {
	sExecutor := executor.TFramedExecutor{}

	fSetWorkers(t, &sExecutor)

	intWorkersExpected := 2

	require.Equal(t, intWorkersExpected, len(sExecutor.SW))

	//t.Logf("Length of ServiceWorker slice is: %d", len(sExecutor.SW))
	//t.Logf("type of first is: %T", sExecutor.SW[0])
	//t.Fatalf("type of second is: %T", sExecutor.SW[1])

	sExecutor.MRun()

}

func fSetWorkers(t *testing.T, psExecutor *executor.TFramedExecutor) {
	pswParseProxy := &parser.TWorkerParse{&parser.TProxy{}, parser.TParser{}}

	psdbmsParserProxy, err := provider.FGetServiceByTitleCached("parser_proxy")
	require.Nil(t, err)
	require.NotNil(t, psdbmsParserProxy)

	psExecutor.SW = append(psExecutor.SW, &executor.TServiceWorker{
		Service: psdbmsParserProxy,
		Worker:  pswParseProxy,
	})

	pswParseDownload := &parser.TWorkerParse{&parser.TDownload{}, parser.TParser{}}

	psdbmsParserDownload, err := provider.FGetServiceByTitleCached("parser_download")
	require.Nil(t, err)
	require.NotNil(t, psdbmsParserDownload)

	psExecutor.SW = append(psExecutor.SW, &executor.TServiceWorker{
		Service: psdbmsParserDownload,
		Worker:  pswParseDownload,
	})
}

/*

func fCheckRunStop(t *testing.T, psdbmService *provider.TSDBMService) {
	time.Sleep(1 * time.Minute)

	boolRunning := psdbmService.MIsRunning()
	require.True(t, boolRunning)

	err := psdbmService.MStop()
	require.Nil(t, err)

	pstPeriod, err := psdbmService.MGetPeriod()
	require.Nil(t, err)
	require.NotNil(t, pstPeriod)

	time.Sleep(*pstPeriod)

	boolRunning = psdbmService.MIsRunning()
	require.False(t, boolRunning)
}

*/

func TestExecutorSuite(t *testing.T) {
	suite.Run(t, new(TTestSuiteExecute))
}
