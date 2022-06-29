FROM alpine:latest

ENV GOOSE_FILE /bin/goose
ENV MIGRATIONS /migrations

ARG GOOSE_VERSION
# Add goose binary file -------------------------
ADD https://github.com/pressly/goose/releases/download/v${GOOSE_VERSION}/goose_linux_x86_64 ${GOOSE_FILE}
RUN chmod +x ${GOOSE_FILE}

WORKDIR ${MIGRATIONS}

ARG WAIT_VERSION
# Add docker-compose-wait tool -------------------
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/${WAIT_VERSION}/wait /wait
RUN chmod +x /wait

CMD /wait && ${GOOSE_FILE} postgres "host=postgres user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DBNAME} sslmode=${POSTGRES_SSLMODE}" up