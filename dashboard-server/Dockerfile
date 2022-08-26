#编译阶段
FROM  golang:alpine3.16 as build
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /build
ADD . .
RUN go build -o app main.go
#运行阶段
FROM alpine3.16
WORKDIR /product
RUN mkdir -p /product/config
COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=build /build/app /product/app
COPY --from=build /build/config/*.yaml          /product/config
COPY --from=build /build/config/*.xml           /product/config
COPY --from=build /build/config/*.properties    /product/config
COPY --from=build /build/config/*.sql           /product/config

CMD [ "/product/app" ]