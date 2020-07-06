### Netbox testing environment

#### Requirement
1. docker
2. docker-compose

#### Build netbox docker containers 
```shell script
make build
```

### Setup a hostname rewrite in hosts file(Optional)
```shell script
127.0.0.1 netbox.k8s.me netbox netbox.you.like
```

#### Start 
```shell script
make compose-up
```

#### Stop 
```shell script
make compose-down
```