FROM golang:1.23.3-alpine3.20 AS build_stage
COPY . /go/src/item_store
WORKDIR /go/src/item_store
RUN go build -o /go/bin/item_store cmd/main.go

FROM alpine AS run_stage
WORKDIR /app_binary
COPY --from=build_stage /go/bin/item_store /app_binary/
RUN chmod +x ./item_store
EXPOSE 8080/tcp
ENTRYPOINT ./item_store

EXPOSE 8080/tcp
CMD [ "item_store" ]