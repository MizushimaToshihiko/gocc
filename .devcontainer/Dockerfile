# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.194.3/containers/ubuntu/.devcontainer/base.Dockerfile

# [Choice] Ubuntu version: hirsute, bionic, focal
ARG VARIANT="focal"
FROM mcr.microsoft.com/vscode/devcontainers/base:0-${VARIANT}
WORKDIR /root/
# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
        build-essential \
        git \
        vim \
        man \
        manpages-dev \
        strace \
        gdb

FROM mcr.microsoft.com/vscode/devcontainers/go:1.17

WORKDIR /go/src/github.com/MizushimaToshihiko/gocc

COPY --from=0 /root/ ./

# Golang 環境構築(任意)
RUN go get golang.org/x/tools/gopls \
  golang.org/x/tools/cmd/godoc

RUN go install github.com/tenntenn/goplayground/cmd/gp@latest
RUN go get -u github.com/kisielk/errcheck