
# Makefile for building Kubera Auth Server
# Reference Guide - https://www.gnu.org/software/make/manual/make.html

#
# Internal variables or constants.
# NOTE - These will be executed when any make target is invoked.
#
IS_DOCKER_INSTALLED = $(shell which docker >> /dev/null 2>&1; echo $$?)
GOLANGCI_LINT := $(shell command -v golangci-lint --version 2> /dev/null)
#docker info
REPONAME ?= mayadataio
IMGNAME ?= kubera-auth

all: deps checks build push

help:
	@echo ""
	@echo "Usage:-"
	@echo "\tmake all   -- [default] runs all checks and builds the kubera auth server image"
	@echo ""

deps:
	@echo "------------------"
	@echo "--> Check the Docker deps"
	@echo "------------------"
	@if [ $(IS_DOCKER_INSTALLED) -eq 1 ]; \
		then echo "" \
		&& echo "ERROR:\tdocker is not installed. Please install it before build." \
		&& echo "" \
		&& exit 1; \
		fi;

checks:
	@echo "------------------"
	@echo "--> Check Module Deps [go mod tidy]"
	@echo "------------------"
	@tidyRes=$$(go mod tidy); \
	if [ -n "$${tidyRes}" ]; then \
		echo "go mod tidy checking failed!" && echo "$${tidyRes}" \
		&& echo "Please ensure you are using $$($(GO) version) for formatting code." \
		&& exit 1; \
	fi

build:
	@echo "------------------"
	@echo "--> Build Kubera Auth Server Image"
	@echo "------------------"
	docker build . -f ./Dockerfile -t $(REPONAME)/$(IMGNAME):$(IMGTAG)

test:
	@echo "------------------"
	@echo "--> Running tests"
	go test ./...

coverage:
	@echo "Avoid running this one in your dev setup"
	@echo "------------------"
	@echo "--> Running tests with coverage"
	# TODO: Fix the code or set the envs via a for loop to help write tests
	# @for i in "JWT_SECRET" "ADMIN_USERNAME" "ADMIN_PASSWORD" "CONFIGMAP_NAME" "DB_SERVER" "PORTAL_URL";    do  let $i="dummy" ; done
	@ADMIN_USERNAME="a" ADMIN_PASSWORD="b" CONFIGMAP_NAME="c" DB_SERVER="d" PORTAL_URL="e" go test ./... -cover -coverprofile=coverage.txt -covermode=atomic
	# Cleanup ENVs like a good citizen
	@ADMIN_USERNAME="" ADMIN_PASSWORD="" CONFIGMAP_NAME="" DB_SERVER="" PORTAL_URL=""


golint:
	# TODO: Ditch golint in favour of revive, see https://revive.run for more details
	# This is a nice target to run on dev machine before raising a PR
	@echo "------------------"
	@echo "--> Running GoLint"
	$(eval PKGS := $(shell go list ./... | grep -v /vendor/))
	@touch golint.tmp
	-@for pkg in $(PKGS) ; do \
		echo `golint $$pkg | grep -v "have comment" | grep -v "comment on exported" | grep -v "lint suggestions"` >> golint.tmp ; \
	done
	@grep -Ev "^$$" golint.tmp || true
	@if [ "$$(grep -Ev "^$$" golint.tmp | wc -l)" -gt "0" ]; then \
		rm -f golint.tmp; echo "golint failure\n"; exit 1; else \
		rm -f golint.tmp; echo "golint success\n"; \
	fi


golangci:
	# curl -sSfL https://github.com/golangci/golangci-lint/releases/download/v1.37.1/golangci-lint-1.37.1-linux-amd64.tar.gz | tar zxf -
	golangci-lint run

push:
	@echo "------------------"
	@echo "--> Push Kubera Auth Server images"
	@echo "------------------"
	REPONAME=$(REPONAME) IMGNAME=$(IMGNAME) IMGTAG=$(IMGTAG) BUILD_TYPE=$(BUILD_TYPE) bash ./hack/push

.PHONY: push golangci golint build test coverage all help deps checks
