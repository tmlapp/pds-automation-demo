## Deploy Details


# Deploy Mysql

``` bash
docker run --detach --name some-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=porxie -d mysql:8.0.31
```

# Deploy Zookeeper and Kafka

``` bash
docker compose up
```