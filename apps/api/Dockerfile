FROM golang:1.26-alpine AS build

WORKDIR /src
COPY apps/api/go.mod apps/api/go.sum* ./
RUN go mod download

COPY apps/api ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/splitr-api ./cmd/api

FROM alpine:3.22

RUN addgroup -S splitr && adduser -S splitr -G splitr
USER splitr
WORKDIR /app
COPY --from=build /out/splitr-api /app/splitr-api

EXPOSE 8080
ENTRYPOINT ["/app/splitr-api"]
