FROM golang:latest AS build-env
RUN mkdir /forum
WORKDIR /forum



COPY . /forum
RUN cd /forum && go mod download && go build -o forum
# Build the binary

FROM debian:buster

RUN mkdir /app
WORKDIR /app
COPY --from=build-env /forum /app

EXPOSE 8080
ENTRYPOINT [ "./forum" ]