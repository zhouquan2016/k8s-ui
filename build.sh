docker build ./dashboard-server -t registry.cn-hangzhou.aliyuncs.com/zhqn/k8s-dashboard-server:1.0.0 && \
docker push registry.cn-hangzhou.aliyuncs.com/zhqn/k8s-dashboard-server:1.0.0 && \
docker build ./dashboard-web -t registry.cn-hangzhou.aliyuncs.com/zhqn/k8s-dashboard-web:1.0.0 && \
docker push registry.cn-hangzhou.aliyuncs.com/zhqn/k8s-dashboard-web:1.0.0