# Base image
FROM golang:1.23-alpine

# Working directory
WORKDIR /app

# Go Modulesの依存関係をキャッシュするためにまずgo.modとgo.sumをコピー
COPY go.mod .
COPY go.sum .
RUN go install github.com/air-verse/air@latest
RUN go mod download

# Copy the source code
COPY . .

# Expose the app on port 8080
EXPOSE 8080

# Run the app
CMD ["air"]
