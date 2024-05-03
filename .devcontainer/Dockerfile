FROM ubuntu:22.04

ARG NVM_VERSION=v0.39.7
ARG GOLANG_VERSION="1.22.2"
#ARG NODE_VERSION="22.0.0"

ENV NVM_DIR /usr/local/nvm
ENV NODE_VERSION="22.0.0"
ENV GOPATH="/go"
ENV PATH="${PATH}:/usr/local/go/bin:/go/bin"


RUN mkdir -p $NVM_DIR

RUN apt-get update && apt-get upgrade -y
RUN apt-get update && apt-get install --no-install-recommends -y \
  ca-certificates \
  curl \
  sudo \ 
  build-essential \
  git

# Install nodejs 

RUN curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash - && \
sudo apt-get install -y nodejs 
# Install nvm 
# RUN curl -fsSL "https://raw.githubusercontent.com/nvm-sh/nvm/${NVM_VERSION}/install.sh" | bash 
# # 

# RUN /bin/bash -c "source $NVM_DIR/nvm.sh && nvm install $NODE_VERSION && nvm use --delete-prefix $NODE_VERSION" 


#Install yarn 
RUN  npm install  --global yarn

# Install go 
RUN curl -fsSLO https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local/ -xf  go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm go${GOLANG_VERSION}.linux-amd64.tar.gz


    