package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

/*
Grant：分配一个租约。
Revoke：释放一个租约。
TimeToLive：获取剩余TTL时间。
Leases：列举所有etcd中的租约。
KeepAlive：自动定时的续约某个租约。
KeepAliveOnce：为某个租约续约一次。
Close：关闭当前客户端建立的所有租约。
*/
func Test_Grant(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)

	kv := clientv3.NewKV(cli)

	//分配一个3秒的租约
	lease, err := cli.Grant(ctx, 3)
	if err != nil {
		log.Fatal(err)
	}

	//put并给到租约
	kv.Put(ctx, "/k2", "一个3秒的租约的数据", clientv3.WithLease(lease.ID))

	//睡眠4秒让租约过期
	time.Sleep(4 * time.Second)
	resp, _ := kv.Get(ctx, "/k2")
	if resp.Count == 1 {
		fmt.Printf("%s\n", "1.没有过期")
	} else {
		fmt.Printf("%s\n", "2.过期了")
	}
}
func Test_Revoke(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)

	kv := clientv3.NewKV(cli)

	//分配一个3秒的租约
	lease, err := cli.Grant(ctx, 3)
	if err != nil {
		log.Fatal(err)
	}

	//put并给到租约
	kv.Put(ctx, "/k2", "一个3秒的租约的数据", clientv3.WithLease(lease.ID))

	// revoking lease expires the key attached to its lease ID
	//释放一个租约,将会导致数据直接过期
	_, err = cli.Revoke(context.TODO(), lease.ID)
	if err != nil {
		log.Fatal(err)
	}
	//查询数据
	resp, _ := kv.Get(ctx, "/k2")

	if resp.Count == 1 {
		fmt.Printf("%s\n", "1.没有过期")
	} else {
		fmt.Printf("%s\n", "2.过期了")
	}
	// Output: 2.过期了
}
func Test_KeepAlive(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)

	kv := clientv3.NewKV(cli)

	//分配一个3秒的租约
	lease, err := cli.Grant(ctx, 3)
	if err != nil {
		log.Fatal(err)
	}

	//put并给到租约
	kv.Put(ctx, "/k2", "一个3秒的租约的数据", clientv3.WithLease(lease.ID))

	//永久自动续租保活KeepAlive
	ch, err := cli.KeepAlive(context.TODO(), lease.ID)
	if err != nil {
		log.Fatal(err)
	}

	ka := <-ch
	if ka != nil {
		fmt.Println("ttl:", ka.TTL)
	} else {
		fmt.Println("Unexpected NULL")
	}

	//休眠4秒看看过期没有
	time.Sleep(time.Second * 4)
	//查询数据
	resp, _ := kv.Get(ctx, "/k2")

	if resp.Count == 1 {
		fmt.Printf("%s\n", "1.没有过期")
	} else {
		fmt.Printf("%s\n", "2.过期了")
	}
	//OUTPUT: 1.没有过期
}

func Test_KeepAliveOnce(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)

	kv := clientv3.NewKV(cli)

	//分配一个3秒的租约
	lease, err := cli.Grant(ctx, 3)
	if err != nil {
		log.Fatal(err)
	}

	//put并给到租约
	kv.Put(ctx, "/k2", "一个3秒的租约的数据", clientv3.WithLease(lease.ID))

	//自动续租一次
	_, err = cli.KeepAliveOnce(context.TODO(), lease.ID)
	if err != nil {
		log.Fatal(err)
	}
	//查询数据
	resp, _ := kv.Get(ctx, "/k2")

	if resp.Count == 1 {
		fmt.Printf("%s\n", "1.没有过期")
	} else {
		fmt.Printf("%s\n", "2.过期了")
	}

	//休眠4秒看看过期没有
	time.Sleep(time.Second * 4)
	//查询数据
	resp, _ = kv.Get(ctx, "/k2")

	if resp.Count == 1 {
		fmt.Printf("%s\n", "1.没有过期")
	} else {
		fmt.Printf("%s\n", "2.过期了")
	}
	//OUTPUT: 1.没有过期
}
