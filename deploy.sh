if [ "$1" == "rebuild" ]; then
  ./build.sh
  if [ $! -ne 0  ]; then
      exit
  fi
fi
kubectl apply -f ./k8s.yaml
for podName in "k8s-dashboard-server", "k8s-dashboard-web" ; do
  while [ true ]; do
    state=$(kubectl get pod | grep k8s-dashboard-web | awk '{print $3}')
    if [ "$state" == "Running" ]; then
      nohup kubectl port-forward --address localhost,192.168.48.136 deploy/k8s-dashboard-server 8080:8080 &
      wait $!
      break
    elif [ "$state" == "Pending"  ]; then
      echo "wait pod running!"
      sleep 1
    else
      break
    fi
  done
  if [ "$state" != "Running" ]; then
      echo "ERROR:k8s-dashboard-web pod state ${state}"
      exit
  fi
done
