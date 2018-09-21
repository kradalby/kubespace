FROM golang:1.11.0-stretch as builder

RUN mkdir -p /app/bin
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/kubespace
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/macos/kubespace
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/windows/kubespace

FROM nginx:alpine 
RUN sed -i '9iautoindex on;' /etc/nginx/conf.d/default.conf
RUN rm -rf /usr/share/nginx/html/*
COPY --from=builder /app/bin /usr/share/nginx/html
