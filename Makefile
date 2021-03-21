# sourced by https://github.com/octomation/makefiles

.DEFAULT_GOAL = init

.PHONY: init
init:
	@git submodule update --init --recursive

.PHONY: pull
pull:
	@git submodule update --recursive --remote

.PHONY: bench
bench:
	@( \
		cd compare/v4; \
		pkg=retry/research/compare/v4; \
		go test -run=NONE -bench=. $$pkg > $(PWD)/v4.out; \
	)
	@( \
		cd compare/v5; \
		pkg=retry/research/compare/v5/async; \
		go test -run=NONE -bench=. $$pkg > $(PWD)/v5.async.out; \
		pkg=retry/research/compare/v5/sync; \
		go test -run=NONE -bench=. $$pkg > $(PWD)/v5.sync.out; \
	)
	@benchcmp v4.out v5.async.out
	@echo "\n---\n"
	@benchcmp v5.async.out v5.sync.out

.PHONY: load
load:
	@hey \
		-c 2 \
		-z 1m \
		-m GET \
		-t 0 \
		-T 'application/json' \
		-H 'X-Timeout: 10ms' \
		http://localhost:8080
