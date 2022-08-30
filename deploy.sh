function portForward {
  podName=$1
  port=$2
  while [ true ]; do
    state=$(kubectl get pod | grep $podName | awk '{print $3}')
    if [ "$state" == "Running" ]; then
      if [ -f "$podName.pid" ]; then
          pid=${cat $podName.pid}
          kill -9 $pid
          rm -f $podName.pid
      fi
      nohup kubectl port-forward --address localhost,192.168.48.136 deploy/$podName port &  echo $! > $podName.pid
      break
    elif [ "$state" == "Pending"  ]; then
      echo "wait pod running!"
      sleep 1
    else
      break
    fi
  done
  if [ "$state" != "Running" ]; then
      echo "ERROR:$podName pod state ${state}"
      exit
  fi
}

if [ "$1" == "rebuild" ]; then
  ./build.sh
  if [ $! -ne 0  ]; then
      exit
  fi
fi
kubectl apply -f ./k8s.yaml && \
portForward k8s-dashboard-server 8080:8080 && \
portForward k8s-dashboard-web 8081:80