# Builder Stage: Build the Go binary for Linux
FROM golang:1.22 AS builder
WORKDIR /CTS-BE-24

# Copy go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire app and build the Go binary
COPY . .

# Copy the necessary YAML files into the correct location
COPY internal/files/questions.yaml internal/files/
COPY internal/files/answer.yaml internal/files/
COPY internal/files/location.yaml internal/files/

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api/main.go

# Production Stage
FROM alpine:latest
WORKDIR /root/

# Copy the built Go application from the builder stage
COPY --from=builder /CTS-BE-24/main .

# Copy the YAML files from the builder stage to the production image
COPY --from=builder /CTS-BE-24/internal/files/questions.yaml /root/internal/files/
COPY --from=builder /CTS-BE-24/internal/files/answer.yaml /root/internal/files/
COPY --from=builder /CTS-BE-24/internal/files/location.yaml /root/internal/files/

# Copy the .env file to the production image
COPY --from=builder /CTS-BE-24/.env .env

# Ensure binary has execute permissions
RUN chmod +x /root/main

# Set environment variables
ENV PORT=8000
ENV MONGO_URI="mongodb+srv://ishaanagarwal2022:Hg2KRVCQL711ohfD@cluster0.xbrih.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
ENV MONGO_DBNAME=C2S
ENV MONGO_USER_COLLECTION="User"
ENV MONGO_QUESTIONS_COLLECTION="Questions"
ENV MONGO_ROOMS_COLLECTION="Rooms"
ENV JWTSECRET="Ishaan"
ENV JWTEXPINSEC=21600

# Expose port
EXPOSE 8000

# Command to run the Go app
CMD ["./main"]
