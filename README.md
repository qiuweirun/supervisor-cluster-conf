# supervisor-cluster-conf
一个由配置来控制web集群机器上的常驻进程的样例。由supervisor+etcd+confd实现

## 前言
使用supervisor来管理常驻进程、维护常驻进程的运行。然后启用/停止进程可通过开源的第三方web-ui实现管理，如：[supervisord-monitor](https://github.com/mlazarov/supervisord-monitor)。
这里有的问题是如果是一个集群机器下，当要变更supervisor的进程配置文件需要在每台机器上执行`supervisorctl update`命令才能生效，那么能否通过配置方式让集群的每台机器生效呢？

## 介绍
![image](https://github.com/qiuweirun/supervisor-cluster-conf/blob/main/docs/images/net.jpg)
---
通过docker构建了一个集群场景，web机器上安装`supervisord`和`confd`服务，`confd`监听`etcd`。当往`etcd`发布配置后，监听的`confd`按照预定义的模板+etcd的配置生成应用进程的`supervisor`配置文件，并且当有变更时自动执行`supervisorctl update`命令。
* etcd提供订阅发布服务
* confd作为etcd的订阅者
* `monitor.ui.com`为[supervisord-monitor](https://github.com/mlazarov/supervisord-monitor)服务，用于控制web机器进程的启用/停止
* `admin.etcd.com`为往etcd服务发布配置的后台，由[etcd-kepper](https://github.com/evildecay/etcdkeeper)组成
