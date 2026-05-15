# Default environment
ENV := dev

# Intercept CLI-style flags like `--dev` or `--prod`
ifeq ($(firstword $(MAKECMDGOALS)),--dev)
  ENV := dev
  override MAKECMDGOALS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
endif

ifeq ($(firstword $(MAKECMDGOALS)),--stage)
  ENV := stage
  override MAKECMDGOALS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
endif

ifeq ($(firstword $(MAKECMDGOALS)),--prod)
  ENV := prod
  override MAKECMDGOALS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
endif


#Targets
.PHONY: run-dev run-staging run-prod build-dev build-stage build-prod

run-dev:
	$(MAKE) run ENV=dev

run-staging:
	$(MAKE) run ENV=stage

run-prod:
	$(MAKE) run ENV=prod

build-dev:
	$(MAKE) build ENV=dev

build-stage:
	$(MAKE) build ENV=stage

build-prod:
	$(MAKE) build ENV=prod

run: build
	@bin/smeag_sms_manager_module

clear_buff:
	@rm -rf ./pb

build:
	@go build -ldflags "-X 'main.BuildEnv=${ENV}'" -o bin/smeag_sms_manager_module cmd/main.go

gen:
	mkdir pb
	protoc -I=proto --go_out=pb --go-grpc_out=pb proto/*.proto