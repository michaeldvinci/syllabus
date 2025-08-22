FROM golang:1.24 AS build
# RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ls -la && ls -la cmd/
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags '-linkmode external -extldflags "-static"' -o syllabus ./cmd/syllabus

FROM golang:1.24-alpine
ENV SYLLABUS_CONFIG=/config/books.yaml
WORKDIR /app
COPY --from=build /app/syllabus /app/syllabus
COPY --from=build /app/app/res /app/res
ENTRYPOINT /app/syllabus "$SYLLABUS_CONFIG"
