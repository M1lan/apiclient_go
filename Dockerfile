FROM golang:1.21 as build
WORKDIR /src
COPY main.go ./main.go
RUN go build -o /bin/apiclient_go ./main.go

FROM scratch
COPY --from=build /bin/apiclient_go /bin/apiclient_go
CMD ["/bin/apiclient_go"]
