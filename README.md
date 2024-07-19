# select-course

## 启动项目

### 1. 克隆项目到你的工作区

这里建议clone整个项目，而不是单独的某个tag

```bash
git clone 
```

### 2. 配置文件

你需要修改[.env.dev](.env.dev)
文件，其中你需要修改服务依赖地址。如果你修改了mysql配置信息那你就需要再修改[init.sql](init.sql)

### 3. 上传服务和部署依赖
> 这里推荐使用docker进行部署安装服务依赖

你需要将本项目上传到你的云服务器/虚拟机上，进入到项目执行如下命令：
```bash
docker-compose up
```
这里会进行在mysql创建用户和所需要的数据库，到这里服务依赖就部署完成了。你可以进入数据库进行查询是否创建成功等操作。

### 4. 如何预热
你需要执行`demo(N)/src/mock/get_course_test.go/TestPreheatMysql2Redis`进行预热操作。为了确保预热成功你需要到redis进行查询操作。

## 常见问题
### 运行测试脚本问题
其中执行测试文件依赖于当前运行环境下的配置文件，如果找不到会panic。如果修改了配置文件那么也需要修改测试脚本下所依赖的.env等
![img.png](assets/q1.png)