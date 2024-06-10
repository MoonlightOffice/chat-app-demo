FROM ubuntu:latest
SHELL ["/bin/bash", "-c"]

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update && apt full-upgrade -y && apt install -y less vim curl git unzip
RUN git clone -b stable https://github.com/flutter/flutter.git /flutter
ENV PATH=$PATH:/flutter/bin
RUN flutter doctor
RUN flutter precache --web

ENV DEBIAN_FRONTEND=dialog
