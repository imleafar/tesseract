FROM golang:alpine
LABEL maintainer="Rafael Teixeira <rafaelteixeiradev@gmail.com>"
WORKDIR /go
COPY . .
EXPOSE 3000
RUN apk update && apk add tesseract-ocr && apk add tesseract-ocr-data-por
CMD [ "go", "run", "main.go" ]