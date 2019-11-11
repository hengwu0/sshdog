all:
	go build -o sshd -mod=vendor
fmt:
	for file in `find -name "*.go" `; do gofmt -l -w $file; done
