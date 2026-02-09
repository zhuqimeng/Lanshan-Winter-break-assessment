# 项目介绍

做一个类似 Quora/zhihu 的在线问答/发帖社区项目。

### 接口文档

使用 postman 调试接口，可访问网页版在线文档。

[Link](https://zhuqimeng-4383979.postman.co/workspace/frog's-Workspace~89de5974-f978-4fdc-8d83-c1ae9bdfc668/collection/49204767-9aa78575-a604-4ffc-97b1-fd507c030477?action=share&creator=49204767)

### 基本要求

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
3. 关注用户，并能够接收关注的人的动态，(**非必要**)对用户下发通知
   1. 关注用户和取消关注
   2. 关注列表和关注者列表。
   3. （**非必要**）可以接收已关注用户的所有动态，如果要做的话建议看看b站上朋友圈/微博 feed 流的设计。
4. 内容的搜索，可以使用 MySQL 的全文索引来做。（**进阶，推荐做**）
   1. 使用 MySQL 的全文索引实现搜索功能。
5. 文章状态(草稿，已删除，未删除)
6. 管理员禁言
7. 使用了滑动窗口算法实现防刷机制(限制点赞频率，评论频率)

### 进阶要求

1. 加入缓存策略
   1. 对基本的文章内容等数据进行缓存，并采取合理的缓存策略。
   2. 使用 go 内置的 singleflight、布隆过滤器等来应对缓存击穿、缓存雪崩这些问题。
2. Docker 打包、部署和安全
   1. 通过编写 dockerfile 文件来将项目打包成镜像。
   2. 推送到 dockerhub 或其他的镜像仓库上。
   3. 部署到自己的服务器上。
   4. 考虑防止 XSS、SQL 注入、CSRF 等常见安全风险。
3. 接口文档以及前端界面
   1. 使用了 postman 保存接口文档。
   2. 借助 agent 工具如 copilot、codex 等扫描你的代码并生成前端界面，便于调试。
4. 站内通知和私信
   1. 互相关注的人可以互相发送消息私信。
   2. 当关注的问题或者有人回复你的评论会发送站内通知。
5. 盐选会员
   1. 加入会员的角色管理。
   2. 非会员只能访问某些文章的部分内容。
6. 文章推荐和内容总结
   1. 可以通过接入大模型实现。
   2. 每个文章的总结只需要生成一次即可，后续均访问该缓存，当更新的时候可以重新生成内容。
   3. 文章推荐可以选择调大模型接口做智能推荐，这里不做重点要求，因为一般是算法干的事。
7. 配置管理和日志管理（已完成）
   1. 通过 viper 或者其他方式加载如 mysql 相关信息的配置文件
   2. 通过 zap、zerolog 等专门的日志库集成日志。
8. 内容分发与访问加速方案等。

### 参考资料

朋友圈 feed 流设计：https://www.bilibili.com/video/BV1jGxLzDE23

Mysql 的全文索引：https://juejin.cn/post/7524910653062463528

怎么调大模型：https://www.cloudwego.io/zh/docs/eino/overview/

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
│     │     │     └─ read.go
│     │     ├─ Serach
│     │     │  ├─ SerachArticle
│     │     │  ├─ SerachQuestion
│     │     │  └─ SerachUser
│     │     └─ User
│     │        ├─ Admin
│     │        ├─ dao
│     │        │  ├─ CreateUser.go
│     │        │  └─ ReadUser.go
│     │        ├─ Follow
│     │        │  ├─ followers
│     │        │  └─ following
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
│        └─ Router.go
├─ go.mod
├─ go.sum
├─ main.go
├─ README.md
├─ Storage
│  ├─ Document
│  │  ├─ Answer
│  │  │  └─ 1770377843667167300-test2.md
│  │  ├─ Article
│  │  │  └─ 1770370796889802200-test1.md
│  │  └─ Question
│  │     └─ 1770372183202854000-test1.md
│  ├─ Log
│  │  └─ ZhiHu.log
│  └─ User
│     ├─ Avatar
│     │  └─ test1_1770370243591934400.png
│     └─ Profile
│        └─ test1_1770370513029485900.md
└─ utils
   ├─ files
   │  ├─ detectType.go
   │  └─ headerSet.go
   ├─ randoms
   │  └─ genSalt.go
   ├─ strings
   │  └─ Hash.go
   └─ tokens
      ├─ create.go
      ├─ model.go
      └─ read.go

```