FROM postgres:13

RUN apt-get update \
    && apt-get install -y \
    postgresql-server-dev-13 \
    build-essential \
    pgxnclient \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
