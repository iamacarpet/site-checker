BUILDTIME:=$(shell date -u +.%Y%m%d.%H%M%S)

build: linux64 linux32
	echo $(BUILDTIME) >> build/build_times.log

linux64:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.minversion=$(BUILDTIME)" -o build/site-checker

linux32:
	GOOS=linux GOARCH=386 go build -ldflags "-X main.minversion=$(BUILDTIME)" -o build/site-checker-32
