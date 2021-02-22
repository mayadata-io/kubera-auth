
# Makefile for building Kubera Auth Server
# Reference Guide - https://www.gnu.org/software/make/manual/make.html

#
# Internal variables or constants.
# NOTE - These will be executed when any make target is invoked.
#
IS_DOCKER_INSTALLED = $(shell which docker >> /dev/null 2>&1; echo $$?)

#docker info
REPONAME ?= mayadataio
IMGNAME ?= kubera-auth

.PHONY: all
all: deps checks build push

.PHONY: help
help:
	@echo ""
	@echo "Usage:-"
	@echo "\tmake all   -- [default] runs all checks and builds the kubera auth server image"
	@echo ""

.PHONY: deps

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

.PHONY: checks

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

.PHONY: build

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



.PHONE: push

push:
	@echo "------------------"
	@echo "--> Push Kubera Auth Server images"
	@echo "------------------"
	REPONAME=$(REPONAME) IMGNAME=$(IMGNAME) IMGTAG=$(IMGTAG) BUILD_TYPE=$(BUILD_TYPE) bash ./hack/push
