FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ls -la && ls -la cmd/
RUN go build -o syllabus ./cmd/syllabus

FROM golang:1.24-alpine
ENV SYLLABUS_CONFIG=/config/books.yaml
WORKDIR /app
COPY --from=build /app/syllabus /app/syllabus
ENTRYPOINT /app/syllabus "$SYLLABUS_CONFIG"
