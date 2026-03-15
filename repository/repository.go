package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"practice5/model"
)

var allowedOrderBy = map[string]bool{
	"id":         true,
	"name":       true,
	"email":      true,
	"gender":     true,
	"birth_date": true,
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetPaginatedUsers(f model.FilterParams) (model.PaginatedResponse, error) {
	args := []interface{}{}
	argIdx := 1

	where := []string{}

	if f.ID != nil {
		where = append(where, fmt.Sprintf("id = $%d", argIdx))
		args = append(args, *f.ID)
		argIdx++
	}
	if f.Name != nil && *f.Name != "" {
		where = append(where, fmt.Sprintf("name ILIKE $%d", argIdx))
		args = append(args, "%"+*f.Name+"%")
		argIdx++
	}
	if f.Email != nil && *f.Email != "" {
		where = append(where, fmt.Sprintf("email ILIKE $%d", argIdx))
		args = append(args, "%"+*f.Email+"%")
		argIdx++
	}
	if f.Gender != nil && *f.Gender != "" {
		where = append(where, fmt.Sprintf("gender = $%d", argIdx))
		args = append(args, *f.Gender)
		argIdx++
	}
	if f.BirthDate != nil && *f.BirthDate != "" {
		where = append(where, fmt.Sprintf("birth_date = $%d", argIdx))
		args = append(args, *f.BirthDate)
		argIdx++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	orderBy := "id"
	if f.OrderBy != "" {
		col := strings.ToLower(f.OrderBy)
		if allowedOrderBy[col] {
			orderBy = col
		}
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM users %s`, whereClause)
	var totalCount int
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return model.PaginatedResponse{}, fmt.Errorf("count query: %w", err)
	}

	offset := (f.Page - 1) * f.PageSize
	dataQuery := fmt.Sprintf(
		`SELECT id, name, email, gender, birth_date FROM users %s ORDER BY %s LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, argIdx, argIdx+1,
	)
	args = append(args, f.PageSize, offset)

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return model.PaginatedResponse{}, fmt.Errorf("data query: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return model.PaginatedResponse{}, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return model.PaginatedResponse{}, err
	}

	return model.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       f.Page,
		PageSize:   f.PageSize,
	}, nil
}

func (r *Repository) GetCommonFriends(userID1, userID2 int) ([]model.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM users u
		JOIN user_friends uf1 ON uf1.friend_id = u.id AND uf1.user_id = $1
		JOIN user_friends uf2 ON uf2.friend_id = u.id AND uf2.user_id = $2
		ORDER BY u.id
	`
	rows, err := r.db.Query(query, userID1, userID2)
	if err != nil {
		return nil, fmt.Errorf("common friends query: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
