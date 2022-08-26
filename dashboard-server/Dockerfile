#编译阶段
FROM  golang:alpine3.16 as build
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /build
ADD . .
RUN go build -o app main.go
#运行阶段
FROM alpine:3.16
WORKDIR /product
RUN mkdir -p /product/config
COPY --from=build /build/app /product/app
COPY --from=build /build/config/*          /product/config

CMD [ "/product/app" ]