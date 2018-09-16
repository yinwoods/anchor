# Liman
[![Build Status](https://travis-ci.org/salihciftci/liman.svg?branch=master)](https://travis-ci.org/salihciftci/liman) [![Go Report Card](https://goreportcard.com/badge/github.com/salihciftci/liman)](https://goreportcard.com/report/github.com/salihciftci/liman)

![alt text](https://img.salih.co/liman/v0.6/logo.png "Liman")

Web application for monitoring docker. Monitor docker inside the docker. Written in Go.

----
![screenshot](https://img.salih.co/liman/v0.6/dashboard.png)

## Features

* Monitoring docker
    * Containers
    * Logs
    * Images
    * Stats
    * Volumes
    * Networks
* Notifications
* Restful API

## Installation

[Download](https://github.com/salihciftci/liman/releases) and run the latest binary or ship with Docker.


```
docker run -it -v /var/run/docker.sock:/var/run/docker.sock salihciftci/liman
```

Note: the `-v /var/run/docker.sock:/var/run/docker.sock` option can be used in Linux environments only. 

## API Usage

Basic usage:
```
curl -i http://localhost:8080/api/status?key=XXX
```

More examples and all end points can be found in [wiki](https://github.com/salihciftci/liman/wiki/API-Usage).


| Screenshots | | |
|:-------------:|:-------:|:-------:|
|![Dashboard](https://img.salih.co/liman/v0.6/dashboard.png)|![Containers](https://img.salih.co/liman/v0.6/containers.png)|![Images](https://img.salih.co/liman/v0.6/images.png)|
|![Stats](https://img.salih.co/liman/v0.6/stats.png)|![Volumes](https://img.salih.co/liman/v0.6/volumes.png)|![Networks](https://img.salih.co/liman/v0.6/networks.png)|
|![Logs](https://img.salih.co/liman/v0.6/logs.png)|![Notifications](https://img.salih.co/liman/v0.6/notifications.png)|![Settings](https://img.salih.co/liman/v0.6/settings.png)|

## License
MIT
# scheduler

Toy scheduler for use in Kubernetes demos.

## Usage

Annotate each node using the annotator command:

```
kubectl proxy
```
```
Starting to serve on 127.0.0.1:8080
```

### Create a deployment

```
kubectl create -f deployments/nginx.yaml
```
```
deployment "nginx" created
```

The nginx pod should be in a pending state:

```
kubectl get pods
```
```
NAME                     READY     STATUS    RESTARTS   AGE
nginx-1431970305-mwghf   0/1       Pending   0          27s
```

### Run the Scheduler

List the nodes and note the price of each node.

```
annotator -l
```
```
gke-k0-default-pool-728d327f-00lq 0.80
gke-k0-default-pool-728d327f-3vzg 0.40
gke-k0-default-pool-728d327f-nmz7 0.40
gke-k0-default-pool-728d327f-pxee 0.05
gke-k0-default-pool-728d327f-xm4i 1.60
gke-k0-default-pool-728d327f-zynj 0.40
```

Run the best price scheduler:

```
scheduler
```
```
2016/08/19 11:16:25 Starting custom scheduler...
2016/08/19 11:16:28 Successfully assigned nginx-1431970305-mwghf to gke-k0-default-pool-728d327f-pxee
2016/08/19 11:16:35 Shutdown signal received, exiting...
2016/08/19 11:16:35 Stopped reconciliation loop.
2016/08/19 11:16:35 Stopped scheduler.
```

> Notice the pending nginx pod is deployed to the node with the lowest cost.

## Run the Scheduler on Kubernetes

```
kubectl create -f deployments/scheduler.yaml
```
``` 
deployment "scheduler" created
```

---

```
          _                               __    
   __  __(_)___ _      ______  ____  ____/ /____
  / / / / / __ \ | /| / / __ \/ __ \/ __  / ___/
 / /_/ / / / / / |/ |/ / /_/ / /_/ / /_/ (__  ) 
 \__, /_/_/ /_/|__/|__/\____/\____/\__,_/____/  
/____/        yinwoods.github.io           
```
