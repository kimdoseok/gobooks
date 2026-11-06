# go build stage
FROM golang:alpine as builder

RUN apk update
RUN apk add --no-cache bash tzdata
RUN addgroup -g 1001 -S app && \
    adduser -u 1001 -S app -G app
ENV GOBIN /build/bin

RUN mkdir /build
COPY . /build
WORKDIR /build

RUN mkdir log
RUN mkdir uploads

RUN go build -o /build/bin/grpcserv

# Final image
FROM alpine
RUN addgroup -g 1001 -S app && \
    adduser -u 1001 -S app -G app
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /build /app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

RUN chown -R app:app /app
WORKDIR /app/bin
USER app
CMD ["./grpcserv"]
