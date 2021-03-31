REGISTRY ?= aramase
IMAGE_NAME ?= csi
VERSION ?= v0.1
IMAGE_TAG ?= $(REGISTRY)/$(IMAGE_NAME):$(VERSION)

## --------------------------------------
## Builds
## --------------------------------------
build-binary:
	@echo "** building binary with version:$(VERSION)"
	@GOOS=linux go build -o _output/csi .

build-image:
	@docker build -t $(IMAGE_TAG) .

push-image:
	@docker push $(IMAGE_TAG)

build-and-push-image: build-binary build-image push-image ## builds binary, image and pushes image
