FROM golang:1.18-alpine AS build
ARG COMMIT_SHA
ARG VERSION
COPY . $GOPATH/src/app
WORKDIR $GOPATH/src/app
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-X 'main.version=$VERSION' -X 'main.commitSha=$COMMIT_SHA'" \
  -a -installsuffix cgo -o app .

FROM scratch
COPY --from=build /go/src/app/app /bin/app
EXPOSE 8080
CMD ["app"]
