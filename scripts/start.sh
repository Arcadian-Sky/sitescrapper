#!/bin/bash

# Останавливаем и удаляем существующие контейнеры
docker-compose down -v

# Запускаем сервисы через docker-compose
docker-compose up -d

# Ждем, пока сервисы будут готовы
echo "Ожидание запуска сервисов..."
sleep 10

# Проверяем статус
docker-compose ps