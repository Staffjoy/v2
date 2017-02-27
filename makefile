.PHONY: build clean test stage promote

.DEFAULT_GOAL = build

stage:
	bash ci/stage.sh

build:
	make clean 
	bash ci/build.sh

test: build
	bash ci/test.sh

promote:
	bash ci/promote.sh

clean:
	rm -rf vendor/

purge: clean
	rm -rf ~/Library/Application\ Support/Unison/
	git clean -fdx .

build-fmt:
	bash ci/build-fmt.sh

protobuf: build
	bash ci/protobuf.sh

#
# jenkins commands
#

jenkins: clean test

jenkins-stage: clean test stage

jenkins-promote: promote

#
# dev commands
#

dev-sync:
	bash vagrant/unison.sh

dev-k8s-fix:
	bash vagrant/dev-k8s-fix.sh

dev-build: build
	bash ci/dev-build.sh

dev:
	bash vagrant/dev-watch.sh
