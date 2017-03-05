FROM golang:onbuild
EXPOSE 8080
ENTRYPOINT /src/main.go
