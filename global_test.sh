#!/bin/bash

# Полное тестирование всех API endpoints

# Базовый URL для API
BASE_URL="http://localhost"
SERVICE_NAME_AUTH="auth"
SERVICE_NAME_NOTES="notes"

# Генерация случайного username
RANDOM_USERNAME="testuser_$(date +%s)_$RANDOM"

echo "🚀 Полное тестирование Auth API"
echo "🎲 Используемый username: $RANDOM_USERNAME"
echo "==============================="

# Тест 1: Регистрация нового пользователя
echo ""
echo "🔍 Регистрация нового пользователя"
echo "Запрос: POST $BASE_URL/$SERVICE_NAME_AUTH/register"
curl -X POST "$BASE_URL/$SERVICE_NAME_AUTH/register" \
     -H "Content-Type: application/json" \
     -d '{
           "username": "'$RANDOM_USERNAME'",
           "password": "password123"
         }' \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 1

# Тест 2: Вход в систему и получение токена
echo ""
echo "🔍 Вход в систему (получение JWT токена)"
echo "Запрос: POST $BASE_URL/$SERVICE_NAME_AUTH/login"
echo "Ответ:"
# Сохраняем ответ логина для извлечения токена
LOGIN_RESPONSE=$(curl -X "POST" "$BASE_URL/$SERVICE_NAME_AUTH/login" \
     -H "Content-Type: application/json" \
     -d '{"username": "'$RANDOM_USERNAME'","password":"password123"}' \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s)

echo "$LOGIN_RESPONSE"

# Извлекаем токен из JSON ответа (поле "access_token")
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

echo "Извлеченный токен: $TOKEN"
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 1

# Тест 3: Получение информации о пользователе
echo ""
echo "🔍 Получение информации о пользователе"
echo "Запрос: GET $BASE_URL/$SERVICE_NAME_AUTH/user"
echo "Ответ:"
curl -X "GET" "$BASE_URL/$SERVICE_NAME_AUTH/user" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 1

# Тест 4: Обновление информации о пользователе
echo ""
echo "🔍 Обновление информации о пользователе"
echo "Запрос: PUT $BASE_URL/$SERVICE_NAME_AUTH/user"
echo "Ответ:"
curl -X "PUT" "$BASE_URL/$SERVICE_NAME_AUTH/user" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"username":"updated_'$RANDOM_USERNAME'","email":"updated@example.com"}' \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 1

# Тест 5: Проверка обновленной информации
echo ""
echo "🔍 Проверка обновленной информации о пользователе"
echo "Запрос: GET $BASE_URL/$SERVICE_NAME_AUTH/user"
echo "Ответ:"
curl -X "GET" "$BASE_URL/$SERVICE_NAME_AUTH/user" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 1

# Тест 6: Удаление пользователя
echo ""
echo "🔍 Удаление пользователя"
echo "Запрос: DELETE $BASE_URL/$SERVICE_NAME_AUTH/user"
echo "Ответ:"
curl -X "DELETE" "$BASE_URL/$SERVICE_NAME_AUTH/user" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s

echo ""
echo "✅ Тестирование сервиса AUTH  завершено!"
echo "-------------------------------------------"

echo ""
echo "🚀 Полное тестирование Notes API"
echo "==============================="


# Тест 1: Создание новой заметки
echo ""
echo "🔍 Создание новой заметки"
echo "Запрос: POST $BASE_URL/$SERVICE_NAME_NOTES/note"
echo "Ответ:"
CREATE_RESPONSE=$(curl -X "POST" "$BASE_URL/$SERVICE_NAME_NOTES/note" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"name":"Test Note","content":"Test Content"}' \
     -w "\n📊 HTTP Статус: %{http_code}\n")
     
# Извлекаем ID созданной заметки из ответа
ID_NOTE=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "ID созданной заметки: $ID_NOTE"
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 2
# Тест 2: Получение списка всех заметок
echo ""
echo "🔍 Получение списка всех заметок"
echo "Запрос: GET $BASE_URL/$SERVICE_NAME_NOTES/notes"
echo "Ответ:"
curl -X "GET" "$BASE_URL/$SERVICE_NAME_NOTES/notes" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 2
# Тест 3: Получение заметки по ID
echo ""
echo "🔍 Получение заметки по ID"
echo "Запрос: GET $BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE"
echo "Ответ:"
curl -X "GET" "$BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 2
# Тест 4: Редактирование заметки по ID
echo ""
echo "🔍 Редактирование заметки по ID"
echo "Запрос: PUT $BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE"
echo "Ответ:"
curl -X "PUT" "$BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"name":"Updated Note","content":"Updated Content"}' \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 2
# Тест 5: Удаление заметки по ID
echo ""
echo "🔍 Удаление заметки по ID"
echo "Запрос: DELETE $BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE"
echo "Ответ:"
curl -X "DELETE" "$BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

echo "✅ Все тесты завершены!"
echo "==============================="