FROM --platform=linux/amd64 golang:1.19
WORKDIR /work
COPY ./trace-app-golang ./
CMD ["/bin/sh", "-c", "./trace-app-golang"]



