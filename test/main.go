package main

import (
	"flag"
	"fmt"
	"github.com/llightos/efind"
	"log"
	"time"
)

func main() {
	f := flag.String("addr", "0.0.0.0", "")
	flag.Parse()
	client := efind.NewClient(efind.Config{
		EtcdAddr: "127.0.0.1:2379",
		TTL:      5,
	})
	session, err := client.NewSession()
	if err != nil {
		log.Panicln(err)
		return
	}
	elect := session.NewElect("etcd/user1", *f)
	//fmt.Println("elect.ReturnLeader():", elect.ReturnLeader())
	go elect.LogReturnLeader()
	fmt.Println("elect.Info()", elect.Info())
	err = elect.Campaign()
	if err != nil {
		log.Panicln(err)
		return
	}

	for {
		time.Sleep(2 * time.Second)
	}
}
