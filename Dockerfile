FROM golang:1.25.3 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app
RUN mkdir /schemes

FROM gcr.io/distroless/static-debian12
COPY --from=build /go/bin/app /
COPY --from=build /schemes /schemes

EXPOSE 8080
CMD ["/app", "-d", "/schemes", "-p", "8080"]
