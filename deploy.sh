if [ "$2" == "rebuild" ]; then
  ./build.sh
  if [ $! -ne 0  ]; then
      exit
  fi
fi
kubectl apply -f ./k8s.yaml && nohup kubectl port-forward --address localhost,192.168.48.136 deploy/k8s-dashboard-server 8080:8080 &
while [ true ]; do
    state=$(kubectl get pod | grep k8s-dashboard-web | awk '{print $3}')
    if [ "$state" == "Running" ]; then
      nohup kubectl port-forward --address localhost,192.168.48.136 deploy/k8s-dashboard-server 8080:8080 &
      break
    elif [ "$state" == "Pending"  ]; then
      echo "wait pod running!"
      sleep 1
    else
      echo "ERROR:k8s-dashboard-web pod state ${state}"
      break
    fi
done