# supervisor-cluster-conf
一个由配置来控制web集群机器上的常驻进程的样例。由supervisor+etcd+confd实现

## 前言
使用supervisor来管理常驻进程、维护常驻进程的运行。然后启用/停止进程可通过开源的第三方web-ui实现管理，如：[supervisord-monitor](https://github.com/mlazarov/supervisord-monitor)。
这里有的问题是如果是一个集群机器下，当要变更supervisor的进程配置文件需要在每台机器上执行`supervisorctl update`命令才能生效，那么能否通过配置方式让集群的每台机器生效呢？

## 架构图
![image](https://github.com/qiuweirun/supervisor-cluster-conf/blob/main/docs/images/net.jpg)
---
通过docker构建了一个集群场景，web机器上安装`supervisord`和`confd`服务，`confd`监听`etcd`。当往`etcd`发布配置后，监听的`confd`按照预定义的模板+etcd的配置生成应用进程的`supervisor`配置文件，并且当有变更时自动执行`supervisorctl update`命令。
* etcd提供订阅发布服务
* confd作为etcd的订阅者
* `monitor.ui.com`为[supervisord-monitor](https://github.com/mlazarov/supervisord-monitor)服务，用于控制web机器进程的启用/停止
* `admin.etcd.com`为往etcd服务发布配置的后台，由[etcd-kepper](https://github.com/evildecay/etcdkeeper)组成
* `nginx-proxy`为子网访问入口，整个子网只暴露80和8080端口，通过域名的访问均通过80端口进行代理转发。

## 启动集群
* 需先安装docker、docker-compose服务，如果使用如下命令
```
docker-compose -f /path/to/docker-compose.yml up
```
* 绑定如下host，然后再通过域名访问相应的服务
```
如果docker安装方式为Docker Toolbox，使用如下即可，否则把192.168.99.100换为你本机分配到的IP
#站点1
192.168.99.100 demo.website1.com
#supervisor原生管理控制后台
192.168.99.100 svadm.website1.com
#开源管理后台
192.168.99.100 monitor.ui.com
#ETCD管理后台
192.168.99.100 admin.etcd.com
```
* `192.168.99.100:8080`为[nginxWebUI](https://github.com/cym1102/nginxWebUI)服务，配置转发域名以及负载均衡

## 效果演示
1) 访问`monitor.ui.com`，初始下三台web机只有`confd`的订阅者服务进程。
---
![image](https://github.com/qiuweirun/supervisor-cluster-conf/blob/main/docs/images/init-view.png)
---
2) 访问`demo.website1.com`，为503错误，通过下面配置让服务跑起来
3) 访问`admin.etcd.com`，发布如下配置（需修改为下图标红处的一样的地址，content为base64编码内容）。
---
![image](https://github.com/qiuweirun/supervisor-cluster-conf/blob/main/docs/images/etcd-keeper.png)
---
```
/services/apps/demo1/ips       172.16.238.11,172.16.238.12,172.16.238.13
/services/apps/demo1/content   Y29tbWFuZCA9IC9wcm9qZWN0L2RlbW9fd2Vic2l0ZQpwcm9jZXNzX25hbWU9JShwcm9ncmFtX25hbWUpcwpzdGFydHNlY3MgPSAwCm51bXByb2NzID0gMQphdXRvc3RhcnQgPSB0cnVlCmF1dG9yZXN0YXJ0ID0gdHJ1ZQ==
```
4) 发布完成后再访问`monitor.ui.com`，已经看到demo1的进程
---
![image](https://github.com/qiuweirun/supervisor-cluster-conf/blob/main/docs/images/pub-view.png)
---
5) 再次访问`demo.website1.com`看到如下页面（这是demo1站点的内容）：
---
![image](https://github.com/qiuweirun/supervisor-cluster-conf/blob/main/docs/images/demo1.png)
---
6) 试试`/services/apps/demo1/ips`删减某些IP看看会发生什么？

## 一些细节
* 集群中，通过指定生效的机器IP地址来实现，`confd`自带的模板方法没有获取服务器IP的，这是一个稍微[二次开发的confd](https://github.com/qiuweirun/confd)服务。安装启动过程可见`dockercompose/dockerfile/Dockerfile-supervisord`的`22-51`行。
* 还可以通过每台服务器维护一个唯一的hostname或者环境变量实现，有`confd`自带的模板方法支持。
* `confd`的模板为`dockercompose/dockerfile/build/confd/app.tmpl`与`etcd`的解码后的`content`组成完整的`supervisor`应用程序配置，并按照`dockercompose/dockerfile/build/confd/app.toml`定义的reload_cmd命令执行使生效。