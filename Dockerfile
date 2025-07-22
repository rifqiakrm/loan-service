FROM rifqiakrm/grpc-go-base-builder:1.0.3-alpine AS base

# define timezone
ENV TZ Asia/Jakarta

# define work directory
WORKDIR /app

# copy the sourcecode
COPY . .

# build exec
RUN cd /app/cmd && go mod vendor && go build -o backend-service

FROM alpine:3.16

WORKDIR /app

COPY --from=base app/cmd/backend-service .

# EXPOSE 8080 is the port that the REST API will be exposed on
EXPOSE 8080

CMD [ "./backend-service" ]