# kubernetes的安装前置要求

## 安装要求

在开始之前，部署Kubernetes集群机器需要满足以下几个条件：

- 一台或多台机器，操作系统 CentOS7.x-86_x64
- 硬件配置：4GB或更多RAM，2个CPU或更多CPU，硬盘50GB或更多
- 集群中所有机器之间网络互通
- 可以访问外网，需要拉取镜像
- 禁止swap分区

**以下操作需要在每台机器上都执行！！**

1. ​	升级centos内核：


```shell
chmod +x ./upgradeCore.sh

./upgradeCore.sh

reboot
```

2. 执行preinstall.sh脚本

```shell
chmod +x ./pre-install.sh

# <host-name>为自己想要设置的主机名,注意有的时候/etc/hosts下的主机名可能会出错，需要留意一下
./pre-install.sh <host-name> 
```

3. 配置静态ip地址

```shell
ifconfig   #查看IP地址


vim /etc/sysconfig/network-scripts/ifcfg-ens33 
TYPE="Ethernet"
PROXY_METHOD="none"
BROWSER_ONLY="no"
BOOTPROTO="static"    #这里要修改为static   
DEFROUTE="yes"
IPV4_FAILURE_FATAL="no"
IPV6INIT="yes"
IPV6_AUTOCONF="yes"
IPV6_DEFROUTE="yes"
IPV6_FAILURE_FATAL="no"
IPV6_ADDR_GEN_MODE="stable-privacy"
NAME="ens33"
UUID="cd4ce59b-0bcf-42b8-ab7d-312644bb46f3"
DEVICE="ens33"
ONBOOT="yes"
IPADDR="192.168.159.143"    #从这一行开始都是要添加的，这里添加上述查看到的ip地址
PREFIX="24"
GATEWAY="192.168.159.2"    #需要与ip地址相对应，如192.168.159
DNS1="202.119.248.66"   #我这里是我学校的DNS，你可以自己选择一个公有DNS
```

4. 重启网络服务

```shell
systemctl restart network
ping www.baidu.com
```



## 参考

[How To Configure Static IP Address In Linux And Unix](https://ostechnix.com/configure-static-ip-address-linux-unix/)
