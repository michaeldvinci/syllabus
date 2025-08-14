FROM golang:1.24 AS build
WORKDIR /app
COPY . .
RUN go build -o syllabus . && ls -al

FROM golang:1.24
WORKDIR /app
COPY --from=build /app/syllabus /app/syllabus
COPY books.yaml /app/books.yaml
# EXPOSE 8080
ENTRYPOINT ["/app/syllabus", "/app/books.yaml"]