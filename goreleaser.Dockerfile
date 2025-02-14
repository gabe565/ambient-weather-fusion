FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY ambient-weather-fusion /
ENTRYPOINT ["/ambient-weather-fusion"]
