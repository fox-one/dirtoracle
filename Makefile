
.PHONY: build
build:
	sh hack/build.sh local

.PHONY: build-dev
build-dev:
	sh hack/build.sh dev

.PHONY: build-prod
build-prod:
	sh hack/build.sh prod

clean:
	rm -rf mtg.*
