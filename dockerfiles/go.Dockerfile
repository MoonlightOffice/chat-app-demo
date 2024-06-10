FROM golang:latest

RUN go install golang.org/x/tools/cmd/godoc@latest && \
	go install golang.org/x/tools/gopls@latest

CMD ["godoc", "-http", "0.0.0.0:6060"]
