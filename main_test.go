package main

import (
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
		DialTimeout: dialTimeout,
	})
	if err != nil {
		//处理错误
	}
	cli = cli2
}
