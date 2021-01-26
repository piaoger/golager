
OS := $(shell uname -s)
THIS_FOLDER := $(shell cd ${0%/*} && echo ${PWD})

TIMESTAMP = $(shell date "+%G%m%d%H%M%S")

export PATH :=${THIS_FOLDER}/bin:${PATH}

# --------------------------------------------
# golang
# --------------------------------------------
export GOROOT :=${THIS_FOLDER}/bin/apps/golang
export GOPATH :=${THIS_FOLDER}/vendor
export PATH :=${THIS_FOLDER}/bin/apps/golang/${OS}:${PATH}

GO := ${THIS_FOLDER}/bin/apps/golang/${OS}/go
GOFMT :=${THIS_FOLDER}/bin/apps/golang/${OS}/gofmt
GOFILES :=$(wildcard ${THIS_FOLDER}/src/*.go ${THIS_FOLDER}/src/**/*.go)

# vendoring go packages
# for golang.org deps, please add to vendor/golang.org manually..
# golang.com/x/net is added during golang bootstrapping
GO_DEPS := \
	"github.com/BurntSushi/toml" \
	"github.com/kardianos/osext" \
	"github.com/aws/aws-sdk-go" \
	"github.com/aliyun/aliyun-oss-go-sdk/oss" \
	"qiniupkg.com/api.v7" \
	"github.com/qiniu/api.v6" \
	"github.com/piaoger/obs-sdk-go"
.PHONY: build


all: bootstrap build test

help:
	@echo "#-------------------------------------"
	@echo "# PROJECT LAGER GO!"
	@echo "#-------------------------------------"
	@echo " bootstrap       - setup devel environment"
	@echo " build           - build the project"
	@echo " clean           - cleanup and refresh"
	@echo " godeps-list     - list dependencies"
	@echo " godeps-clean    - clean dependencies"
	@echo " godeps-update   - update dependencies"
	@echo " godeps-upgrade  - upgrade dependencies"
	@echo " test            - run go test"
	@echo " dockerize       - create docker container"
	@echo " makebundle      - build package"


bootstrap:
	@echo 'provision golang ...'
	@${THIS_FOLDER}/bin/provision/golang.sh

build:
	@rm -rf build
	@echo 'golang code auto formatting ...'
	@${GOFMT} -s -w ${GOFILES}
	@echo 'building lager tools ...'
	@${GO} build -o build/bin/lager-cli ./src/lager-cli
	@${GO} build -o build/bin/lager-s3 ./src/lager-cli.s3
	@${GO} build -o build/bin/lager-oss ./src/lager-cli.oss
	@${GO} build -o build/bin/lager-qiniu ./src/lager-cli.qiniu
	@${GO} build -o build/bin/lager-obs ./src/lager-cli.obs

build-obs:
	@echo 'building obs tools ...'
	@${GO} build -o build/bin/lager-obs ./src/lager-cli.obs

fmt:
	@echo 'golang code auto formatting ...'
	@${GOFMT} -s -w ${GOFILES}

clean:
	@echo 'clean up and refresh ...'
	@rm -rf build
	@rm -rf bundle

test:
	@echo "TODO: health checking ..."

makebundle:
	@echo "build package ..."
	@rm -rf ${THIS_FOLDER}/bundle && \
	 mkdir -p ${THIS_FOLDER}/bundle/lager-${OS}-${TIMESTAMP} && \
	 cp -r ${THIS_FOLDER}/build/*  ${THIS_FOLDER}/bundle/lager-${OS}-${TIMESTAMP} &&  \
	 cd ${THIS_FOLDER}/bundle && \
	 zip -r lager-${OS}-${TIMESTAMP}.zip lager-${OS}-${TIMESTAMP}

godeps-update:
	@echo 'update dependencies ...'
	@for godep in ${GO_DEPS}; do \
		${GO} get -v -d $${godep}; \
	done;

	@# remove hidden .git stuffs, so that they can be put into git repo
	@echo 'deps are disconnectted from git ...'
	@bash -c "find vendor -name '.git'       -and -type d|xargs rm -rf"
	@bash -c "find vendor -name '.gitignore' -and -type f|xargs rm -f"

godeps-clean:
	@echo 'clean dependencies ...'
	@rm -rf vendor/src/github.com vendor/src/gopkg.in vendor/pkg

godeps-upgrade: godeps-clean godeps-update

godeps-list:
	@${GO} list -f '{{join .Deps "\n"}}' ./... | xargs ${GO} list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'


