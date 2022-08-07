PROJECT_DIR:=${CURDIR}
DOCKER_FILE:=${PROJECT_DIR}
IMAGE:=adershrp85/kube-app

all: docker-push

.PHONY: build
build:
	$(MAKE) -C ${PROJECT_DIR}/src/ manager

.PHONY: docker-build
docker-build:
	docker build -t ${IMAGE} ${DOCKER_FILE}

.PHONY: docker-push
docker-push: build docker-build
	docker push ${IMAGE}

.PHONY: deploy
deploy:
	kubectl delete --ignore-not-found=true -f deploy
	kubectl create -f deploy
