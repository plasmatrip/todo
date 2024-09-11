FROM golang:1.22.3

EXPOSE 7540

ENV WEB_DIR ./web/
ENV TODO_PORT 7540
ENV TODO_DB_DIR ./db/
ENV TODO_DBFILE scheduler.db
ENV TODO_DATE_LAYOUT 20060102
ENV APP_LOG_DIR ./log/
ENV APP_LOG_FILE app.log

WORKDIR /usr/src/app

COPY ./api ./api
COPY ./cmd ./cmd 
COPY ./configs ./configs
COPY ./model ./model
COPY ./repository ./repository
COPY ./service ./service
COPY ./web ./web

COPY go.mod go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todo ./cmd/main.go

CMD ["./todo"]