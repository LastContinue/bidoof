# Dockerfile
FROM golang:alpine
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .
RUN adduser -S -D -H -h /app appuser
USER appuser
#This is just to give the other services a chance to start
CMD ["./main"]