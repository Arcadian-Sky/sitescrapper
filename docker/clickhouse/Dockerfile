FROM clickhouse/clickhouse-server:latest

# Копируем SQL-скрипты для инициализации базы данных
COPY migrations/001_create_products_table.sql /docker-entrypoint-initdb.d/

# Создаем скрипт инициализации
RUN echo "CREATE DATABASE IF NOT EXISTS scrapper;" > /docker-entrypoint-initdb.d/000_create_database.sql

# Настройки для ClickHouse
ENV CLICKHOUSE_USER=admin
ENV CLICKHOUSE_PASSWORD=admin123
ENV CLICKHOUSE_DB=scrapper 