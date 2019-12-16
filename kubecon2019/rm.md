

#etcd metrics (google ectd ,maintainer) 

  raft talk   6.26  11.20
 
```
    curl  -l :2379/metrics
    
```
  ## service metrics
  has_leader      whether or not a leader exists     (import)    leader change  
  leader_changes_seen_total                  
  proposal_committed_total,applied_total,pending,failed_total  
  
  ## disk    metrics      
   1. 请求的wal -> disk
   2. snapshot ->  disk
   wal_fsync_duration_seconds  
   
   +
  ## network metrics
  
  peer_round_trip_time_seconds     1.5    10X ->ELECTION TIME
  
  ## new   metrics
  
     ### snapshot metrics
     ### peers healthiness
         1. ettcd_network_active_peers
     ### database size metrics
     ### storage layer metrics 
         etcd_servcer_sholw_aaply_total   ->   apply log  timeout  maybe disk showly
         etcd_disk_backend_defrag_duration_seconds   ->   frag 整理   disk gragbe    （read all exist kv and sort hash）  service break when frag 整理
     ### server side metrics
         etcd_server_read_indexes_failed_total
         线性读取   -->     (local read   need to ask leader to sync)
     ### learner metrics
      
     learner:  1 node -> 3 node 
               1. call API add_member
               2. starup etcd service in new nod 
               
               when add member uncorrect info , 不可逆事物 , use learner  (not vote)  etcd github/etcd/learner  pull request 10645
        
               
     
  ##  example 
  
  1. Apply entry took too long 
     (
        request too large
        slow disk: backend_commit_duration_seconds)  <10ms
        CPU starvation, memory swapping
     
     )
  2. Client request timeout
      context deadline exceeded
      Can cluster make progress
      
      
      
  ### QA    
  auditing metric ? 
  http duration ? 
  3.0  
  
       
       
       
       
 
 
 
 
 
#   commuity


#   大规模容器故障 (alibaba)  


#  containerd 

#  619 ....
   
data collector  get data 1. pod resource 2. pod predict resource  
data aggregator 
->  policy engine -> 绕开k8s 直接操作 cgroup （调整cgroup 就可以调整resource）


组件化 容器化 -> 成为k8s 的一种资源

将来：  与hpa vpa  合作 

Policy engine : (大脑 核心组件) 

api server  
command center  
  定期从data collector拿到数据
  如果node 正常  -> 遍历所有的容器 触发executor  
executor  


稳定性： 一个控制器只调整一种cgroup 资源 以免震荡  
        容器性能指标采样值： 统计数据 -> 来决定是否进行 executor
        
        
mainly control：  
  CPU
  
 

经验  
  1. 尽量微服务  
  2. 不要用不稳定接口  
  3.  QoS 资源动态调整???
  
2019/09 open source      
      
 

CRI !!!! IMPORT




# Deployment and Management in the age of Cloud Distributed Application 

##
  
  
  
  

       
       
       
       
  
      
  
  
  
  
  
  