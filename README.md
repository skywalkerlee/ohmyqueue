# OhMyQueue
omq(OhMyQueue) is a distributed message queue

#### 架构图
   ![image](./doc/arch.png)

#### Features：
    1）使用gRPC进行通信,protobuf作为数据交换格式,支持TLS(后续支持Thirft,RESTful因为相比之下使用复杂而且性能较差暂不考虑,
       毕竟是业余时间搞的小项目,时间条件不充裕(o(╯□╰)o))
    2）使用RocksDB进行数据持久化(后续支持LevelDB,性能至上^_^)
    3）后续支持文件映射(mmap),消息直接写磁盘(抄袭kafka,O(∩_∩)O)
    4）弱一致性,后续提供强一致性选项(通过raft算法实现)
    5）每个topic选举出一个leader对读写进行服务,同时cluster中除了leader以外的所有节点成为follower,对消息进行备份,
       如果leader不幸死亡马上从followers中选举新的leader
    6）暂不打算支持kafka的在topic内分partition的机制,原因:1.本人对kafka的设计理念了解还不够,2.时间原因o(╯□╰)o
    7）producer注册需要pub的topic,然后cluster进行leader选举,producer获取leader broker
    8）producer把topic内的消息推送到对应的leader broker,并且更新当前topic的状态
    9）consumer在注册需要sub的topic,获取topic对应的leader broker
    10）consumer监控所订阅topic的状态,发生改变时去对应的broker进行拉取


#### Already Done:
    1）gRPC通信框架 Done
    2）单topic发布&订阅 Done
    3）单topic的leader选举&follower备份 Done
    4）单topic的failover Done