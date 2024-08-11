#### 目录说明

```shell
.
├── cmd                   # 存放命令行程序的目录
│   ├── agent            # 代理相关的命令
│   ├── alert-webhook    # 警报 Webhook 相关的命令
│   └── server           # 服务器相关的命令
├── deploy               # 部署脚本和配置文件
├── kube-configs         # Kubernetes 配置文件
├── local_yaml_dir       # 本地 YAML 配置文件
├── logs                 # 日志文件
├── scripts              # 各种脚本文件
├── sql                  # 数据库 SQL 脚本
└── src                  # 源代码目录
    ├── alert-webhook    # 警报 Webhook 相关代码
    │   └── cron         # 警报 Webhook 的定时任务代码
    ├── cache            # 缓存相关代码
    ├── common           # 公共模块和工具代码
    ├── config           # 配置相关代码
    ├── models           # 数据模型代码
    ├── pbms             # PBMS（Protocol Buffers Message Schema）相关代码
    ├── rpc              # 远程过程调用（RPC）相关代码
    ├── src              # 源代码目录下的 src 目录（需要移除或合并内容）
    │   ├── agent        # 代理相关代码
    │   ├── cron         # 定时任务代码
    │   └── job          # 任务相关代码
    └── web              # Web 应用相关代码
        ├── middleware   # 中间件代码
        ├── view_alertwebhook # 警报 Webhook 视图代码
        └── view_server  # 服务器视图代码
```

# cloudops

#### 介绍

运维CMDB平台

#### 软件架构
- 语言：Golang
- 数据库：MySQL
- 缓存：Redis

### 软件架构说明
 - 基础底座
 - 服务树
 - CMDB
 - 工单系统
 - 任务执行中心
 - Prometheus监控
 - k8s管理平台
 - CICD


#### 安装教程

1.  待补充
2.  开发中


#### 使用说明

1.  待补充
2.  开发中

#### 参与贡献

1.  待补充
2.  开发中


## 功能展示
### 系统管理

- 菜单管理

![image-20240801224619260](https://p.ipic.vip/xxsfov.png)

- 用户管理

![image-20240801224713351](https://p.ipic.vip/e6s1ax.png)

- 角色管理

![image-20240801225116033](https://p.ipic.vip/5vvi0v.png)

- 修改密码

![image-20240801225155661](https://p.ipic.vip/an7na2.png)

- API接口管理

![image-20240801225255675](https://p.ipic.vip/1r7asj.png)

### 服务树与CMDB

-  服务树
![img.png](https://p.ipic.vip/nzx05p.png)

更多功能正在持续开发中... 敬请期待

2. 开发中