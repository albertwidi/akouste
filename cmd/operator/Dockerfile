FROM debian:latest

RUN apt-get update \
    && yes Y | apt-get install curl \
    && curl https://releases.hashicorp.com/consul-template/0.20.0/consul-template_0.20.0_linux_amd64.tgz | tar -xzv -C /usr/local/bin \
    && mkdir /template

COPY ./template/* /template/
COPY appoperator /usr/local/bin

CMD ["consul-template"]