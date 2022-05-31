package efind

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

// EClient 是应用的中心，拥有大多数方法
type EClient struct {
	cli        *clientv3.Client
	serverName string
	*Config
}

// Config Etcd的服务地址
type Config struct {
	EtcdAddr string
	TTL      int
	//EleName string
}

// 返回的kv
type ResKV struct {
	key string
	val string
}

func NewClient(config Config) *EClient {
	var c EClient
	c.Config = &config
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.EtcdAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Println(err)
		return &c
	}

	c.cli = client
	return &c
}

// MatchAll serverName如jingdong/user,返回所有以其为
func (ec *EClient) MatchAll(serverName string) (kv []ResKV, err error) {
	res, err := ec.cli.Get(context.TODO(), serverName, clientv3.WithPrefix())
	ec.serverName = serverName
	if err != nil {
		return kv, err
	}
	kv = make([]ResKV, 1)
	//var kv []ResKV
	for _, v := range res.Kvs {
		kv = append(kv, ResKV{
			key: string(v.Key),
			val: string(v.Value),
		})
	}
	return
}

// MatchAServer 返回所有选定服务列表的第一个服务
func (ec *EClient) MatchAServer(serverName string) (kv ResKV, err error) {
	res, err := ec.cli.Get(context.TODO(), serverName, clientv3.WithPrefix())
	if err != nil {
		return kv, err
	}
	return ResKV{
		key: string(res.Kvs[0].Key),
		val: string(res.Kvs[0].Value),
	}, nil
}
