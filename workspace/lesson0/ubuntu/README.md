# Ubuntu 20.4下安装 K8s开发环境的脚本

## 为Ubuntu设置root密码

```shell
sudo passwd
```



## Ubuntu 开启22和6443端口并且允许ssh连接root用户

```shell
sudo apt-get install openssh-server
sudo apt-get install ufw
sudo ufw enable
sudo ufw allow 22
sudo ufw allow 6443
sudo ufw allow 6443/tcp
sudo ufw allow 6443/udp

sudo vi /etc/ssh/sshd_config
#找到PermitRootLogin，将其后面的内容改为yes即可
sudo systemctl restart sshd
```



## 为Ubuntu换阿里源

```sh
sudo mv /etc/apt/sources.list /etc/apt/sources.list.bak  
sudo vim /etc/apt/sources.list  
```
​	用如下阿里源替换该文件的已有内容： 
```sh   
deb http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse  
deb-src http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse  
deb http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse  
deb-src http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse  
deb http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse  
deb-src http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse  
deb http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse  
deb-src http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse  
deb http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse  
deb-src http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse  

sudo apt-get update  
```

## 安装GNU  
```sh  
sudo apt install build-essential  
```

## 安装Docker
```sh  
sudo apt-get update  
sudo apt-get install ca-certificates curl gnupg lsb-release  

sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg  

sudo echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null  

sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin
```

## 修改ContainerD 所用的镜像库地址  
```sh  
sudo containerd config default > ~/config.toml  
```
然后编辑～config.toml去添加信息
```sh  
#修改镜像源
sudo sed -i "s#k8s.gcr.io/pause#registry.aliyuncs.com/google_containers/pause#g" /etc/containerd/config.toml

#镜像加速，修改config.toml，内容如下：
[plugins."io.containerd.grpc.v1.cri".registry.mirrors]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
          endpoint = ["https://je0sfs52.mirror.aliyuncs.com"]

sudo mv ~/config.toml /etc/containerd/config.toml  
sudo systemctl restart containerd  
```

## 安装rsync:
```sh
cd ~/Downloads  
sudo wget https://github.com/WayneD/rsync/archive/refs/tags/v3.2.4.tar.gz  
sudo tar -xf v3.2.4.tar.gz  
cd rsync-3.2.4  
```
安装一些工具包  
```sh  
sudo apt install -y gcc g++ gawk autoconf automake python3-cmarkgfm  
sudo apt install -y acl libacl1-dev  
sudo apt install -y attr libattr1-dev  
sudo apt install -y libxxhash-dev  
sudo apt install -y libzstd-dev  
sudo apt install -y liblz4-dev  
sudo apt install -y libssl-dev  
```
编译，安装  
```sh  
sudo ./configure  
sudo make  
sudo cp ./rsync /usr/local/bin/  
sudo cp ./rsync-ssl /usr/local/bin/  
```

## 安装jq：  
```sh  
sudo apt-get install jq  
```

## 安装pyyaml:  
```sh
sudo apt install python3-pip  
sudo pip install pyyaml 
```

## 安装etcd：  
```sh  
ETCD_VER=v3.5.4  
#到etcd的github release中下载对应版本

mkdir ~/etcd  
tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C ~/etcd --strip-components=1  
rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz  

sudo vim ~/.bashrc  
```
最后加入：export PATH="/home/<用户名>/etcd:${PATH}"  
```sh  
source ~/.bashrc  
```

## 安装golang (1.24及以上需要golang 1.18)：  
```sh  
cd ~/Downloads  
wget https://golang.google.cn/dl/go1.18.2.linux-amd64.tar.gz  
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.18.2.linux-amd64.tar.gz  

mkdir ~/go  
mkdir ~/go/src  
mkdir ~/go/bin  
sudo vim ~/.bashrc  
```
最后加入如下几行：  
```sh  
export GOPATH="/home/<用户名>/go"  
export GOBIN="/home/<用户名>/go/bin"  
export PATH="/usr/local/go/bin:$GOPATH/bin:${PATH}"  

source ~/.bashrc  

sudo chmod u+w /etc/sudoers
sudo vim /etc/sudoers

#修改完/etc/sudoers文件后
sudo chmod u-w /etc/sudoers
```
​	在secure_path一行加入如下目录：  
​		/usr/local/go/bin （这个是$GOPATH/bin目录）  
​		/home/<用户名>/etcd （这个是etcd命令所在目录）  
​		/home/<用户名>/go/bin （这个是go get安装的程序所在位置）  

### 设置golang代理：  
```sh  
go env -w GO111MODULE="on"  
go env -w GOPROXY="https://goproxy.cn,direct"  
```

## 安装CFSSL：  
```sh  
go install github.com/cloudflare/cfssl/cmd/...@latest  
```

## 下载kubernetes代码：  
```sh  
mkdir $GOPATH/src/k8s.io  && cd $GOPATH/src/k8s.io
git clone https://github.com/kubernetes/kubernetes.git  
git checkout -b kube1.24 v1.24.0  
```

## 编译启动本地单节点集群：  
```sh  
cd $GOPATH/src/k8s.io/kubernetes
编译单个组建：sudo make WHAT="cmd/kube-apiserver"  
编译所有组件：sudo make all  
启动本地单节点集群： sudo ./hack/local-up-cluster.sh  
```

## 开启本地debug功能

```sh
cd $GOPATH/src/k8s.io/kubernetes
# kubernetes go编译文件
sudo vi ./hack/lib/golang.sh
# 查找build_binaries()函数 vi语法
:/build_binaries()
```

### 找到一下bebug判断，注释，一直开启debug能力

> ```sh
> 	gogcflags="all=-trimpath=${trimroot} ${GOGCFLAGS:-}"
>     if [[ "${DBG:-}" == 1 ]]; then
>         # Debugging - disable optimizations and inlining.
>         gogcflags="${gogcflags} -N -l"
>     fi
> 
>     goldflags="all=$(kube::version::ldflags) ${GOLDFLAGS:-}"
>     if [[ "${DBG:-}" != 1 ]]; then
>         # Not debugging - disable symbols and DWARF.
>         goldflags="${goldflags} -s -w"
>     fi
> ```
>
> 注释判断，将debug直接放在下面， 再保存即可
>
> ```sh
> 	gogcflags="all=-trimpath=${trimroot} ${GOGCFLAGS:-}"
>     # if [[ "${DBG:-}" == 1 ]]; then
>     #     # Debugging - disable optimizations and inlining.
>     #     gogcflags="${gogcflags} -N -l"
>     # fi
> 	gogcflags="${gogcflags} -N -l"
>     goldflags="all=$(kube::version::ldflags) ${GOLDFLAGS:-}"
>     # if [[ "${DBG:-}" != 1 ]]; then
>     #     # Not debugging - disable symbols and DWARF.
>     #     goldflags="${goldflags} -s -w"
>     # fi
> ```