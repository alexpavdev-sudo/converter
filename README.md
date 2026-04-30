# О проекте

Конвертер файлов

# Сборка и запуск

1. Скопировать `.env.dist` в `.env` и заполнить переменные нужными значениями. Этот файл содержит переменные
   для `docker-compose`.
2. Сделайте файл `generate-ssl-cert.sh` исполняемым и запустити для генерации сертификата https
3. Запустить приложение с помощью `docker-compose up -d`
4. Провести миграции базы данных: `make migrate-up`
5. Открыть в браузере `https://localhost`

# Структура docker-compose

* `docker-compose.yml` - основной файл, содержит настройки для prod-режима.
* `docker-compose.dev.yml` - настройки для запуска в dev-режиме

Для того, чтобы изменить режим, необходимо в `.env` указать `COMPOSE_FILE` и перечислить используемые compose-файлы:

Пример для разработки:

```
COMPOSE_FILE="docker-compose.yml:docker-compose.dev.yml"
```

Пример для prod:

```
COMPOSE_FILE="docker-compose.yml"
```