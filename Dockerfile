FROM golang:alpine as build

RUN apk add --update gcc musl-dev make

WORKDIR /app
COPY . .
RUN make build

FROM alpine:3.10

COPY --from=build /app/bin/shopping-list /shopping-list
ENTRYPOINT [ "/shopping-list" ]
CMD [ "--dbPath", "/shopping.sqlite" ]

