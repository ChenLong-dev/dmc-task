# Docker

## Mysql
```
docker run --name mysql 
           -p 3306:3306 
           -e MYSQL_ROOT_PASSWORD=shanhai888888 
           -d mysql
```

```
docker run --name mysql -p 23306:3306 -e MYSQL_ROOT_PASSWORD=Shanhai*123 -d mysql
```

## Redis
```
docker run --name redis 
           -p 6379:6379 
           -v /docker/redis/data:/data 
           --restart unless-stopped  
           --appendonly yes 
           --requirepass 'shanhai888888'
           -d redis
```
- docker run：
    - 启动一个新的 Docker 容器。
- --name wh-redis：
    - 给容器指定一个名称，容器名为 wh-redis。指定名称便于后续操作，例如停止或启动容器时，可以通过名字直接操作容器。
- -p 6379:6379：
    - 将宿主机的端口 6379 映射到容器内的端口 6379。
    - 6379 是 Redis 默认的监听端口，因此这个映射会允许你通过宿主机的 6379 端口访问容器内的 Redis 服务。例如，你可以通过 localhost:6379 来连接 Redis。
- -v /root/RedisData:/data：
  - 使用 Docker 的 -v 参数来进行目录挂载，指定宿主机的目录 /root/RedisData 映射到容器中的 /data 目录。
  - Redis 会将它的持久化文件（如 dump.rdb 和 appendonly.aof）保存到 /data 目录中。通过挂载，Redis 的数据将保存到宿主机的 /root/RedisData 目录，从而实现数据持久化，即使容器删除或重启，数据也不会丢失。
- -d：
  - 让容器在后台运行（即“分离模式”），启动后不占用当前终端窗口。
- --restart unless-stopped：
    - 设置容器的重启策略为 unless-stopped，这意味着：
        - 如果容器意外停止（例如由于系统重启或 Docker 守护进程重启），容器将自动重启。
        - 但如果你手动停止容器（例如使用 docker stop 命令），容器将不会自动重启，除非你再次手动启动它。
- redis：
    - 使用官方的 Redis 镜像来创建和启动容器。Docker Hub 上有官方维护的 Redis 镜像，默认使用最新版本的 Redis。
- --appendonly yes：
    - 启用 Redis 的 AOF（Append Only File）持久化模式。AOF 记录每次写操作，确保数据实时保存到磁盘。即使 Redis 崩溃，AOF 也能恢复最近的操作记录。
    - 默认情况下，Redis 只使用 RDB 持久化（定期生成快照）。通过 --appendonly yes，Redis 将每次写入操作记录到 appendonly.aof 文件中，这比仅使用 RDB 持久化更可靠。
- --requirepass 'Your-password'：这个选项会告诉 Redis 容器启动时，设置密码为 Your-password。任何访问 Redis 的客户端都需要提供该密码。

```
# 启动 Redis 容器
docker run --name redis -p 6379:6379 -d --restart unless-stopped redis --requirepass 'shanhai888888'
```


## Mongo
```
docker run --name wh-mongo -p 27017:27017 -v /root/mongo-data:/data/db -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD='password' -d --restart unless-stopped mongo
```

- docker run：
  - 启动一个新的 Docker 容器。
- --name wh-mongo：
  - 为新容器指定一个名称 wh-mongo。这样可以方便地使用容器名称进行后续的管理，例如停止、启动或删除容器。
- -p 27017:27017：
  - 将宿主机的端口 27017 映射到容器内的端口 27017。
  - 27017 是 MongoDB 的默认端口，这样你可以通过访问 localhost:27017 来连接到 MongoDB 服务。
- -v /root/mongo-data:/data/db：
  - 使用 -v 参数将宿主机的 /root/mongo-data 目录挂载到容器的 /data/db 目录。
  - /data/db 是 MongoDB 默认的数据存储目录。通过挂载，MongoDB 的数据将保存在宿主机的 /root/mongo-data 目录中，从而实现数据的持久化，即使容器停止或删除，数据依然会保留在宿主机上。
- -e MONGO_INITDB_ROOT_USERNAME=admin：
  - 使用 -e 参数设置环境变量 MONGO_INITDB_ROOT_USERNAME 为 admin。这会创建一个名为 admin 的管理员用户。
- -e MONGO_INITDB_ROOT_PASSWORD='password'：
  - 使用 -e 参数设置环境变量 MONGO_INITDB_ROOT_PASSWORD 为 'password'。这会将管理员用户 admin 的密码设置为 password。
  - 注意：在实际生产环境中，请务必使用强密码以增强安全性。
- -d：
  - 让容器在后台运行（即“分离模式”），这使得命令行不会被容器的输出所占用。
- --restart unless-stopped：
  - 设置容器的重启策略为 unless-stopped。这意味着：
    - 如果容器因错误停止，Docker 将自动重启容器。
    - 但如果 Docker 守护进程重启，容器也会自动重启。 
    - 如果你手动停止容器（例如通过 docker stop），容器将不会自动重启，除非你手动再次启动它。
- mongo：
  - 使用官方的 MongoDB 镜像来创建和启动容器。默认情况下，它会拉取最新版本的 MongoDB


## Jaeger
```bash
docker run -d --name jaeger \
    -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
    -p 5775:5775/udp \
    -p 6831:6831/udp \
    -p 6832:6832/udp \
    -p 5778:5778 \
    -p 16686:16686 \
    -p 14268:14268 \
    -p 9411:9411 \
    --restart=always \
    jaegertracing/all-in-one:latest
```