REGISTRY ?= aramase
IMAGE_NAME ?= csi
VERSION ?= v0.3
IMAGE_TAG ?= $(REGISTRY)/$(IMAGE_NAME):$(VERSION)

## --------------------------------------
## Builds
## --------------------------------------
build-binary:
	@echo "** building binary with version:$(VERSION)"
	@GOOS=linux GOARCH=${ARCH} go build -o _output/${ARCH}/csi .

build-image:
	@docker buildx build --no-cache --platform "linux/amd64,linux/arm64" -t $(IMAGE_TAG) --push .

push-image:
	@docker push $(IMAGE_TAG)

build-and-push-image: build-binary build-image push-image ## builds binary, image and pushes image
