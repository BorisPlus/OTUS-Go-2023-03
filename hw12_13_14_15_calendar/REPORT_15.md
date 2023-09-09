"user=hw12user password=hw12user host='127.0.0.1' database=hw12calendar search_path=hw12calendar

PG_DSN="host=172.18.0.1 port=5437 dbname=hw15 user=hw15user password=hw15user " make migrate-sql-up

make migrate-rmq-up: 

make compose-up



sudo iptables -P INPUT ACCEPT
sudo iptables -P FORWARD ACCEPT
sudo iptables -P OUTPUT ACCEPT
sudo iptables -t nat -F
sudo iptables -t mangle -F
sudo iptables -F
sudo iptables -X

service networking restart
sudo service networking restart
systemctl restart networking
docker rm -f $(docker ps -a -q)
docker rmi -f $(docker images -q)


# docker build --no-cache --build-arg MICROSERVICE=calendar --progress plain -f ./docker/Dockerfile --tag=hw15/calendar:v1 .
# docker build --no-cache --build-arg MICROSERVICE=sheduler --progress plain -f ./docker/Dockerfile --tag=hw15/sheduler:v1 .
# docker build --no-cache --build-arg MICROSERVICE=archiver --progress plain -f ./docker/Dockerfile --tag=hw15/archiver:v1 .
# docker build --no-cache --build-arg MICROSERVICE=sender   --progress plain -f ./docker/Dockerfile --tag=hw15/sender:v1 .
# docker builder prune
# docker build --build-arg MICROSERVICE=calendar --progress plain -f ./docker/Dockerfile --tag=hw15/calendar:v2 .
# docker build --build-arg MICROSERVICE=sheduler --progress plain -f ./docker/Dockerfile --tag=hw15/sheduler:v2 .
# docker build --build-arg MICROSERVICE=archiver --progress plain -f ./docker/Dockerfile --tag=hw15/archiver:v2 .
# docker build --build-arg MICROSERVICE=sender   --progress plain -f ./docker/Dockerfile --tag=hw15/sender:v2 .
