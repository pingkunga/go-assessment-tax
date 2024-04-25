FROM golang:1.21.9-alpine3.19 AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN go build -o /bin/app

FROM alpine:3.19.1

COPY --from=build /bin/app /bin

# https://stackoverflow.com/questions/50178013/docker-expose-using-run-time-environment-variables
EXPOSE 8080

RUN adduser -D -u 2024 appuser
USER appuser

CMD [ "/bin/app" ]