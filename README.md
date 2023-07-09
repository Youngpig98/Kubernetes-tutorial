





------



# 深入理解StatefulSet（三）：有状态应用实践

​	本节我们介绍一个非常典型的主从模式的 MySQL 集群

![4](.\4.webp)



## 思考题

如果我们现在的需求是：所有的读请求，只由 Slave 节点处理；所有的写请求，只由 Master 节点处理。那么，你需要在今天这篇文章的基础上再做哪些改动呢？















![image-20220903212805475](C:\Users\Young\AppData\Roaming\Typora\typora-user-images\image-20220903212805475.png)

![image-20220903212750165](C:\Users\Young\AppData\Roaming\Typora\typora-user-images\image-20220903212750165.png)

![image-20220903212846374](C:\Users\Young\AppData\Roaming\Typora\typora-user-images\image-20220903212846374.png)

![image-20220903233551848](C:\Users\Young\AppData\Roaming\Typora\typora-user-images\image-20220903233551848.png)

![image-20220903233852371](C:\Users\Young\AppData\Roaming\Typora\typora-user-images\image-20220903233852371.png)

![image-20220903235142410](C:\Users\Young\AppData\Roaming\Typora\typora-user-images\image-20220903235142410.png)

![image-20220903235324744](C:\Users\Young\AppData\Roaming\Typora\typora-user-images\image-20220903235324744.png)

![image-20220904000202139](C:\Users\Young\AppData\Roaming\Typora\typora-user-images\image-20220904000202139.png)

![image-20220904000402845](C:\Users\Young\AppData\Roaming\Typora\typora-user-images\image-20220904000402845.png)

​	resync到底有啥用？我看代码就是把index中的数据全拿出来再放入deltafifo中，然后事件是sync执行了一遍onupdate的自定义操作，是怕前面的自定义操作没执行成功吗？我看网上很多说是为了保证数据一致性，但是根本没保证呀

​	在处理 SharedInformer 事件回调时，可能存在处理失败的情况，定时的 Resync 能让处理失败的事件重新处理。







![image-20220904001252325](C:\Users\Young\AppData\Roaming\Typora\typora-user-images\image-20220904001252325.png)

![image-20220904001307127](C:\Users\Young\AppData\Roaming\Typora\typora-user-images\image-20220904001307127.png)
