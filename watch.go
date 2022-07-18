package efind

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EWatcher struct {
	*EClient
	wChan clientv3.WatchChan
}

// newWatcher 新建一个watcher 监视值的使用，返回的c应该作为父级节点，
// 如果调用此方法的位置监视到wChan的改变，会向wChan发送
func (ec *EClient) newWatcher(ew *EWatcher, ctx context.Context) {
	ew = &EWatcher{}
	ew.EClient = ec
	ew.wChan = ew.cli.Watch(ctx, ec.serverName, clientv3.WithPrefix())
	ctx, _ = context.WithCancel(context.Background())
}

func (ew *EWatcher) PrintMessage() {
	for {
		select {
		case ch := <-ew.wChan:
			for _, event := range ch.Events {
				fmt.Println(event)
				fmt.Println(event.Kv.Value)
			}
		}
	}
}

//func (ec *EClient) resetWatcher() {
//
//}
