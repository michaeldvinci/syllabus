FROM golang:1.24 AS build
WORKDIR /app
COPY . .
RUN go build -o syllabus . && ls -al

FROM golang:1.24
ENV SYLLABUS_CONFIG=/config/books.yaml
WORKDIR /app
COPY --from=build /app/syllabus /app/syllabus
ENTRYPOINT /app/syllabus "$SYLLABUS_CONFIG"
