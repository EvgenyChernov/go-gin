# Используем официальный образ Go как базовый
FROM golang:1.25.5-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go mod файлы notes сервиса
COPY notes/go.mod notes/go.sum ./notes/

# Копируем pkg директорию для локального модуля
COPY pkg ./pkg

# Переходим в папку notes
WORKDIR /app/notes

# Загружаем зависимости
RUN go mod download

# Копируем исходный код notes сервиса
COPY notes/ .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный этап - минимальный образ
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS
RUN apk --no-cache add ca-certificates

# Создаем рабочую директорию
WORKDIR /root/

# Копируем исполняемый файл из builder этапа
COPY --from=builder /app/notes/main .

# Открываем порт
EXPOSE 8103

# Запускаем приложение
CMD ["./main"]