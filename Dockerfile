FROM golang:1.21.9-alpine3.19 AS build-stage

WORKDIR /app

COPY . ./
RUN go mod download

RUN go build -o /bin/app

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...


# Run go vet in the container
FROM build-stage AS run-sec-stage
RUN go vet

# Final image
FROM alpine:3.19.1 AS final-stage

COPY --from=build-stage /bin/app /bin

# https://stackoverflow.com/questions/50178013/docker-expose-using-run-time-environment-variables
EXPOSE 8080

RUN adduser -D -u 2024 appuser
USER appuser

CMD [ "/bin/app" ]