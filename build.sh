echo "Building monkey-ops docker image"

TAG="0.1.1"
PROXY=""

if [ "$#" -gt 0 ]; then
    PROXY="--build-arg https_proxy=$1"
fi

echo "+-- Building binary"
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./image/monkey-ops ./go/*.go 

if [ $? = 0 ]
then
    cd ./image
    echo "+-- Building docker image"
	docker build ${PROXY} -t startx/monkey-ops:${TAG} .
    cd -
fi