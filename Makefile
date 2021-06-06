IMAGE = "mock:latest"

# h - help
h help:
	@echo "h help 	- this help"
	@echo "docker 	- run docker image build"
.PHONY: h


# docker build
docker:
	docker build -f ./build/Dockerfile -t $(IMAGE) .
.PHONY: docker
