kubectl delete -f test/autoscale.yaml
kubectl create -f test/autoscale.yaml
make --directory anchor/scheduler run
