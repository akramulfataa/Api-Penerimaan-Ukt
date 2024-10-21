FROM golang:1.23

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /api-penerimaan-ukt

EXPOSE 9393

CMD ["/api-penerimaan-ukt"]