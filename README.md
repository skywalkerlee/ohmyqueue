# OhMyQueue
omq(OhMyQueue) is a distributed message queue

#### 架构图
   ![image](./doc/arch.png)

    ohmq是一个基于发布订阅机制的分布式消息队列

#### DONE list：
* Pub/Sub
* clean expired data
* backend storage使用rocks
* 水平扩展

#### TODO list:
* consumer group (load balance)
* backend storage支持bitcask
* cluster manager
* 重构数据结构
