FROM golang:1.13 AS build-env

WORKDIR /blockpropeller

ADD go.mod /blockpropeller/go.mod
ADD go.sum /blockpropeller/go.sum
RUN go mod download

ADD . /blockpropeller

RUN go build -a -o blockpropeller-api ./blockpropeller/cmd/blockpropeller-api

FROM golang:1.13

RUN echo "deb http://ppa.launchpad.net/ansible/ansible/ubuntu trusty main" >> /etc/apt/sources.list && \
    apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 93C4A3FD7BB9C367 && \
    apt update && \
    apt install -y ansible unzip

RUN wget https://releases.hashicorp.com/terraform/0.12.12/terraform_0.12.12_linux_amd64.zip && \
    unzip terraform_0.12.12_linux_amd64.zip && \
    mkdir -p /usr/local/bin && \
    mv terraform /usr/local/bin/ && \
    rm terraform_0.12.12_linux_amd64.zip

WORKDIR /blockpropeller

COPY --from=build-env /blockpropeller/blockpropeller-api .
COPY --from=build-env /blockpropeller/playbooks /blockpropeller/playbooks

ENTRYPOINT ./blockpropeller-api
