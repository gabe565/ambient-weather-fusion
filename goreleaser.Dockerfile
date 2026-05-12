FROM gcr.io/distroless/static:nonroot
WORKDIR /
ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/ambient-weather-fusion /
ENTRYPOINT ["/ambient-weather-fusion"]
