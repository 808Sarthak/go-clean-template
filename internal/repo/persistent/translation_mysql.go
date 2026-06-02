package persistent

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/cache"
	"github.com/evrone/go-clean-template/pkg/mysql"
)

const _defaultEntityCap = 64

// TranslationRepo -.
type TranslationRepo struct {
	*mysql.MySQL
	cache *cache.Cache
}

// NewTranslationRepo -.
func NewTranslationRepo(db *mysql.MySQL, cache *cache.Cache) *TranslationRepo {
	return &TranslationRepo{db, cache}
}

// GetHistory -.
func (r *TranslationRepo) GetHistory(ctx context.Context, userID string) ([]entity.Translation, error) {
	cacheKey := fmt.Sprintf("translation_history:%s", userID)
	var entities []entity.Translation

	if r.cache != nil {
		err := r.cache.Get(ctx, cacheKey, &entities)
		if err == nil {
			return entities, nil
		}
	}

	query, args, err := r.Builder.
		Select("source, destination, original, translation").
		From("history").
		Where("user_id = ?", userID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TranslationRepo - GetHistory - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("TranslationRepo - GetHistory - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities = make([]entity.Translation, 0, _defaultEntityCap)

	for rows.Next() {
		e := entity.Translation{}

		err = rows.Scan(&e.Source, &e.Destination, &e.Original, &e.Translation)
		if err != nil {
			return nil, fmt.Errorf("TranslationRepo - GetHistory - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	if r.cache != nil {
		_ = r.cache.Set(ctx, cacheKey, entities, 10*time.Minute)
	}

	return entities, nil
}

// Store -.
func (r *TranslationRepo) Store(ctx context.Context, userID string, t entity.Translation) error {
	query, args, err := r.Builder.
		Insert("history").
		Columns("user_id, source, destination, original, translation").
		Values(userID, t.Source, t.Destination, t.Original, t.Translation).
		ToSql()
	if err != nil {
		return fmt.Errorf("TranslationRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("TranslationRepo - Store - r.Pool.Exec: %w", err)
	}

	if r.cache != nil {
		_ = r.cache.Delete(ctx, fmt.Sprintf("translation_history:%s", userID))
	}

	return nil
}
