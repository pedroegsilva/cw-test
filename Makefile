SHELL := bash

IMAGE_TAG=cw-test:local
PWD := $(shell pwd)

run-docker:  
	docker build . -t ${IMAGE_TAG}
	docker run -e OUT_HUMAN=true ${IMAGE_TAG}

run-docker-json-out:  
	docker build . -t ${IMAGE_TAG}
	docker run -e OUT_JSON=true ${IMAGE_TAG}

run:  
	export OUT_HUMAN=true; go run main.go -i ${PWD}/input/qgames.log

run-json:  
	export OUT_JSON=true; go run main.go -i ${PWD}/input/qgames.log	

test:
	go test ./parser