package efind

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
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

// ResKV 返回的kv
type ResKV struct {
	Key string
	Val string
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
			Key: string(v.Key),
			Val: string(v.Value),
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
		Key: string(res.Kvs[0].Key),
		Val: string(res.Kvs[0].Value),
	}, nil
}

//ReadConfig 实现了通过文件路径编写config文件配置,读到文件位置即可
func ReadConfig(pwd string) *EClient {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(pwd)      // path to look for the config file in
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicln(err)
	}
	serverName := viper.GetString("etcd.serverName")
	etcdAddr := viper.GetString("etcd.etcdAddr")
	ttl := viper.GetInt("etcd.ttl")
	ec := new(EClient)
	ec.Config = new(Config)
	ec.serverName = serverName
	ec.EtcdAddr = etcdAddr
	ec.TTL = ttl

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAddr},
		DialTimeout: time.Duration(int64(ttl)) * time.Second,
	})

	ec.cli = client

	fmt.Println("etcd.serverName", serverName)
	return ec
}
