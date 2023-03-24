# 什么是RESTful风格？



## 一、概述

​	REST（英文：Representational State Transfer，简称REST）是一种软件架构风格、设计风格，而不是标准，只是提供了一组设计原则和约束条件。它主要用于客户端和服务器交互类的软件。基于这个风格设计的软件可以更简洁，更有层次，更易于实现缓存等机制。



## 二、什么是RESTful

​	REST 指的是一组架构（约束条件）和原则。满足这些（约束条件）和（原则）的应用程序或设计就是 RESTful。它就是一个**资源定位、资源操作**的风格。不是标准也不是协议，只是一种风格。基于这个风格设计的软件可以更简洁，更有层次，更易于实现缓存等机制。



## 三、restful有什么特点

1. 每一个URI代表一种资源，独一无二
2. 客户端和服务器之间，传递这种资源的某种表现层
3. 客户端通过四个HTTP动词，对服务器端资源进行操作，实现"表现层状态转化"。

## 四、具体用例

​	RESTful架构风格规定，数据的元操作，即CRUD(create, read, update和delete,即数据的增删查改)操作，分别对应于HTTP方法：`GET`用来获取资源，`POST`用来新建资源（也可以用于更新资源），`PUT`用来更新资源，`DELETE`用来删除资源，这样就统一了数据操作的接口，仅通过HTTP方法，就可以完成对数据的所有增删查改工作。

即：

- GET（SELECT）：从服务器取出资源（一项或多项）。
- POST（CREATE）：在服务器新建一个资源。
- PUT（UPDATE）：在服务器更新资源（客户端提供完整资源数据）。
- PATCH（UPDATE）：在服务器更新资源（客户端提供需要修改的资源数据）。
- DELETE（DELETE）：从服务器删除资源。

 

使用RESTful操作资源 ：

- 【GET】 /users # 查询用户信息列表
- 【GET】 /users/1001 # 查看某个用户信息
- 【POST】 /users # 新建用户信息
- 【PUT】 /users/1001 # 更新用户信息(全部字段)
- 【PATCH】 /users/1001 # 更新用户信息(部分字段)
- 【DELETE】 /users/1001 # 删除用户信息



**RESTful的CRUD**

@RequestMapping：通过设置method属性的CRUD，可以将同一个URL映射到不同的HandlerMethod方法上。

@GetMapping、@PostMapping、@PutMapping、@DeleteMapping注解同@RequestMapping注解的method属性设置。

 

我们都知道在没有出现RESTful风格之前，我们的代码是这样的：

<img src="https://img2020.cnblogs.com/blog/1827620/202007/1827620-20200721212516730-1457422532.png" alt="img" style="zoom:50%;" />

那么RESTful风格要求的是这样的：

<img src="https://img2020.cnblogs.com/blog/1827620/202007/1827620-20200721212559468-1673462679.png" alt="img" style="zoom:50%;" />

 

## 五、传统风格与RestFul风格对比

1. 传统方式操作资源

   ​	通过不同的参数来实现不同的效果！方法单一！

- http://127.0.0.1/item/queryItem.action?id=1 （查询,GET）
- http://127.0.0.1/item/saveItem.action （新增,POST）
- http://127.0.0.1/item/updateItem.action （更新,POST）
- http://127.0.0.1/item/deleteItem.action?id=1 （删除,GET或POST）

2. RestFul方式操作资源
   可以通过不同的请求方式来实现不同的效果！
   如下：请求地址一样，但是功能可以不同！

- http://127.0.0.1/item/1 （查询,GET）
- http://127.0.0.1/item （新增,POST）
- http://127.0.0.1/item （更新,PUT）
- http://127.0.0.1/item/1 （删除,DELETE）



## 六、k8s restful API 结构分析

​	k8s的api-server组件负责提供restful api访问端点, 并且将数据持久化到etcd server中. 那么k8s是如何组织它的restful api的?

### namespaced resources

​	所谓的namespaced resources,就是这个resource是从属于某个namespace的, 也就是说它不是cluster-scoped的资源. 比如pod, deployment, service都属于namespaced resource. 那么我们看一下如何请求一个namespaced resources：

```
http://localhost:8080/api/v1/namespaces/default/pods/test-pod
```

​	可以看出, 该[restful](https://so.csdn.net/so/search?q=restful&spm=1001.2101.3001.7020) api的组织形式是:

| api  | api版本 | namespaces | 所属的namespace | 资源种类 | 所请求的资源名称 |
| ---- | ------- | ---------- | --------------- | -------- | ---------------- |
| api  | v1      | namespaces | default         | pods     | test-pod         |

​	这里api version如果是v1的话,表示这是一个很稳定的版本了, 以后不会有大的修改,并且当前版本所支持的所有特性以后都会兼容. 而如果版本号是v1alpha1, v1beta1之类的,则不保证其稳定性。

### non-namespaced resources

```
http://localhost:8080/apis/rbac.authorization.k8s.io/v1/clusterroles/test-clusterrole
```

​	这里可以观察到它clusterrole与pod不同, apis表示这是一个非核心api。 rbac.authorization.k8s.io指代的是api-group, 另外它没有namespaces字段, 其他与namespaced resources类似.不再赘述。

### non-resource url

​	这类资源和pod, clusterrole都不同. 例如：

```
http://localhost:8080/healthz/etcd
```

​	这就是用来确认etcd服务是不是健康的。它不属于任何namespace,也不属于任何api版本。

​	总结, k8s的REST API的设计结构为:

```
[api/apis] / api-group              / api-version / namespaces / namespace-name / resource-kind / resource-name

apis      /  rbac.authorization.k8s.io / v1       / namespaces / default         / roles       / test-role
```





## 七、结论

​	RESTful风格要求每个资源都使用 URI (Universal Resource Identifier) 得到一个唯一的地址。所有资源都共享统一的接口，以便在客户端和服务器之间传输状态。使用的是标准的 HTTP 方法，比如 GET、PUT、POST 和 DELETE。
​	总之就是REST是一种写法上规范，获取数据或者资源就用GET，更新数据就用PUT，删除数据就用DELETE，然后规定方法必须要传入哪些参数，每个资源都有一个地址。
