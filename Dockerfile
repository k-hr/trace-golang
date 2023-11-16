FROM --platform=linux/amd64 golang:1.19
WORKDIR /work
COPY ./trace-app-golang ./
RUN mkdir -p db/migrations
COPY db/migrations db/migrations
CMD ["/bin/sh", "-c", "./trace-app-golang"]



