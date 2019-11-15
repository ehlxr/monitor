BUILD_VERSION   := $(shell cat version)
BUILD_TIME		:= $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD)
ROOT_DIR    	:= $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DIST_DIR 		:= $(ROOT_DIR)/dist/

#VERSION_PATH	:= $(shell cat `go env GOMOD` | awk '/^module/{print $$2}')/
VERSION_PATH	:= main
LD_APP_NAMW		:= -X '$(VERSION_PATH).AppName=$(shell basename `pwd`)'
LD_GIT_COMMIT	:= -X '$(VERSION_PATH).GitCommit=$(COMMIT_SHA1)'
LD_BUILD_TIME	:= -X '$(VERSION_PATH).BuildTime=$(BUILD_TIME)'
LD_GO_VERSION	:= -X '$(VERSION_PATH).GoVersion=`go version`'
LD_VERSION		:= -X '$(VERSION_PATH).Version=$(BUILD_VERSION)'
LD_FLAGS		:= "$(LD_APP_NAMW) $(LD_GIT_COMMIT) $(LD_BUILD_TIME) $(LD_GO_VERSION) $(LD_VERSION) -w -s"

RELEASE_VERSION = $(version)
REGISTRY_URL    = $(url)

ifeq ("$(RELEASE_VERSION)","")
	RELEASE_VERSION	:= $(shell echo `date "+%Y%m%d_%H%M%S"`)
endif

.PHONY : build release clean install upx docker-push docker

build:
ifneq ($(shell type gox >/dev/null 2>&1;echo $$?), 0)
	@echo "Can't find gox command, will start installation..."
	cd ~ && go get -v -u github.com/mitchellh/gox && cd $(ROOT_DIR)
endif
	@# $(if $(findstring 0,$(shell type gox >/dev/null 2>&1;echo $$?)),,echo "Can't find gox command, will start installation...";GO111MODULE=off go get -v -u github.com/mitchellh/gox)
	gox -ldflags $(LD_FLAGS) -osarch="darwin/amd64 linux/386 linux/amd64 windows/amd64" \
		-output="$(DIST_DIR){{.Dir}}_{{.OS}}_{{.Arch}}"

docker: build upx
ifneq ("$(REGISTRY_URL)","")
	@echo ========== current docker tag is: $(RELEASE_VERSION) ==========

	docker build -t $(REGISTRY_URL)/monitor_server:$(RELEASE_VERSION) -f Dockerfile .
else
	@echo "url arg should not be empty"
endif

docker-push: docker
	docker push $(REGISTRY_URL)/monitor_server:$(RELEASE_VERSION)

clean:
	rm -rf $(DIST_DIR)*

install:
	go install -ldflags $(LD_FLAGS)

# 如果一个规则是以 .IGNORE 作为目标的，那么这个规则中所有命令都将会忽略错误
.IGNORE:
	upx

# 压缩。需要安装 https://github.com/upx/upx
upx:
	@# 在命令前面加上 "-"，表示不管该命令出不出错，后面的命令都将继续执行下去
	@# -upx $(DIST_DIR)**
	upx $(DIST_DIR)**

release: build upx
ifneq ($(shell type ghr >/dev/null 2>&1;echo $$?), 0)
	@echo "Can't find ghr command, will start installation..."
	cd ~ && go get -v -u github.com/tcnksm/ghr && cd $(ROOT_DIR)
endif
	@# $(if $(findstring 0,$(shell type ghr >/dev/null 2>&1;echo $$?)),,echo "Can't find ghr command, will start installation...";GO111MODULE=off go get -v -u github.com/tcnksm/ghr)
	ghr -u ehlxr -t $(GITHUB_RELEASE_TOKEN) -replace -delete --debug ${BUILD_VERSION} $(DIST_DIR)

# this tells 'make' to export all variables to child processes by default.
.EXPORT_ALL_VARIABLES:

GO111MODULE = on
GOPROXY = https://goproxy.cn,direct
GOSUMDB = sum.golang.google.cn
