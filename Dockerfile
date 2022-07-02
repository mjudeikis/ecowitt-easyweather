ARG TARGET_GOOS
ARG TARGET_GOARCH
FROM golang:1.18 as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    TARGET_GOOS=${TARGET_GOOS} \
    TARGET_GOARCH=${TARGET_GOARCH}
WORKDIR /app

# <- COPY go.mod and go.sum files to the workspace
COPY go.mod .
#COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# COPY the source code as the last step
COPY . .

RUN make build

FROM alpine
RUN apk --update add ca-certificates

COPY --from=builder /app/ingestor /bin

# Create a group and user
ENV USER_ID=1000
ENV GROUP_ID=1000
ENV USER_NAME=ingestor
ENV GROUP_NAME=ingestor

RUN addgroup -g $USER_ID $GROUP_NAME && \
    adduser --shell /sbin/nologin --disabled-password \
    --no-create-home --uid $USER_ID --ingroup $GROUP_NAME $USER_NAME

USER $USER_ID

ENTRYPOINT ["/bin/ingestor"]
