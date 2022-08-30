if [ "$2" == "rebuild" ]; then
  ./build.sh
  if [ $! -ne 0  ]; then
      exit
  fi
fi
kubectl apply -f ./k8s.yaml && nohup kubectl port-forward --address localhost,192.168.48.136 deploy/k8s-dashdoard-server 8080:8080 &