FROM ubuntu:18.04

RUN apt-get update
RUN apt-get install -y postgresql-client
COPY wait-for-postgres.sh /bin/wait-for-postgres

ENTRYPOINT ["wait-for-postgres"]