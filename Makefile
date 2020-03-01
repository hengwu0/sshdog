ifeq ($(OS),Windows_NT)     # is Windows_NT on XP, 2000, 7, Vista, 10...
    a.out := sshd.exe
else
    a.out := sshd
endif

all:
	go build -o $(a.out) -ldflags "-s -w" -mod=vendor

fmt:
	for file in `find -name "*.go" `; do gofmt -l -w $$file; done
