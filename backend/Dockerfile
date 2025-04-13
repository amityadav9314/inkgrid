FROM golang:1.22.6-alpine
RUN mkdir /opt/goinkgridBE
WORKDIR /opt/goinkgridBE
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8034
CMD ["./main"]