echo "Building monkey-ops docker image"

REGISTRY=$1
REPOSITORY=$2
IMAGE=$3
TAG=$4
PROXY=""
if [ "$#" -gt 4 ]; then
    PROXY="--build-arg https_proxy=$5"
fi

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./image/monkey-ops ./go/*.go 

if [ $? = 0 ]
then
	docker build --no-cache=true ${PROXY} -t ${REGISTRY}/${REPOSITORY}/${IMAGE}:${TAG} -f ./image/Dockerfile ./image
fi