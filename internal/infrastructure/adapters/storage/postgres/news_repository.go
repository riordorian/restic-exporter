package postgres

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
	"restic-exporter/internal/domain/news"
	"restic-exporter/internal/infrastructure/db"
	"restic-exporter/internal/shared/dto"
	"strings"
)

type NewsRepository struct {
	Db *db.Db
}

func (r NewsRepository) GetList(ctx context.Context, req dto.ListRequest) ([]news.News, error) {
	var result []news.News

	query := map[string]interface{}{
		"sort":   req.Sort,
		"status": req.Status,
	}

	author := req.Author

	switch author {
	case uuid.Nil:
		query["author"] = " IS NOT NULL"
	default:
		// TODO: Move stringify to serializer
		query["author"] = author.String()
	}

	queryString := strings.Join(
		[]string{
			"SELECT * FROM news WHERE created_by=:author AND status=:status",
			"ORDER BY created_at",
			req.Sort,
		},
		" ")

	rows, err := r.Db.NamedQuery(ctx, queryString, query)

	if err != nil {
		return nil, err
	}

	var newItem news.News
	for rows.Next() {
		if errScan := rows.StructScan(&newItem); err != nil {
			return nil, errScan
		} else {
			result = append(result, newItem)
		}
	}

	return result, nil
}

func (NewsRepository) GetById(uuid uuid.UUID) (news.News, error) {
	return news.News{}, nil
}

func (r NewsRepository) Insert(ctx context.Context, fields news.News) (uuid.UUID, error) {
	media, err := json.Marshal(fields.Media)
	if err != nil {
		return uuid.Nil, err
	}
	queryFields := map[string]interface{}{
		"title":      fields.Title,
		"text":       fields.Text,
		"created_by": fields.CreatedBy,
		"status":     fields.Status,
		"media":      media,
	}

	var uid uuid.UUID
	sql := "INSERT INTO news (title, text, created_by, status, media) VALUES (:title, :text, :created_by, :status, :media) RETURNING id"
	preparedQuery, err := r.Db.PrepareNamed(sql)
	if err != nil {
		return uuid.Nil, err
	}
	err = preparedQuery.Get(&uid, queryFields)
	if err != nil {
		return uuid.New(), err
	}
	return uid, nil
}

func (NewsRepository) Update(id uuid.UUID, fields news.News) (bool, error) {
	return false, nil
}

func (NewsRepository) Delete(id uuid.UUID) (bool, error) {
	return false, nil
}
