dc := docker-compose

up:
	$(dc) up -d $${PARAMS}

down:
	$(dc) down

logs:
	@$(dc) logs -f

setup: up deps

deps:
	rm -rf src/vendor
	$(dc) run --rm app sh -c 'cd src && glide install'

deps_get:
	rm -rf src/vendor
	$(dc) run --rm app sh -c "cd src && glide get $${PACKAGE}"

build:
	$(dc) run --rm app sh -c 'go build -o bin/wc2018-slack-bot src/main/main.go'

run:
	$(dc) run --rm app sh -c "\
	WC2018_POLLING_INTERVAL=$${POLL_INT} \
	WC2018_CURRENT_MATCH_THRESHOLD=$${CURR_TH} \
	WC2018_SLACK_TOKEN=$${TOKEN} \
	WC2018_SLACK_CHANNEL=$${CHANNEL} \
	go run src/main/main.go"
