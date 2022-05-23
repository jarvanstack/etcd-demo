# etcd-demo

etcd-demo 代码用例

## 安装

1.docker 安装

使用本项目下的 docker-compose.yaml 启动一个 3 个 etcd 节点的小集群

映射的端口分别是

```
{"localhost:32379", "localhost:22379", "localhost:12379"}
```



 docker-compose up -d

2.下载可执行文件TODO

## etcdctl todo




### client todo

* get
* put
* del
* watch
* lease



修改 main_test.go  的etcd集群配置

然后把所有 example_**_test.go 测试全部跑通没有报错就行



## 参考

官网: https://github.com/etcd-io/etcd


client: https://zhuanlan.zhihu.com/p/149805165