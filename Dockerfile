# Use the official golang image, based on Debian.
FROM golang:1.20.0-buster AS builder

# Create and change to the app directory.
WORKDIR /app

# Copies go.mod and go.sum if present.
COPY go.* ./
# Retrieve application dependencies.
RUN go mod download

# Copy local code to container image.
COPY . ./

# Build the binary.
# -v prints the names of the packages as they are compiled.
# -o sets the binary name and output path.
RUN go build -v -o shopa

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
#    /var/lib/apt/lists/* contains the packages downloaded for install
#    it's safe to delete at the end and it saves space.
FROM debian:buster-slim
RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/shopa /app/shopa

# Run the web service on container startup.
CMD ["/app/shopa"]
