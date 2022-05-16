# CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -tags=containers
# docker build -t pomo/pomo:latest -f containers/Dockerfile .
# docker run --rm -it pomo/pomo

FROM alpine:latest
RUN mkdir /app && adduser -h /app -D pomo
WORKDIR /app
COPY --chown=pomo /pomo/pomo .
CMD ["/app/pomo"]