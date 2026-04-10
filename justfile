set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

IMAGE_PREFIX := env_var_or_default("IMAGE_PREFIX", "localhost/borg")
IMAGE_VERSION := env_var_or_default("IMAGE_VERSION", "latest")

TOOLS := "droid siegfried tika magika mediainfo jhove verapdf odf-validator ooxml-validator"

default:
	@just --list

build-gui:
	podman build -t "{{IMAGE_PREFIX}}/gui:{{IMAGE_VERSION}}" ./gui

build-server:
	podman build -t "{{IMAGE_PREFIX}}/server:{{IMAGE_VERSION}}" ./server

build-tools:
	for tool in {{TOOLS}}; do \
		podman build -t "{{IMAGE_PREFIX}}/$$tool:{{IMAGE_VERSION}}" "./tools/$$tool"; \
	done

build-all: build-server build-gui build-tools

push-gui:
	podman push "{{IMAGE_PREFIX}}/gui:{{IMAGE_VERSION}}"

push-server:
	podman push "{{IMAGE_PREFIX}}/server:{{IMAGE_VERSION}}"

push-tools:
	for tool in {{TOOLS}}; do \
		podman push "{{IMAGE_PREFIX}}/$$tool:{{IMAGE_VERSION}}"; \
	done

push-all: push-server push-gui push-tools

build-and-push-all: build-all push-all