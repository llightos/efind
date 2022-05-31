package efind

import (
	"context"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
)

type Election struct {
	ele *concurrency.Election
	//选举的目标（参与哪场选举）
	eleName string
	//选举的值（就是val）
	eleVal string
}

type NodeInfo struct {
	Name string
	//LeaderVal
	Value string
}

// NewElect 返回一个Election对象，但进行任何操作
func (es *ESession) NewElect(eleName, val string) (ele *Election) {
	election := concurrency.NewElection(es.session, eleName)
	ele = new(Election)
	ele.eleName = eleName
	ele.eleVal = val
	ele.ele = election
	return
}

// Campaign 表示开始从事竞选活动
func (ele *Election) Campaign() (err error) {
	log.Println("节点正在竞选", ele.Info())
	err = ele.ele.Campaign(context.TODO(), ele.eleVal)
	if err != nil {
		log.Panicln("参与选举失败")
		return err
	}
	log.Println("选举成功！", ele.Info())
	return nil
}

// Info 返回当前节点的信息
func (ele *Election) Info() *NodeInfo {
	nodeInfo := new(NodeInfo)
	nodeInfo.Name = ele.eleName
	nodeInfo.Value = ele.eleVal
	return nodeInfo
}

// ReturnLeader 返回当前选举目标,的领导人
func (ele *Election) ReturnLeader() (l *NodeInfo) {
	leader, err := ele.ele.Leader(context.TODO())
	if err != nil {
		return nil
	}
	l = new(NodeInfo)
	//fmt.Println(len(leader.Kvs))
	//fmt.Println(leader.Kvs)
	l.Name = string(leader.Kvs[0].Key)
	l.Value = string(leader.Kvs[0].Value)
	return
}

// LogReturnLeader 若领导人改变，log一下
func (ele *Election) LogReturnLeader() {
	observe := ele.ele.Observe(context.TODO())
	for {
		select {
		case leader := <-observe:

			log.Println("新的leader：", leader.Kvs[0])
		}
	}
}
