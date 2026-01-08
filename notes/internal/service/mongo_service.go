package service

import (
	"context"
	"encoding/json"
	"fmt"
	"notes/internal/caching"
	"notes/internal/config"
	"notes/internal/database"
	"notes/internal/errors"
	"notes/internal/models"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoService - реализация интерфейса Service для работы с базой данных MongoDB
type MongoService struct {
	db         *mongo.Client     // Указатель на клиент MongoDB для выполнения операций с базой данных
	collection *mongo.Collection // Коллекция заметок в MongoDB
	caching    *redis.Client     // Клиент Redis для кэширования данных
}

// Проверка, что MongoService реализует интерфейс Service
// Это позволяет гарантировать, что MongoService соответствует контракту интерфейса Service
var _ Service = (*MongoService)(nil)

// NewService - конструктор для создания нового экземпляра MongoService
// Он принимает конфигурацию, используется для подключения к MongoDB
func NewService(cfg *config.Config) (Service, error) {
	db, err := database.NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}
	// Создаем новый клиент Redis для кэширования
	cache, err := caching.NewCaching(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrCacheConnection, err)
	}

	// Получаем коллекцию заметок
	collection := db.Database(cfg.DB_NAME).Collection(cfg.DB_COLLECTION)

	return &MongoService{
		db:         db,
		collection: collection,
		caching:    cache,
	}, nil
}

// Close закрывает соединение с базой данных
func (m *MongoService) Close() error {
	if m.caching != nil {
		if err := m.caching.Close(); err != nil {
			return fmt.Errorf("%w: %v", errors.ErrCacheClose, err)
		}
	}
	return database.CloseDB(m.db, &config.Config{Timeout: 10})
}

// getCacheKey генерирует ключ кэша для заметок автора
func (m *MongoService) getCacheKey(authorID int) string {
	return fmt.Sprintf("notes:author:%d", authorID)
}

// invalidateAuthorCache инвалидирует кэш для конкретного автора
func (m *MongoService) invalidateAuthorCache(authorID int) {
	if m.caching != nil {
		cacheKey := m.getCacheKey(authorID)
		m.caching.Del(cacheKey)
		fmt.Println("Кэш для автора с ID", authorID, "был успешно инвалидирован")
	}
	fmt.Println("Кэш для автора с ID", authorID, "не найден или не был инвалидирован")
}

// cacheNotes сохраняет заметки в кэш
func (m *MongoService) cacheNotes(authorID int, notes []models.Note) {
	if m.caching != nil {
		cacheKey := m.getCacheKey(authorID)
		notesJSON, err := json.Marshal(notes)
		if err == nil {
			// Устанавливаем TTL для кэша (100 минут)
			m.caching.Set(cacheKey, notesJSON, 100*time.Minute)
			fmt.Println("Заметки для автора с ID", authorID, "успешно сохранены в кэш")
		}
	}
}

// getCachedNotes получает заметки из кэша
func (m *MongoService) getCachedNotes(authorID int) ([]models.Note, bool) {
	// Получаем ключ кэша для заметок автора
	cacheKey := m.getCacheKey(authorID)
	// Получаем заметки из кэша
	cachedData, err := m.caching.Get(cacheKey).Result()
	if err != nil {
		fmt.Printf("Ошибка при получении кэша для автора с ID %d: %v\n", authorID, err)
		return nil, false
	}
	// Разбираем JSON в слайс заметок
	var cachedNotes []models.Note
	if err := json.Unmarshal([]byte(cachedData), &cachedNotes); err != nil {
		fmt.Printf("Ошибка при разборе кэша для автора с ID %d: %v\n", authorID, err)
		return nil, false
	}
	// Возвращаем заметки и true, если заметки найдены
	fmt.Println("Заметки для автора с ID", authorID, "успешно получены из кэша")
	return cachedNotes, true
}

// Create создает новую заметку в базе данных
func (m *MongoService) Create(ctx context.Context, note models.Note) (*models.Note, error) {
	// Создаем новый ObjectId для заметки
	result, err := m.collection.InsertOne(ctx, bson.M{
		"name":      note.Name,
		"content":   note.Content,
		"author_id": note.AuthorID,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrNoteCreation, err)
	}

	// Получаем ID созданной заметки
	insertedID := result.InsertedID.(primitive.ObjectID)
	note.ID = insertedID.Hex() // Преобразуем ObjectID в строку

	// Инвалидируем кэш для этого автора
	m.invalidateAuthorCache(note.AuthorID)

	return &note, nil
}

// GetByID получает заметку по идентификатору
func (m *MongoService) GetByID(ctx context.Context, id string) (*models.Note, error) {

	// Преобразуем строку ID в ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrInvalidNoteID, err)
	}

	var note models.Note
	err = m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&note)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, id)
		}
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}

	note.ID = objectID.Hex()

	return &note, nil
}

// GetAll получает все заметки из базы данных для конкретного автора
func (m *MongoService) GetAll(ctx context.Context, authorId int) ([]models.Note, error) {
	// Попытаемся получить заметки из кэша
	if cachedNotes, found := m.getCachedNotes(authorId); found {
		// Если заметки найдены в кэше, возвращаем их
		// Если нет, то продолжаем с запросом к базе данных
		fmt.Println("Заметки получены из кэша")
		return cachedNotes, nil
	}

	// Запрос к базе данных
	filter := bson.M{"author_id": authorId}

	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}
	defer cursor.Close(ctx)

	var notes []models.Note
	for cursor.Next(ctx) {
		var note models.Note
		if err := cursor.Decode(&note); err != nil {
			return nil, fmt.Errorf("%w: %v", errors.ErrDecodeNote, err)
		}
		notes = append(notes, note)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrIterationNotes, err)
	}
	// Если изначально кэш не был найден, сохраняем результат в кэш
	m.cacheNotes(authorId, notes)

	return notes, nil
}

// Update обновляет существующую заметку
func (m *MongoService) Update(ctx context.Context, note models.Note) (*models.Note, error) {
	// Преобразуем строку ID в ObjectID
	objectID, err := primitive.ObjectIDFromHex(note.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrInvalidNoteID, err)
	}

	// Сначала получаем существующую заметку для получения AuthorID
	var existingNote models.Note
	err = m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingNote)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, note.ID)
		}
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}

	update := bson.M{
		"$set": bson.M{
			"name":    note.Name,
			"content": note.Content,
		},
	}

	result, err := m.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrNoteUpdate, err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, note.ID)
	}

	// Инвалидируем кэш
	m.invalidateAuthorCache(existingNote.AuthorID)

	// Устанавливаем AuthorID из существующей заметки
	note.AuthorID = existingNote.AuthorID

	return &note, nil
}

// Delete удаляет заметку по идентификатору
func (m *MongoService) Delete(ctx context.Context, id string) error {
	// Преобразуем строку ID в ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrInvalidNoteID, err)
	}

	// Сначала получаем заметку для получения AuthorID перед удалением
	var existingNote models.Note
	err = m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingNote)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, id)
		}
		return fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}

	result, err := m.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrNoteDeletion, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, id)
	}
	// Инвалидируем кэш отдельной заметки
	m.invalidateAuthorCache(existingNote.AuthorID)

	return nil
}
