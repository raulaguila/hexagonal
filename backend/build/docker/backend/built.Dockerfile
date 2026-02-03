FROM alpine:3.21.3

RUN apk add --no-cache bash tzdata dumb-init

WORKDIR /opt/api

# Copy files to docker image
COPY config/.env config/.env
COPY bin/backend .

RUN chmod +x backend

ENTRYPOINT [ "/usr/bin/dumb-init", "--" ]

CMD [ "./backend" ]
