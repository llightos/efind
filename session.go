package efind

import (
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
)

type ESession struct {
	// 此session的父类包含了v3.Client
	session *concurrency.Session
}

func (ec *EClient) NewSession() (es *ESession, err error) {
	session, err := concurrency.NewSession(ec.cli, concurrency.WithTTL(ec.TTL))
	if err != nil {
		log.Panicln(err)
		return nil, err
	}
	ess := new(ESession)
	ess.session = session
	return ess, nil
}
