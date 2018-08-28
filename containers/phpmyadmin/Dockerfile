FROM phpmyadmin/phpmyadmin:4.8.3

RUN apk add --no-cache curl curl-dev bash vim

HEALTHCHECK --interval=5s --retries=5 CMD curl --fail http://localhost || exit 1

