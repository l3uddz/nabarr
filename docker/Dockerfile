FROM sc4h/alpine-s6overlay:3.15

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

ENV \
  PATH="/app/nabarr:${PATH}" \
  APP_CONFIG="/config/config.yml" \
  APP_CACHE="/config/cache" \
  APP_LOG="/config/activity.log" \
  APP_VERBOSITY="0"

# Binary
COPY ["dist/nabarr_${TARGETOS}_${TARGETARCH}${TARGETVARIANT:+_7}/nabarr", "/app/nabarr/nabarr"]

# Add root files
COPY ["docker/etc", "/etc"]

# Volume
VOLUME ["/config"]
