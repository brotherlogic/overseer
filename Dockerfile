# syntax=docker/dockerfile:1

FROM golang:1.23.5 AS build

WORKDIR $GOPATH/src/github.com/brotherlogic/overseer

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN mkdir proto
COPY proto/*.go ./proto/

RUN CGO_ENABLED=0 go build -o /overseer

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /overseer /overseer

EXPOSE 8080
EXPOSE 8081

USER nonroot:nonroot

ENTRYPOINT ["/overseer"]