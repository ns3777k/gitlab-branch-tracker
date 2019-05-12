FROM golang:1.12.5-alpine3.9 as build
WORKDIR /go/src/github.com/ns3777k/gitlab-branch-tracker
ENV CGO_ENABLED=0 \
	GO111MODULE=on
RUN apk add --no-cache git make

COPY . .
RUN make build

FROM alpine:3.9
COPY --from=build /go/src/github.com/ns3777k/gitlab-branch-tracker/gitlab-branch-tracker /gitlab-branch-tracker
ENTRYPOINT ["/gitlab-branch-tracker"]
CMD ["help"]
