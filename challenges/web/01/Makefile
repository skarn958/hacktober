TARGET_NAME=chall
DOCKER_NAME=hacktober_web_challenge_01
SRCS=main.go

GO=go

all: build

build: $(TARGET_NAME)

$(TARGET_NAME): $(SRCS)
	CGO_ENABLED=0 $(GO) build -o $@ $(SRCS)

run:
	FLAG=this_is_only_a_test_flag ./$(TARGET_NAME)

archive:
	tar -czvf $(TARGET_NAME).tar.gz $(SRCS) Dockerfile Makefile

dockerfile:
	podman build -t $(DOCKER_NAME) .

clean:
	rm -f $(TARGET_NAME) $(TARGET_NAME).tar.gz