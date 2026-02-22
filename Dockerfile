# 构建阶段
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置Go环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 复制go mod文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download -x

# 复制源代码
COPY . .

# 编译应用（main.go在根目录）
RUN go build -o zhihu-app main.go

# 最终运行阶段
FROM alpine:latest

# 安装必要的工具
RUN apk --no-cache add ca-certificates tzdata

# 创建必要的目录
RUN mkdir -p /app/Storage/{Document/{Answer,Article,Question},Log,User/{Avatar,Profile}} \
    && mkdir -p /app/app/api/configs

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/zhihu-app /app/

# 复制配置文件
COPY --from=builder /app/app/api/configs/config.yaml /app/app/api/configs/

# 设置工作目录
WORKDIR /app

# 暴露应用端口（根据您的实际端口修改）
EXPOSE 8080

# 设置时区
ENV TZ=Asia/Shanghai

# 运行应用
CMD ["./zhihu-app"]