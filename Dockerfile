FROM golang:1.25.1-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /bin/node ./cmd/node

FROM alpine:3.20
COPY --from=build /bin/node /bin/node
EXPOSE 8080
ENTRYPOINT ["/bin/node"]
