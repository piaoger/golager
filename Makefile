
OS := $(shell uname -s)
THIS_FOLDER := $(shell cd ${0%/*} && echo ${PWD})
NOW = $(shell date "+%G%m%d%H%M%S")

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
# e.g, golang.com/x/net is added during golang bootstrapping
GO_EXTDEPS := \
	"github.com/BurntSushi/toml" \
	"github.com/kardianos/osext" \
	"github.com/aws/aws-sdk-go" \
	"github.com/aliyun/aliyun-oss-go-sdk/oss" \
	"qiniupkg.com/api.v7"

.PHONY: build

help:
	@echo "#-------------------------------------"
	@echo "# PROJECT LAGER"
	@echo "#-------------------------------------"
	@echo " bootstrap    - setup devel environment"
	@echo " build        - build the project"
	@echo " clean        - cleanup and refresh"
	@echo " update       - update dependencies"
	@echo " upgrade      - upgrade dependencies"
	@echo " test         - run go test"
	@echo " dockerize    - create docker container"
	@echo " makebundle   - build package"

bootstrap:
	@echo 'provision golang ...'
	@${THIS_FOLDER}/bin/provision/golang.sh

build:
	@rm -rf build
	@echo 'golang code auto formatting ...'
	@${GOFMT} -s -w ${GOFILES}
	@echo 'building lager-cli ...'
	@${GO} build -o build/bin/lager-cli ./src/lager-cli

clean:
	@echo 'clean up and refresh ...'
	@rm -rf build
	@rm -rf bundle


update:
	@echo 'update dependencies ...'
	@for godep in ${GO_EXTDEPS}; do \
		${GO} get -v -d $${godep}; \
	done ;

	@# remove hidden .git stuffs, so that they can be put into git repo
	@echo 'deps are disconnectted from git ...'
	@bash -c "find vendor -name '.git'       -and -type d|xargs rm -rf"
	@bash -c "find vendor -name '.gitignore' -and -type f|xargs rm -f"

upgrade:
	@echo 'clean dependencies ...'
	@rm -rf vendor/src/github.com vendor/src/gopkg.in vendor/pkg

	@echo 'update dependencies ...'
	@for godep in ${GO_EXTDEPS}; do \
		${GO} get -v -d $${godep}; \
	done ;

	# remove hidden .git stuffs, so that they can be put into git repo
	@echo 'deps are disconnectted from git ...'
	@bash -c "find vendor -name '.git'       -and -type d|xargs rm -rf"
	@bash -c "find vendor -name '.gitignore' -and -type f|xargs rm -f"

test:
	@echo "TODO: health checking ..."