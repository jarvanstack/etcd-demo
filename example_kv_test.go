package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	clientv3 "go.etcd.io/etcd/client/v3"
)

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

//事务 transaction
func TestKV_Txn(t *testing.T) {
	kvc := clientv3.NewKV(cli)
	//控制超时
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := kvc.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision("/k1"), "=", 0)).
		Then(clientv3.OpPut("/k1", "不存在的情况设置值")).
		Else(clientv3.OpPut("/k1", "存在的情况设置值")).
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

//watch
