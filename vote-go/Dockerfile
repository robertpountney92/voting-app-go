FROM golang:1.16-buster AS build

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o vote-go

##
## Deploy
##

FROM alpine:3.14

WORKDIR /vote-go

RUN adduser -D appuser && chown -R appuser ./
USER appuser

COPY --from=build /app/vote-go ./
COPY --from=build /app/templates ./templates
COPY --from=build /app/static ./static


EXPOSE 8080

ENTRYPOINT ["./vote-go"]