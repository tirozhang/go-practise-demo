docker ps // 查看所有正在运行容器
docker stop containerId // containerId 是容器的ID

docker ps -a // 查看所有容器
docker ps -a -q // 查看所有容器ID

docker start $(docker ps -a -q) // start启动所有停止的容器
docker stop $(docker ps -a -q) // stop停止所有容器
docker rm $(docker ps -a -q) // remove删除所有容器

docker rmi $(docker images -q) //删除所有镜像
docker rm $(docker ps -qf status=exited) //删除所有已经停止的容器
