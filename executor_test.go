package executor_test

import (
	"bitbucket.org/rcmonitor/proxy-parse/parser"
	"bitbucket.org/rcmonitor/proxy-provider/provider"
	"github.com/rcmonitor/executor"
	"github.com/rcmonitor/executor/testdata"
	"github.com/romana/rlog"
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

func (suite *TTestSuiteExecute) SetupTest() {
	suite.Executor = &executor.TFramedExecutor{}
}

func (suite *TTestSuiteExecute) ExcludedTestFramedExecutorTestdata() {

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


func (suite *TTestSuiteExecute) TestTFramedExecutorProduction() {

	suite.mSetProductionWorkers()

	intWorkersExpected := 2
	suite.Require().Equal(intWorkersExpected, len(suite.Executor.SW))

	err := suite.Executor.MRun()
	suite.Require().Nil(err)

}

func (suite *TTestSuiteExecute) mSetProductionWorkers() {
	pswParseProxy := &parser.TWorkerParse{&parser.TProxy{}, parser.TParser{}}

	psdbmsParserProxy, err := provider.FGetServiceByTitleCached("parser_proxy")
	suite.Require().Nil(err)
	suite.Require().NotNil(psdbmsParserProxy)

	suite.Executor.SW = append(suite.Executor.SW, &executor.TServiceWorker{
		Service: psdbmsParserProxy,
		Worker:  pswParseProxy,
	})

	pswParseDownload := &parser.TWorkerParse{&parser.TDownload{}, parser.TParser{}}

	psdbmsParserDownload, err := provider.FGetServiceByTitleCached("parser_download")
	suite.Require().Nil(err)
	suite.Require().NotNil(psdbmsParserDownload)

	suite.Executor.SW = append(suite.Executor.SW, &executor.TServiceWorker{
		Service: psdbmsParserDownload,
		Worker:  pswParseDownload,
	})
}


func TestExecutorSuite(t *testing.T) {
	suite.Run(t, new(TTestSuiteExecute))
}
