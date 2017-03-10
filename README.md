# OhMyQueue
omq(OhMyQueue) is a distributed message queue

#### 架构图
   ![image](./doc/arch.png)

    ohmq是一个基于发布订阅机制的分布式消息队列

#### Features：
* 使用gRPC进行通信,protobuf作为数据交换格式,支持TLS
* 消息具有生命周期,过期数据会被清理,在生命周期内可以多次被消费
* 使用RocksDB进行数据持久化
* 后续支持文件映射(mmap),消息直接写磁盘
* 弱一致性,后续提供强一致性选项(通过三阶段提交或模仿raft算法实现)
* 每个topic选举出一个leader对读写进行服务,同时cluster中除了leader以外的所有节点成为follower,对消息进行备份,如果leader不幸死亡马上从followers中选举出新的leader
* 支持水平扩展
* 支持故障恢复
