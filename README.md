# 项目介绍

做一个类似 Quora/ZhiHu 的在线问答/发帖社区项目。

### 接口文档

使用 postman 调试接口，可访问网页版在线文档。

[Link](https://zhuqimeng-4383979.postman.co/workspace/frog's-Workspace~89de5974-f978-4fdc-8d83-c1ae9bdfc668/collection/49204767-9aa78575-a604-4ffc-97b1-fd507c030477?action=share&creator=49204767)

### 基本功能

1. 基本的用户系统（已完成）
   1. 用户的注册与登录。
   2. 用户的基本信息如头像、个人简介等内容。
   3. 基本的用户鉴权。
   4. 密码的加盐加密。
   5. 用户个人信息的修改
2. 发布文章、问题；以及能够进行问题的回答、文章的评论。（已完成）
   1. 发布/获取文章、问题
   2. 文章格式规范为markdown，且文章可以插入markdown格式的图片链接
   3. 能够对文章、问题进行回复/评论
   4. 对问题进行关注、可以对文章、问题回答进行点赞。
   5. 按照问题关注度进行排序，做一个热度榜
3. 关注用户，并能够接收关注的人的动态，对用户下发通知（已完成）
   1. 关注用户和取消关注
   2. 关注列表和关注者列表。
   3. 可以接收已关注用户的所有动态，采用了异步写扩散方式的 feed 流设计。
4. 内容的搜索，可以使用 MySQL 的全文索引来做。（已完成）
   1. 使用 MySQL95 的全文索引实现搜索功能。
   2. 只对文章的标题和智能总结进行搜索，若搜索词不直接出现在文章中可能无法检测到。
5. 使用了滑动窗口算法实现限制点赞、评论等操作频率（已完成）

### 进阶功能

1. 加入缓存策略（已完成）
   1. 对基本的文章内容等数据进行缓存，并采取合理的缓存策略。
   2. 使用 go 内置的 singleflight、布隆过滤器等来应对缓存击穿、缓存雪崩这些问题。
2. Docker 打包、部署和安全
   1. 通过编写 dockerfile 文件来将项目打包成镜像。
   2. 推送到 dockerhub 或其他的镜像仓库上。
   3. 部署到自己的服务器上。
   4. 考虑防止 XSS、SQL 注入、CSRF 等常见安全风险。
3. 接口文档以及前端界面（已完成）
   1. 使用了 postman 保存接口文档。
4. 站内私信（使用了 WebSocket 实现）
   1. 互相关注的人可以实时发送消息私信。
   2. 接受者不在线时会把消息暂存在数据库中，当接收者上线时会自动收到未读消息。
5. 内容总结（已完成）
   1. 通过接入deepseek大模型实现。
   2. 每个文章的总结只需要生成一次即可，后续均访问该缓存。
6. 配置管理和日志管理（已完成）
   1. 通过 viper 加载如 mysql 相关信息的配置文件
   2. 通过 zap 日志库集成日志。

### 参考资料

朋友圈 feed 流设计：https://www.bilibili.com/video/BV1jGxLzDE23

Mysql 的全文索引：https://juejin.cn/post/7524910653062463528

怎么调大模型：https://www.cloudwego.io/zh/docs/eino/overview/

特别鸣谢： deepseek

### 项目架构

```
LanshanWinterProject
├─ .idea
│  ├─ dictionaries
│  │  └─ project.xml
│  ├─ LanshanWinterProject.iml
│  ├─ modules.xml
│  ├─ vcs.xml
│  └─ workspace.xml
├─ app
│  └─ api
│     ├─ configs
│     │  ├─ config.yaml
│     │  ├─ llm.go
│     │  ├─ logger.go
│     │  └─ sql.go
│     ├─ internal
│     │  ├─ middleware
│     │  │  └─ Auth
│     │  │     ├─ AuthStatu.go
│     │  │     ├─ AuthUser.go
│     │  │     └─ AuthVip.go
│     │  ├─ model
│     │  │  ├─ Document
│     │  │  │  ├─ answer.go
│     │  │  │  ├─ article.go
│     │  │  │  ├─ comment.go
│     │  │  │  └─ question.go
│     │  │  └─ User
│     │  │     ├─ feed.go
│     │  │     ├─ like.go
│     │  │     ├─ message.go
│     │  │     ├─ relation.go
│     │  │     └─ user.go
│     │  └─ service
│     │     ├─ Document
│     │     │  ├─ Answer
│     │     │  │  ├─ create.go
│     │     │  │  └─ read.go
│     │     │  ├─ Comment
│     │     │  │  ├─ create.go
│     │     │  │  └─ read.go
│     │     │  └─ dao
│     │     │     ├─ browse.go
│     │     │     ├─ create.go
│     │     │     ├─ like.go
│     │     │     └─ read.go
│     │     ├─ Message
│     │     │  ├─ client.go
│     │     │  ├─ controller.go
│     │     │  ├─ hub.go
│     │     │  └─ model.go
│     │     ├─ Serach
│     │     │  ├─ SerachArticle
│     │     │  ├─ SerachQuestion
│     │     │  └─ SerachUser
│     │     └─ User
│     │        ├─ dao
│     │        │  ├─ CreateUser.go
│     │        │  └─ ReadUser.go
│     │        ├─ Follow
│     │        │  ├─ feed.go
│     │        │  ├─ follow.go
│     │        │  └─ get.go
│     │        ├─ homepage.go
│     │        ├─ Sign
│     │        │  ├─ SignIn.go
│     │        │  └─ SignUp.go
│     │        └─ Upload
│     │           ├─ Homepage.go
│     │           └─ PasswordChange.go
│     └─ router
│        ├─ react.go
│        ├─ refresh.go
│        ├─ Router.go
│        └─ websocket.go
├─ go.mod
├─ go.sum
├─ main.go
├─ README.md
├─ Storage
│  ├─ Document
│  │  ├─ Answer
│  │  ├─ Article
│  │  │  ├─ 1770984344547407100-test1.md
│  │  │  ├─ 1771053263647365700-test1.md
│  │  │  └─ 1771053652519400600-test1.md
│  │  └─ Question
│  ├─ Log
│  │  └─ ZhiHu.log
│  └─ User
│     ├─ Avatar
│     └─ Profile
└─ utils
   ├─ files
   │  ├─ detectType.go
   │  └─ headerSet.go
   ├─ randoms
   │  └─ genSalt.go
   ├─ Strings
   │  ├─ Hash.go
   │  ├─ mdToPlain.go
   │  └─ split.go
   └─ tokens
      ├─ create.go
      ├─ model.go
      └─ read.go

```