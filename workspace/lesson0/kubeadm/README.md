# kubeadm安装k8s

​	kubeadm是官方社区推出的一个用于快速部署kubernetes集群的工具

​	安装版本可依据个人情况，**本文使用的是1.20.0**。

​	这个工具能通过两条指令完成一个kubernetes集群的部署：

```shell
# 创建一个 Master 节点
$ kubeadm init

# 将一个 Node 节点加入到当前集群中
$ kubeadm join <Master节点的IP和端口 >
```

 

## 安装Docker【所有节点】

​	Kubernetes默认CRI（容器运行时）为Docker，因此先安装Docker（注意版本与Kubernetes版本的兼容性）。**不过在Kubernetes1.24中，已经正式弃用Docker，而使用containerd作为CRI。**

```shell
wget https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo -O /etc/yum.repos.d/docker-ce.repo
yum install -y docker-ce-19.03.5-3.el7 docker-ce-cli-19.03.5-3.el7 containerd.io
systemctl enable docker && systemctl start docker
```

​	配置Docker镜像下载加速器，并且修改Docker的cgroup driver：

```shell
cat > /etc/docker/daemon.json << EOF
{
  "registry-mirrors": ["https://tzksttqp.mirror.aliyuncs.com"]
}
EOF

vim /usr/lib/systemd/system/docker.service
#在ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock  这一行添加 --exec-opt native.cgroupdriver=systemd  参数


systemctl daemon-reload &&  systemctl restart docker


docker info | grep Cgroup   #查看docker的 cgroup driver
```



## 安装kubeadm/kubelet/kubectl

### 添加阿里云YUM软件源【所有节点】

​	为了能够在国内下载kubeadm、kubelet和kubectl

```shell
cat > /etc/yum.repos.d/kubernetes.repo << EOF
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=0
repo_gpgcheck=0
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
```

### 安装kubeadm，kubelet和kubectl【所有节点】

​	由于版本更新频繁，这里可以指定版本号部署（不指定即下载最新版本，这里使用的是1.20版本）：

```shell
yum install -y kubelet-1.20.0 kubeadm-1.20.0 kubectl-1.20.0

systemctl enable kubelet
systemctl start kubelet
```

### 部署Kubernetes Master

​	在Mster节点（192.168.159.143）执行：

```shell
kubeadm init --apiserver-advertise-address=192.168.159.143 --image-repository registry.aliyuncs.com/google_containers --kubernetes-version v1.20.0 --service-cidr=10.96.0.0/12 --pod-network-cidr=10.244.0.0/16 --ignore-preflight-errors=all
```

- --apiserver-advertise-address 集群通告地址

- --image-repository  由于默认拉取镜像地址k8s.gcr.io国内无法访问，这里指定阿里云镜像仓库地址

- --kubernetes-version K8s版本，**与上面kubeadm等安装的版本一致**

- --service-cidr 集群内部虚拟网络，可以通过 `kubectl get svc` 查看svc的CLUSTER-IP字段

- --pod-network-cidr Pod网络，可以通过 `kubectl get pods -o wide` 查看pod的IP字段。**与下面部署的CNI网络组件yaml中保持一致**

- --ignore-preflight-errors=all   忽视一些不是很重要的警告

  

**PS：可能会出现It seems like the kubelet isn't running or healthy的错误，此时可以参考这三篇博客：**

- https://blog.csdn.net/boling_cavalry/article/details/91306095
- https://blog.csdn.net/weixin_41298721/article/details/114916421

- http://www.manongjc.com/detail/23-umxjtqmublnyjwl.html


​	或者使用配置文件引导：

```shell
vi kubeadm.conf

apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
kubernetesVersion: v1.20.0
imageRepository: registry.aliyuncs.com/google_containers 
networking:
  podSubnet: 10.244.0.0/16 
  serviceSubnet: 10.96.0.0/12 



kubeadm init --config kubeadm.conf --ignore-preflight-errors=all  
```



#### kubeadmin init步骤：

1. [preflight]  环境检查
2. [kubelet-start]   准备kublet配置文件并启动/var/lib/kubelet/config.yaml
3. [certs]   证书目录 /etc/kubernetes/pki
4. [kubeconfig]   kubeconfig是用于连接k8s的认证文件 
5. [control-plane]  静态pod目录  /etc/kubernetes/manifests  启动组件用的
6. [etcd]   etcd的静态pod目录
7. [upload-config]  kubeadm-config存储到kube-system的命名空间中
8. [mark-control-plane]  给master节点打污点，不让pod分配
9. [bootstrp-token]   用于引导kubernetes的证书



​	使用`kubeadm init`初始化成功后，可以拷贝 kubectl 使用的连接k8s认证文件到默认路径：

```shell
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
kubectl get nodes

#  NAME         STATUS   ROLES    AGE   VERSION
#  k8s-master   Ready    master   2m    v1.20.0
```



###  worker节点加入Kubernetes Node

​	在Worker节点（192.168.159.144/145）执行。

​	向集群添加新节点，执行在node结点上：

```shell
kubeadm join 192.168.159.143:6443 --token esce21.q6hetwm8si29qxwn --discovery-token-ca-cert-hash sha256:00603a05805807501d7181c3d60b478788408cfe6cedefedb1f97569708be9c5
```

​	PS：**默认token有效期为24小时，当过期之后，上述token就不可用了**。这时就需要重新创建token，在master节点中操作：

```shell
#方法1
kubeadm token create
kubeadm token list    
openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //'

#方法2
kubeadm token create --print-join-command
```

​	在worker节点中使用最新的token：

```shell
kubeadm join 192.168.159.143:6443 --token nuja6n.o3jrhsffiqs9swnu --discovery-token-ca-cert-hash sha256:63bca849e0e01691ae14eab449570284f0c3ddeea590f8da988c07fe2729e924
```

**PS：`kubeadm token create --print-join-command`        可用于生成 `kubeadm join` 命令**





### 部署容器网络（CNI）

​	在master节点中部署。注意：只需要部署下面其中一个，推荐Calico。

#### Calico（推荐）

​	Calico是一个纯三层的数据中心网络方案，Calico支持广泛的平台，包括Kubernetes、OpenStack等。它在每一个计算节点利用 Linux Kernel 实现了一个高效的虚拟路由器（ vRouter） 来负责数据转发，而每个 vRouter 通过 BGP 协议负责把自己上运行的 workload 的路由信息向整个 Calico 网络内传播。此外，Calico  项目还实现了 Kubernetes 网络策略，提供ACL功能。

​	可以使用本仓库中提供的[calico.yaml](./calico.yaml)，这里使用的是3.18版本（官网：https://docs.tigera.io/archive/v3.18/getting-started/kubernetes/self-managed-onprem/onpremises）

​	下载完后还需要修改里面配置项：

- 定义Pod网络（CALICO_IPV4POOL_CIDR），与前面kubeadmin.conf文件中的podSubnet配置一样
- 选择工作模式（CALICO_IPV4POOL_IPIP），支持**BGP（Never）**、**IPIP（Always）**、**CrossSubnet**（开启BGP并支持跨子网）。选择默认的Always即可。

​	修改完后应用清单：

```shell
kubectl apply -f calico.yaml
kubectl get pods -n kube-system
```

**PS：可能会出现calico pod处于not ready状态**

​	解决方法：修改kube-proxy为IPVS模式：

```shell
yum install ipset -y
yum install ipvsadm -y
cat > /etc/sysconfig/modules/ipvs.modules <<EOF
#!/bin/bash
modprobe -- ip_vs
modprobe -- ip_vs_rr
modprobe -- ip_vs_wrr
modprobe -- ip_vs_sh
modprobe -- nf_conntrack_ipv4
EOF
# 使用lsmod | grep -e ip_vs -e nf_conntrack_ipv4命令查看是否已经正确加载所需的内核模块。
chmod 755 /etc/sysconfig/modules/ipvs.modules && bash /etc/sysconfig/modules/ipvs.modules && lsmod | grep -e ip_vs -e nf_conntrack_ipv4

kubectl edit cm kube-proxy -n kube-system   #将其中的mode=""  改为mode="ipvs"
kubectl rollout restart daemonset kube-proxy -n kube-system 
```

​	也可以参考这篇博客：

​	https://blog.csdn.net/qq_39698985/article/details/123960741

#### Flannel（备选）

​	Flannel是CoreOS维护的一个网络组件，Flannel为每个Pod提供全局唯一的IP，Flannel使用ETCD来存储Pod子网与Node IP之间的关系。flanneld守护进程在每台主机上运行，并负责维护ETCD信息和路由数据包。

​		可以使用本仓库中提供的[flannel.yaml](./flannel.yaml)

```shell
#  wget https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml

kubectl apply -f flannel.yml
```





## 测试kubernetes集群

1. 验证Pod工作

   在Kubernetes集群中创建一个pod，验证是否正常运行：

   ```shell
   kubectl create deployment nginx --image=nginx
   
   kubectl expose deployment nginx --port=80 --type=NodePort
   
   kubectl get pod
   ```

   访问地址：http://<NodeIP>:<NodePort>  

   

2. 验证DNS解析

   临时启动个pod，进去之后ping一下外网，看看能不能ping通，再nslookup kube-dns.kube-system：

```shell
kubectl run dns-test -it --rm --image=busybox:1.28.4  -- sh
```



​	到此，集群就搭建好了，可以开始在K8s的海洋中遨游了~





## 附录

### 一、解决k8s集群 Unable to connect to the server: x509: certificate is valid for xxx, not xxx问题

​	错误：Unable to connect to the server: x509: certificate is valid for xxx, not xxx 的解决方案

​	为了能使本地能连接k8s集群更好的测试client-go的功能，我在服务器上为本地签发了kubeconfig文件，放到本地之后出现如下的错误。

```css
➜  ~ kubectl get node
Unable to connect to the server: x509: certificate is valid for 10.96.0.1, 172.25.1.100, not 10.8.5.5
```

​	通过查阅资料发现了一个kubectl的参数`--insecure-skip-tls-verify`，加上这个参数之后确实好使了，但是，总是感觉治标不治本，所以经过一番查阅是apiserver的证书中没有添加`10.8.5.5`这个ip导致的，需要重新生成一下证书，具体操作如下：

#### 1.查看apiserver证书信息

```shell
cd /etc/kubernetes/pki
openssl x509 -noout -text -in apiserver.crt |grep IP
            DNS:k8s-master,DNS:kubernetes,DNS:kubernetes.default,DNS:kubernetes.default.svc,DNS:kubernetes.default.svc.cluster.local, IP Address:10.96.0.1, IP Address:172.25.1.100
```

​	从上面可以看出ip中并没有报错信息中的`10.8.5.5`这个IP地址，所以需要重新生成。

#### 2.删除旧证书

​	为了保险起见，这里选择将证书移动到其他位置。

```shell
mkdir -pv /opt/cert
mv apiserver.* /opt/cert
```

#### 3.生成新的apiserver证书

```shell
kubeadm init phase certs apiserver \
--apiserver-advertise-address 172.25.1.100 \
--apiserver-cert-extra-sans  10.96.0.1 \
--apiserver-cert-extra-sans 10.8.5.5 \
--apiserver-cert-extra-sans 10.8.5.6 \
--apiserver-cert-extra-sans 10.8.5.7 \
--apiserver-cert-extra-sans 10.8.5.8 \
--apiserver-cert-extra-sans 10.8.5.9 \
--apiserver-cert-extra-sans 10.8.5.10 \
--apiserver-cert-extra-sans 10.8.5.11 \
--apiserver-cert-extra-sans 10.8.5.12 \
--apiserver-cert-extra-sans 10.8.5.13 \
--apiserver-cert-extra-sans 10.8.5.14 \
--apiserver-cert-extra-sans 10.8.5.15 \
--apiserver-cert-extra-sans 10.8.5.16 \
--apiserver-cert-extra-sans 10.8.5.17
```

`--apiserver-cert-extra-sans`参数后可以加上需要添加的IP地址，这里为了省事儿一次性添加了多个，具体情况按需添加即可。

#### 4.检查证书

```shell
ls apiserver.*
apiserver.crt  apiserver.key
```

​	通过检查可以看到新的证书已经成了，现在只需要重启apiserver即可。如果出现问题，可以删除新的证书，将老的证书移回原位，重启apiserver即可。

#### 5.重启服务

```css
systemctl restart kubelet.service
```

#### 6.验证

```shell
➜  ~ kubectl get pod
NAME                                      READY   STATUS    RESTARTS   AGE
kucc4                                     3/3     Running   0          39d
nfs-client-provisioner-585486cc88-wqmzj   1/1     Running   6          40d
nginx-kusc00401                           1/1     Running   0          39d
web-server                                1/1     Running   0          39d
```

​	可以看到现在不加参数也不出现报错了，到这里就已经大功告成了。

### 二、Kubernetes一键部署利器：kubeadm

​	Kubernetes 有三个 Master 组件 kube-apiserver、kube-controller-manager、kube-scheduler，而它们都会被使用 Pod 的方式部署起来。

​	你可能会有些疑问：这时，Kubernetes 集群尚不存在，难道 kubeadm 会直接执行 docker run 来启动这些容器吗？

​	当然不是。

​	**在 Kubernetes 中，有一种特殊的容器启动方法叫做“Static Pod”。它允许你把要部署的 Pod 的 YAML 文件放在一个指定的目录里。这样，当这台机器上的 kubelet 启动时，它会自动检查这个目录，加载所有的 Pod YAML 文件，然后在这台机器上启动它们**。

​	从这一点也可以看出，kubelet 在 Kubernetes 项目中的地位非常高，在设计上它就是一个完全独立的组件，而其他 Master 组件，则更像是辅助性的系统容器。

​	在 kubeadm 中，Master 组件的 YAML 文件会被生成在 /etc/kubernetes/manifests 路径下。比如，kube-apiserver.yaml：

```yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    scheduler.alpha.kubernetes.io/critical-pod: ""
  creationTimestamp: null
  labels:
    component: kube-apiserver
    tier: control-plane
  name: kube-apiserver
  namespace: kube-system
spec:
  containers:

  - command:
    - kube-apiserver
    - --authorization-mode=Node,RBAC
    - --runtime-config=api/all=true
    - --advertise-address=10.168.0.2
      ...
    - --tls-cert-file=/etc/kubernetes/pki/apiserver.crt
    - --tls-private-key-file=/etc/kubernetes/pki/apiserver.key
      image: k8s.gcr.io/kube-apiserver-amd64:v1.11.1
      imagePullPolicy: IfNotPresent
      livenessProbe:
      ...
      name: kube-apiserver
      resources:
      requests:
        cpu: 250m
      volumeMounts:
    - mountPath: /usr/share/ca-certificates
      name: usr-share-ca-certificates
      readOnly: true
      ...
      hostNetwork: true
      priorityClassName: system-cluster-critical
      volumes:
  - hostPath:
    path: /etc/ca-certificates
    type: DirectoryOrCreate
    name: etc-ca-certificates
      ...
```

​	而一旦这些 YAML 文件出现在被 kubelet 监视的 /etc/kubernetes/manifests 目录下，kubelet 就会自动创建这些 YAML 文件中定义的 Pod，即 Master 组件的容器。

​	1. Linux 下生成证书，主流的选择应该是 OpenSSL，还可以使用 GnuGPG，或者 keybase。 2. Kubernetes 组件之间的交互方式：HTTP/HTTPS、gRPC、DNS、系统调用等。



## Reference

- https://kubernetes.io/zh/docs/reference/setup-tools/kubeadm/kubeadm-init/#config-file 
- https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/#initializing-your-control-plane-node
- https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-join 
-  https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/#pod-network 
-  https://docs.projectcalico.org/getting-started/kubernetes/quickstart 

