FROM golang:latest
MAINTAINER "as"

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct"


# 移动到用于存放生成的二进制文件
WORKDIR $GOPATH/src03/summer_exam
COPY . $GOPATH/src03/summer_exam


RUN go build .

EXPOSE 8080 3306

ENTRYPOINT ["./exam"]