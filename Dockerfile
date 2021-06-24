FROM golang:1.16 as builder

# Add Maintainer Info
LABEL maintainer="Bwire Peter <bwire517@gmail.com>"

WORKDIR /
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the local package files to the container's workspace.
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ocr_binary .

FROM ubuntu:20.04

RUN apt-get update -qq
RUN apt-get install -y ca-certificates libtesseract-dev libleptonica-dev  tesseract-ocr-eng

COPY --from=builder /ocr_binary /ocr
COPY --from=builder /migrations /migrations
WORKDIR /

# Run the service command by default when the container starts.
ENTRYPOINT ["/ocr"]

# Document the port that the service listens on by default.
EXPOSE 7012
