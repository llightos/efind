# efind使用简介

efind用于快速对一个etcd节点进行配置，让子节点进行选举，目前暂不支持etcd集群

```go
func New(){
//设置etcd节点名
	client := efind.NewClient(efind.Config{
		EtcdAddr: "127.0.0.1:2379",
		TTL:      5,
	})
    // 按照go.etcd.io/etcd/client/v3/concurrency包的安排，会先指定一个事件
	session, err := client.NewSession()
	if err != nil {
		log.Panicln(err)
		return
	}
    // 选举名和选举内容
	elect := session.NewElect("etcd/user1", *f)
	//fmt.Println("elect.ReturnLeader():", elect.ReturnLeader())
    // 打印当前选举目标的leader
	go elect.LogReturnLeader()
	fmt.Println("elect.Info()", elect.Info())
    // 阻塞进行选举
	err = elect.Campaign()
	if err != nil {
		log.Panicln(err)
		return
	}
    for{
        
    }
}
```

