package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

/*
Watch 监控一个值
watchWithPrefix: 监控一个前缀的值
watchWithRange: 范围监控
watchWithProgressNotify: todo
https://vimsky.com/examples/detail/golang-ex-github.com.coreos.etcd.clientv3---WithProgressNotify-function.html
*/
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
func TestKV_WatchWithPrefix(t *testing.T) {
	//watch
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	//开一个协程进行增加和删除操作
	go func() {
		for i := 0; i < 3; i++ {
			cli.Put(ctx, fmt.Sprintf("/k%d", i), fmt.Sprintf("/v%d", i))
			time.Sleep(time.Second)
			cli.Delete(ctx, fmt.Sprintf("/k%d", i))
		}
	}()
	watchKey := cli.Watch(ctx, "/k", clientv3.WithPrefix())
	for resp := range watchKey {
		for _, e := range resp.Events {
			fmt.Printf("%s %q : %q\n", e.Type, e.Kv.Key, e.Kv.Value)
		}
	}
	cancel()
}
func TestKV_WatchWithRange(t *testing.T) {
	//watch
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	//开一个协程进行增加和删除操作
	go func() {
		for i := 0; i < 3; i++ {
			cli.Put(ctx, fmt.Sprintf("/k%d", i), fmt.Sprintf("/v%d", i))
			time.Sleep(time.Second)
			cli.Delete(ctx, fmt.Sprintf("/k%d", i))
		}
	}()
	//监控一个范围 [start,end) 和切片类似
	watchKey := cli.Watch(ctx, "/k0", clientv3.WithRange("/k3"))
	for resp := range watchKey {
		for _, e := range resp.Events {
			fmt.Printf("%s %q : %q\n", e.Type, e.Kv.Key, e.Kv.Value)
		}
	}
	cancel()
}
