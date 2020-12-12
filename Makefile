 
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

.PHONE: push

push:
	@echo "------------------"
	@echo "--> Push Kubera Auth Server images"
	@echo "------------------"
	REPONAME=$(REPONAME) IMGNAME=$(IMGNAME) IMGTAG=$(IMGTAG) BUILD_TYPE=$(BUILD_TYPE) bash ./hack/push
