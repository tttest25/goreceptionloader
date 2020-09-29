# multi-stage build so that:
#    golang builder is not needed on host
#    golang builder remnants not required in Docker image


#
# builder image
#
FROM golang:1.15-alpine3.12 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o goreceptionloader .
# do not need extldflags set to static, because no external linker (CGO disabled)
#RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o golang-memtest .


#
# generate clean, final image for end users
#
FROM alpine:3.12

# copy golang binary into container
COPY --from=builder /build/goreceptionloader .

# executable
ENTRYPOINT [ "./goreceptionloader" ]
# arguments that can be overridden
# 3Mb, 300 milliseconds between allocation
CMD [ "3", "300" ]

