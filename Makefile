
EXEC=bsblookup

all: help

help: ## This help.
	@grep -E -h "^[a-zA-Z_-]+:.*?## " $(MAKEFILE_LIST) \
	  | sort \
	  | awk -v width=36 'BEGIN {FS = ":.*?## "} {printf "\033[36m%-*s\033[0m %s\n", width, $$1, $$2}'

## Production

build/$(EXEC): *.go lambda/*.go
	env GOOS=linux GOARCH=amd64 go build -o build/$(EXEC) ./lambda

build: build/$(EXEC).zip ## Build zip file for AWS Lambda.

build/$(EXEC).zip: build/$(EXEC)
	zip -j build/$(EXEC).zip build/$(EXEC)

deploy: build ## Deploy to AWS Lambda.  AWS_PROFILE=timix make deploy
	aws lambda update-function-code --function-name $(EXEC) --zip-file fileb://$$PWD/build/$(EXEC).zip


## Local Dev

$(EXEC): *.go
	docker run --rm \
	-v "$$PWD":/go/src/handler \
	lambci/lambda:build-go1.x \
	sh -c 'go mod download && go build -o lambda/$(EXEC) .'

dev-build: $(EXEC) ## Local build.

dev: dev-build ## Local build and run.
	docker run --rm \
	-e DOCKER_LAMBDA_STAY_OPEN=1 \
	-p 9001:9001 \
	-v "$$PWD":/var/task \
	lambci/lambda:go1.x lambda/$(EXEC) '{"bsb": "012-023"}'

ping:
	curl -d '{}' http://localhost:9001/2015-03-31/functions/$(EXEC)/invocations


.PHONY: clean

clean: ## Remove build files.
	rm -f $(EXEC)
	rm -rf build
