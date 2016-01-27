#Chaos

##QuickStart
+ 如需要可执行程序，请执行go install opensource/chaos/server (相对于$GOPATH/src下的项目放置位置。目前开发版本放置在示例的目录下)，之后到bin目录下找到server的可执行程序双击之后。打开浏览器输入`localhost:8080`即可登录私有云平台界面。同时rest api也全部生效。
+ 如希望直接执行，则请执行`go run $YOUR_PATH_PREFIX/server/*.go`之后即全部生效。

**目前界面提供如下功能**
+ 私有云首页入口
+ 罗列服务集群现状概况，并提供搜索排序功能

**目前提供如下rest api**
+ /deploy/apps/rollback 快速回滚到某一个版本或者回滚到最近的版本
+ /deploy/apps 对单个服务进行一键部署
+ /deploy/apps/updater 对单个或者多个服务进行新增或者更新
+ /deploy/groups 对单个或者多个组进行一键部署
+ /info 获取所有服务信息

##Introdution
chaos，卡俄斯，是希腊神话最初始的神——混沌。象征天地伊始，开辟鸿蒙。这也是chaos私有云平台所期望能够带给目前的服务现状的剧烈改变。

## Todo Recently
+ 一键部署相关的rest接口完善，并提供异步化
+ 结合marathon和mesos以及docker，提取服务的真实ip等关键信息，并进行视图呈现
+ 搭建事件动态刷新前后端server架构。对于服务变更进行第一时间监控和实时刷新。
+ 对centos6.5镜像进行完善。对网络问题进行优化。

## Contact Us
inf@zufangit.cn

## Changelog

**v0.3** —— **2016-01-25**
+ 提供一键部署的四个基本rest api
+ 模块抽象和拆解。代码优化

**v0.2** —— **2016-01-25**
+ 前后端架构重构。改为rest + angularjs-ajax的方式
+ 前端代码抽象和拆解。
+ 提供获取所有服务信息的接口

**v0.1** —— **2016-01-13**
+ chaos初始化