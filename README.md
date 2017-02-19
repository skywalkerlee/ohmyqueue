# OhMyQueue

## omq(OhMyQueue) is a distributed message queue

#### 架构图
   ![image](./doc/arch.png)

#### 目前的想法：
    1）使用gRPC进行通信，protobuf作为IDL（Interface Definition Language,接口定义语言）
    2）使用leveldb进行数据持久化
    3）broker启动时在etcd进行服务注册以及拉取每个每个broker要管理的topic（监控topic和broker的关系）
    4）producer在etcd注册需要pub的topic，并且选择当前负载最小的broker进行管理新的topic
    5）producer把topic内的消息推送到对应的broker，并且在etcd更新对应的topic的状态
    6）consumer在etcd注册需要sub的topic，获取对应的broker
    7）consumer监控所订阅topic的状态，发生改变时去对应的broker进行拉取