# Binary name
BINARY=go_view_file
# Builds the project
build:
		go build -o ${BINARY}
		go test -v
# Installs our project: copies binaries
release-local:
		goreleaser --snapshot --skip-publish --snapshot --rm-dist

release:
		goreleaser

docker-build-base:
		docker buildx build --platform linux/amd64,linux/arm64 -f Dockerfile_Base3 -t douguohai/goviewfile-base:v1 . --push
