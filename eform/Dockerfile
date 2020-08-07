FROM golang:1.13.8-alpine AS build
COPY ./ /go/src/github.com/eform
WORKDIR /go/src/github.com/eform
RUN go build -o chaincode -v .

FROM alpine:3.11 as prod
ENV ISEXTERNAL=true
COPY --from=build /go/src/github.com/eform/chaincode /app/chaincode
USER 1000
WORKDIR /app
CMD ./chaincode
