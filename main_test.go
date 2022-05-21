package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	//控制dial超时
	dialTimeout = 5 * time.Second
	//控制request超时
	requestTimeout = 10 * time.Second
)

var cli *clientv3.Client

//初始化 client 对象
//TODO: 在这里配置下你的 etcd 的地址
func init() {
	cli2, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:32379", "localhost:22379", "localhost:12379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		//处理错误
	}
	cli = cli2
}

func TestKV_Put(t *testing.T) {
	//控制超时
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := cli.Put(ctx, "/k1", "v11")
	cancel()
	if err != nil {
		log.Fatal(err)
	}
}

//获取单个值
func TestKV_Get(t *testing.T) {
	//控制超时
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, "/k1")
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	//输出结果可能有多个
	for _, kv := range resp.Kvs {
		fmt.Printf("key: %s,value:%s\n", kv.Key, kv.Value)
	}
}
func TestKV_Delete(t *testing.T) {
	//控制超时
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := cli.Delete(ctx, "/k1")
	cancel()
	if err != nil {
		log.Fatal(err)
	}
}
func TestKV_Watch(t *testing.T) {
	//watch
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	//开一个协程进行增加和删除操作
	go func() {
		for i := 0; i < 3; i++ {
			cli.Put(ctx, "/k1", "v1")
			time.Sleep(time.Second)
			cli.Delete(ctx, "/k1")
		}
	}()
	watchKey := cli.Watch(ctx, "/k1")
	for resp := range watchKey {
		for _, e := range resp.Events {
			fmt.Printf("%s %q : %q\n", e.Type, e.Kv.Key, e.Kv.Value)
		}
	}
	cancel()
}

//事务 transaction
func TestKV_Txn(t *testing.T) {
	kvc := clientv3.NewKV(cli)
	//控制超时
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := kvc.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision("/k1"), "=", 0)).
		Then(clientv3.OpPut("/k1", "不存在")).
		Else(clientv3.OpPut("/k1", "存在")).
		Commit()
	cancel()
	if err != nil {
		log.Fatal(err)
	}

	//获取结果
	gresp, err := kvc.Get(context.TODO(), "/k1")
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range gresp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}

/*
Grant：分配一个租约。
Revoke：释放一个租约。
TimeToLive：获取剩余TTL时间。
Leases：列举所有etcd中的租约。
KeepAlive：自动定时的续约某个租约。
KeepAliveOnce：为某个租约续约一次。
Close：关闭当前客户端建立的所有租约。
*/
func Test_Lease(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)

	kv := clientv3.NewKV(cli)

	//分配一个3秒的租约
	lease, err := cli.Grant(ctx, 3)
	if err != nil {
		log.Fatal(err)
	}

	//put并给到租约
	kv.Put(ctx, "/k2", "一个3秒的租约", clientv3.WithLease(lease.ID))

	//睡眠4秒让租约过期
	time.Sleep(4 * time.Second)
	resp, _ := kv.Get(ctx, "/k2")
	if resp.Count == 1 {
		fmt.Printf("%s\n", "1.没有过期")
	} else {
		fmt.Printf("%s\n", "2.过期了")
	}
}
