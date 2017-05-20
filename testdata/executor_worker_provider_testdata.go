package testdata

import (
	"time"
	"github.com/romana/rlog"
	"math/big"
	"crypto/rand"
	"errors"
)


type TPrimalWorker struct {
	Passes int
	ExecutesInFrame int
	executes int
}

func (psWorker *TPrimalWorker) MExecute() error {
	tStart := time.Now()
	i := 0

	if psWorker.executes == psWorker.ExecutesInFrame {
		psWorker.executes ++
		strErrorMsg := "Worker 'primary' does not fit in frame; should sleep for 1minute"
		stMinute, _ := time.ParseDuration("1m")
		return &TFrameError{errors.New(strErrorMsg), true, &stMinute}
	}

	for i < psWorker.Passes {
		_, err := rand.Prime(rand.Reader, 1024)
		if err != nil {
			return err
		}
		//rlog.Tracef(10, "Prime number is: '%v' \n", p)
		i ++
	}
	tEnd := time.Now()
	tPassed := tEnd.Sub(tStart)
	rlog.Infof("%d Passes for prime took %s \n", psWorker.Passes, tPassed.String())

	psWorker.executes ++

	return nil
}

type TEncryptWorker struct {
	Passes int
}

func (psEncrypt *TEncryptWorker) MExecute() error {
	tStart := time.Now()

	for i := 0; i < psEncrypt.Passes; i ++ {

		bintMaxMsg := new(big.Int).Exp(big.NewInt(10), big.NewInt(2048), nil)
		bintMaxPow := new(big.Int).Exp(big.NewInt(10), big.NewInt(128), nil)

		bintMsg, err := rand.Int(rand.Reader, bintMaxMsg)
		bintE, err := rand.Int(rand.Reader, bintMaxPow)
		bintN, err := rand.Int(rand.Reader, bintMaxPow)
		if err != nil { return err }

		bintCipherText := new(big.Int).Exp(bintMsg, bintE, bintN)

		rlog.Tracef(10, "Got cyphertext: %d \n", bintCipherText)
	}

	tEnd := time.Now()
	rlog.Infof("%d Passes for encrypt took %s \n", psEncrypt.Passes, tEnd.Sub(tStart).String())

	return nil
}
