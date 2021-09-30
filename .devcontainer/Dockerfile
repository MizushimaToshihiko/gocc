# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.194.3/containers/ubuntu/.devcontainer/base.Dockerfile

# [Choice] Ubuntu version: hirsute, bionic, focal
ARG VARIANT="focal"
FROM mcr.microsoft.com/vscode/devcontainers/base:0-${VARIANT}

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
        strace

FROM mcr.microsoft.com/vscode/devcontainers/go:1.17

# Golang 環境構築(任意)
RUN go get github.com/go-delve/delve/cmd/dlv \
  golang.org/x/lint/golint \
  golang.org/x/tools/gopls \
  golang.org/x/tools/cmd/godoc

RUN go install github.com/tenntenn/goplayground/cmd/gp@latest
RUN go get -u github.com/kisielk/errcheck