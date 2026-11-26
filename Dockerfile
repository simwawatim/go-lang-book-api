FROM golang:1.22-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server \
 && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
 && adduser -D -u 10001 app
WORKDIR /app

COPY --from=build /out/server /out/migrate ./
COPY entrypoint.sh ./
RUN chmod +x entrypoint.sh

ENV PORT=8743
USER 10001
EXPOSE 8743
ENTRYPOINT ["./entrypoint.sh"]
