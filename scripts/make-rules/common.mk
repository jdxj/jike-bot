FILENAME := jike-bot.out
DOCKER := docker
OUTPUT := output

GIT_TAG := $(shell git describe --tags --abbrev=0)
DOCKER_TAG := jdxj/jike-bot:$(GIT_TAG)

.PHONY: clean
clean:
	@rm -rf $(OUTPUT)
	@rm -rf $(DOCKER)/$(FILENAME)