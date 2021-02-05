#Докер файл для создания образа gprc сервера

# STEP 1 build executable binary
#FROM golang:alpine as builder
# Create appuser
#RUN adduser -D -g '' appuser
#COPY . /home/app/
#WORKDIR /home/app/
#ENV GO111MODULE=on
#build the binary
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /home/app/bin/server -ldflags="-w -s" /home/app/cmd/server/main.go

# STEP 2 build a small image
FROM scratch
#COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
#COPY --from=builder /home/app/bin/server /go/bin/server
#USER appuser
COPY bin/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]