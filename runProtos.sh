docker stop $(docker ps -aq)
docker-compose up --build -d
export PATH=$PATH:$HOME/go/bin
export PATH=$PATH:/usr/local/go/bin
protoc -I proto --go_out=./proto-output --go-grpc_out=./proto-output ./proto/*.proto
