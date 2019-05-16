FROM alpine:3.9

RUN apk update && apk add --no-cache \
	bash \
	sudo \
	curl

COPY ./configdownloader /bin/configdownloader
RUN chmod +x /bin/configdownloader

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

WORKDIR /
ENTRYPOINT ["/entrypoint.sh"]
