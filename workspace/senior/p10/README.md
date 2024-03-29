# Kubernetes“弃用Docker”是怎么回事？

## 什么是 CRI

​	要了解 Kubernetes 为什么要“弃用 Docker”，还得追根溯源，回头去看 Kubernetes 的发展历史。

​	2014 年，Docker 正如日中天，在容器领域没有任何对手，而这时 Kubernetes 才刚刚诞生，虽然背后有 Google 和 Borg 的支持，但还是比较弱小的。所以，Kubernetes 很自然就选择了在 Docker 上运行，毕竟“背靠大树好乘凉”，同时也能趁机“养精蓄锐”逐步发展壮大自己。

​	时间一转眼到了 2016 年，CNCF 已经成立一年了，而 Kubernetes 也已经发布了 1.0 版，可以正式用于生产环境，这些都标志着 Kubernetes 已经成长起来了，不再需要“看脸色吃饭”。于是它就宣布加入了 CNCF，成为了第一个 CNCF 托管项目，想要借助基金会的力量联合其他厂商，一起来“扳倒”Docker。

​	那它是怎么做的呢？在 2016 年底的 1.5 版里，Kubernetes 引入了一个新的接口标准：CRI ，Container Runtime Interface。

​	CRI 采用了 ProtoBuffer 和 gPRC，规定 kubelet 该如何调用容器运行时去管理容器和镜像，但这是一套全新的接口，和之前的 Docker 调用完全不兼容。

​	Kubernetes 意思很明显，就是不想再绑定在 Docker 上了，允许在底层接入其他容器技术（比如 rkt、kata 等），随时可以把 Docker“踢开”。

​	但是这个时候 Docker 已经非常成熟，而且市场的惯性也非常强大，各大云厂商不可能一下子就把 Docker 全部替换掉。所以 Kubernetes 也只能同时提供一个“折中”方案，在 kubelet 和 Docker 中间加入一个“适配器”，把 Docker 的接口转换成符合 CRI 标准的接口:



![](../../img/cri.jpg)

​	因为这个“适配器”夹在 kubelet 和 Docker 之间，所以就被形象地称为是“shim”，也就是“垫片”的意思。

​	有了 CRI 和 shim，虽然 Kubernetes 还使用 Docker 作为底层运行时，但也具备了和 Docker 解耦的条件，从此就拉开了“弃用 Docker”这场大戏的帷幕。

## 什么是 containerd

​	面对 Kubernetes“咄咄逼人”的架势，Docker 是看在眼里痛在心里，虽然有苦心经营了多年的社区和用户群，但公司的体量太小，实在是没有足够的实力与大公司相抗衡。

​	不过 Docker 也没有“坐以待毙”，而是采取了“断臂求生”的策略，推动自身的重构，**把原本单体架构的 Docker Engine 拆分成了多个模块，其中的 Docker daemon 部分就捐献给了 CNCF，形成了 containerd。**

​	containerd 作为 CNCF 的托管项目，自然是要符合 CRI 标准的。但 Docker 出于自己诸多原因的考虑，它只是在 Docker Engine 里调用了 containerd，外部的接口仍然保持不变，也就是说还不与 CRI 兼容。

​	由于 Docker 的“固执己见”，这时 Kubernetes 里就出现了两种调用链：

- 第一种是用 CRI 接口调用 dockershim，然后 dockershim 调用 Docker，Docker 再走 containerd 去操作容器。
- 第二种是用 CRI 接口直接调用 containerd 去操作容器。



![](../../img/cei2.jpg)

​	显然，由于都是用 containerd 来管理容器，所以这两种调用链的最终效果是完全一样的，但是第二种方式省去了 dockershim 和 Docker Engine 两个环节，更加简洁明了，损耗更少，性能也会提升一些。

​	在 2018 年 Kubernetes 1.10 发布的时候，containerd 也更新到了 1.1 版，正式与 Kubernetes 集成，同时还发表了一篇[博客文章](https://kubernetes.io/blog/2018/05/24/kubernetes-containerd-integration-goes-ga/)，展示了一些性能测试数据:



![](../../img/cri3.jpg)

​	从这些数据可以看到，containerd1.1 相比当时的 Docker 18.03，Pod 的启动延迟降低了大约 20%，CPU 使用率降低了 68%，内存使用率降低了 12%，这是一个相当大的性能改善，对于云厂商非常有诱惑力。

## 正式“弃用 Docker”

​	有了 CRI 和 containerd 这两件强大的武器，胜利的天平已经明显向 Kubernetes 倾斜了。

​	又是两年之后，到了 2020 年，Kubernetes 1.20 终于正式向 Docker“宣战”：kubelet 将弃用 Docker 支持，并会在未来的版本中彻底删除。

​	但由于 Docker 几乎成为了容器技术的代名词，而且 Kubernetes 也已经使用 Docker 很多年，这个声明在不断传播的过程中很快就“变味”了，“kubelet 将弃用 Docker 支持”被简化成了更吸引眼球的“Kubernetes 将弃用 Docker”。

​	这自然就在 IT 界引起了恐慌，“不明真相的广大群众”纷纷表示震惊：用了这么久的 Docker 突然就不能用了，Kubernetes 为什么要如此对待 Docker？之前在 Docker 上的投入会不会就全归零了？现有的大量镜像该怎么办？

​	其实，如果你理解了前面讲的 CRI 和 containerd 这两个项目，就会知道 Kubernetes 的这个举动也没有什么值得大惊小怪的，一切都是“水到渠成”的：**它实际上只是“弃用了 dockershim”这个小组件，也就是说把 dockershim 移出了 kubelet，并不是“弃用了 Docker”这个软件产品。**

​	所以，“弃用 Docker”对 Kubernetes 和 Docker 来说都不会有什么太大的影响，因为他们两个都早已经把下层都改成了开源的 containerd，原来的 Docker 镜像和容器仍然会正常运行，唯一的变化就是 Kubernetes 绕过了 Docker，直接调用 Docker 内部的 containerd 而已。这个关系你可以参考下面的这张图来理解:



<img src="../../img/cri4.jpg" style="zoom:67%;" />

​	当然，影响也不是完全没有。如果 Kubernetes 直接使用 containerd 来操纵容器，那么它就是一个与 Docker 独立的工作环境，彼此都不能访问对方管理的容器和镜像。换句话说，使用命令 docker ps 就看不到在 Kubernetes 里运行的容器了。

​	这对有的人来说可能需要稍微习惯一下，改用新的工具 crictl，不过用来查看容器、镜像的子命令还是一样的，比如 ps、images 等等，适应起来难度不大（但如果我们一直用 kubectl 来管理 Kubernetes 的话，这就是没有任何影响了）。

​	“宣战”之后，Kubernetes 原本打算用一年的时间完成“弃用 Docker”的工作，但它也确实低估了 Docker 的根基，到了 1.23 版还是没能移除 dockershim，不得已又往后推迟了半年，终于在2022年 5 月份发布的 1.24 版把 dockershim 的代码从 kubelet 里删掉了。自此，Kubernetes 彻底和 Docker“分道扬镳”，今后就是“大路朝天，各走一边”。

## Docker 的未来

​	那么，Docker 的未来会是怎么样的呢？难道云原生时代就没有它的立足之地了吗？

​	这个问题的答案很显然是否定的。

​	作为容器技术的初创者，Docker 的历史地位无人能够质疑，虽然现在 Kubernetes 不再默认绑定 Docker，但 Docker 还是能够以其他的形式与 Kubernetes 共存的。

​	首先，因为**容器镜像格式已经被标准化**了（OCI 规范，Open Container Initiative），Docker 镜像仍然可以在 Kubernetes 里正常使用，原来的开发测试、CI/CD 流程都不需要改动，我们仍然可以拉取 Docker Hub 上的镜像，或者编写 Dockerfile 来打包应用。

​	其次，**Docker 是一个完整的软件产品线**，不止是 containerd，它还包括了镜像构建、分发、测试等许多服务，甚至在 Docker Desktop 里还内置了 Kubernetes。

​	单就容器开发的便利性来讲，Docker 还是暂时难以被替代的，广大云原生开发者可以在这个熟悉的环境里继续工作，利用 Docker 来开发运行在 Kubernetes 里的应用。

​	再次，虽然 Kubernetes 已经不再包含 dockershim，但 Docker 公司却把这部分代码接管了过来，另建了一个叫 [cri-dockerd](https://github.com/mirantis/cri-dockerd)的项目，作用也是一样的，把 Docker Engine 适配成 CRI 接口，这样 kubelet 就又可以通过它来操作 Docker 了，就仿佛是一切从未发生过。

​	综合来看，Docker 虽然在容器编排战争里落败，被 Kubernetes 排挤到了角落，但它仍然具有强韧的生命力，多年来积累的众多忠实用户和数量庞大的应用镜像是它的最大资本和后盾，足以支持它在另一条不与 Kubernetes 正面交锋的道路上走下去。

​	而对于我们这些初学者来说，Docker 方便易用，具有完善的工具链和友好的交互界面，市面上很难找到能够与它媲美的软件了，应该说是入门学习容器技术和云原生的“不二之选”。至于 Kubernetes 底层用的什么，我们又何必太过于执着和关心呢？
