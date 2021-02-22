FROM sc4h/alpine-s6overlay:3.12

ENV \
  PATH="/app/nabarr:${PATH}" \
  APP_CONFIG="/config/config.yml" \
  APP_CACHE="/config/cache" \
  APP_LOG="/config/activity.log" \
  APP_VERBOSITY="0"

# Binary
COPY ["dist/nabarr_${TARGETOS}_${TARGETVARIANT}/nabarr", "/app/nabarr/nabarr"]

# Add root files
COPY ["docker/run", "/etc/services.d/nabarr/run"]

# Volume
VOLUME ["/config"]
