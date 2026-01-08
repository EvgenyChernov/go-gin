#!/bin/bash

# Простой скрипт для тестирования всех API endpoints
# Автоматически извлекает JWT токен и тестирует все маршруты

# Базовый URL для API
BASE_URL="http://localhost/notes"
AUTH_BASE_URL="http://localhost/auth"

echo "🚀 Полное тестирование Notes API"
echo "==============================="


# Получаем токен из auth API
echo ""
echo "🔍 Вход в систему (получение JWT токена)"
echo "Запрос: POST $AUTH_BASE_URL/login"
echo "Ответ:"
# Сохраняем ответ логина для извлечения токена - отдельно тело и статус
LOGIN_RESPONSE=$(curl -X "POST" "$AUTH_BASE_URL/login" \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser","password":"password123"}' \
     -s)

echo "$LOGIN_RESPONSE"

# Получаем статус отдельно
LOGIN_STATUS=$(curl -X "POST" "$AUTH_BASE_URL/login" \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser","password":"password123"}' \
     -w "%{http_code}" \
     -s -o /dev/null)

echo "📊 HTTP Статус: $LOGIN_STATUS"

# Извлекаем токен из JSON ответа (поле "access_token")
# Используем jq для более надежного парсинга JSON, если недоступен - используем grep
if command -v jq &> /dev/null; then
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token')
else
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
fi
# Удаляем возможные переносы строк и пробелы
TOKEN=$(echo "$TOKEN" | tr -d '\n\r ' | xargs)
echo "Извлеченный токен: $TOKEN"
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 3

# Тест 1: Создание новой заметки
echo ""
echo "🔍 Создание новой заметки"
echo "Запрос: POST $BASE_URL/note"
echo "Ответ:"
# Сохраняем ответ создания заметки для извлечения ID
# Используем sed для удаления последней строки (статус код) - работает на macOS и Linux
TEMP_FILE=$(mktemp)
curl -X "POST" "$BASE_URL/note" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"name":"Test Note","content":"Test Content","author_id":1}' \
     -w "\n%{http_code}" \
     -s > "$TEMP_FILE"

# Извлекаем тело ответа (все кроме последней строки) и статус (последняя строка)
# Используем sed для совместимости с macOS
CREATE_RESPONSE=$(sed '$d' "$TEMP_FILE")
CREATE_STATUS=$(tail -n 1 "$TEMP_FILE")
rm "$TEMP_FILE"

echo "$CREATE_RESPONSE"
echo "📊 HTTP Статус: $CREATE_STATUS"

# Извлекаем ID заметки из JSON ответа (поле "note.id")
# Используем jq для более надежного парсинга JSON, если недоступен - используем grep
if command -v jq &> /dev/null; then
    NOTE_ID=$(echo "$CREATE_RESPONSE" | jq -r '.note.id // empty')
else
    NOTE_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
fi
# Удаляем возможные переносы строк и пробелы
NOTE_ID=$(echo "$NOTE_ID" | tr -d '\n\r ' | xargs)

if [ -z "$NOTE_ID" ] || [ "$NOTE_ID" = "null" ] || [ "$NOTE_ID" = "" ]; then
    echo "⚠️  Не удалось извлечь ID заметки из ответа"
    echo "⚠️  Последующие тесты с ID будут пропущены"
    NOTE_ID=""
else
    echo "📝 Извлеченный ID заметки: $NOTE_ID"
fi
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 2
# Тест 2: Получение списка всех заметок
echo ""
echo "🔍 Получение списка всех заметок"
echo "Запрос: GET $BASE_URL/notes"
echo "Ответ:"
curl -X "GET" "$BASE_URL/notes" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 2
# Тест 3: Получение заметки по ID
if [ -n "$NOTE_ID" ]; then
    echo ""
    echo "🔍 Получение заметки по ID"
    echo "Запрос: GET $BASE_URL/note/$NOTE_ID"
    echo "Ответ:"
    curl -X "GET" "$BASE_URL/note/$NOTE_ID" \
         -H "Authorization: Bearer $TOKEN" \
         -H "Content-Type: application/json" \
         -w "\n📊 HTTP Статус: %{http_code}\n"
    echo "-------------------------------------------"
else
    echo ""
    echo "⏭️  Пропуск теста получения заметки по ID (ID не извлечен)"
    echo "-------------------------------------------"
fi

# Небольшая пауза между запросами
sleep 2
# Тест 4: Редактирование заметки по ID
if [ -n "$NOTE_ID" ]; then
    echo ""
    echo "🔍 Редактирование заметки по ID"
    echo "Запрос: PUT $BASE_URL/note/$NOTE_ID"
    echo "Ответ:"
    curl -X "PUT" "$BASE_URL/note/$NOTE_ID" \
         -H "Authorization: Bearer $TOKEN" \
         -H "Content-Type: application/json" \
         -d '{"name":"Updated Note","content":"Updated Content"}' \
         -w "\n📊 HTTP Статус: %{http_code}\n"
    echo "-------------------------------------------"
else
    echo ""
    echo "⏭️  Пропуск теста редактирования заметки (ID не извлечен)"
    echo "-------------------------------------------"
fi

# Небольшая пауза между запросами
sleep 2
# Тест 5: Удаление заметки по ID
if [ -n "$NOTE_ID" ]; then
    echo ""
    echo "🔍 Удаление заметки по ID"
    echo "Запрос: DELETE $BASE_URL/note/$NOTE_ID"
    echo "Ответ:"
    curl -X "DELETE" "$BASE_URL/note/$NOTE_ID" \
         -H "Authorization: Bearer $TOKEN" \
         -H "Content-Type: application/json" \
         -w "\n📊 HTTP Статус: %{http_code}\n"
    echo "-------------------------------------------"
else
    echo ""
    echo "⏭️  Пропуск теста удаления заметки (ID не извлечен)"
    echo "-------------------------------------------"
fi

echo "✅ Все тесты завершены!"
echo "==============================="