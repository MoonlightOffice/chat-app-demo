FROM ubuntu:latest

RUN apt update && apt full-upgrade -y && apt install -y ca-certificates
COPY ./backend/app /backend/app

CMD ["/backend/app"]
