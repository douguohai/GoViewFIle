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
