package repository

import (
	"context"
	"database/sql"
	"iter"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/model"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsrRepository(db *sql.DB) UsersRepository {
	return UsersRepository{db}
}

func (repo *UsersRepository) CreateUser(ctx context.Context, user model.User) error {
	if _, err := repo.db.ExecContext(ctx,
		`insert into users(id, clearance, username, password, fullname)
        values ($1, $2, $3, $4, $5)`,
		user.Id,
		user.Clearance,
		user.Username,
		user.Passhash,
		user.Fullname); err != nil {
		return err
	}

	return nil
}

func (repo *UsersRepository) FindUser(ctx context.Context, userId uuid.UUID) (model.User, error) {
	row := repo.db.QueryRowContext(ctx,
		`select id, clearance, username, password, fullname from users
        where id = $1`, userId)

	var user model.User
	if err := row.Scan(
		&user.Id,
		&user.Clearance,
		&user.Username,
		&user.Passhash,
		&user.Fullname); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (repo *UsersRepository) FindAllUsers(
	ctx context.Context,
	cursor uuid.UUID,
) (iter.Seq2[model.User, error], error) {
	query := "select id, clearance, username, password, fullname from users"
	params := make([]interface{}, 0, 1)

	if (cursor != uuid.UUID{}) {
		query += " where id > $1"
        params = append(params, cursor)
	}

	query += " order by id"
	query += " limit 10"

	rows, err := repo.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	it := func(yield func(model.User, error) bool) {
		defer rows.Close()
		for rows.Next() {
			var user model.User
			err := rows.Scan(
				&user.Id,
				&user.Clearance,
				&user.Username,
				&user.Passhash,
				&user.Fullname)

			if !yield(user, err) {
				return
			}
		}
	}

	return it, nil
}

func (repo *UsersRepository) UpdateUser(ctx context.Context, user model.User) error {
	if _, err := repo.db.ExecContext(ctx,
		`update users set
        username = coalesce(nullif($1, ''), username),
        fullname = coalesce(nullif($2, ''), fullname),
        password = coalesce(nullif($3, ''), password)
        where id = $4`, user.Username, user.Fullname, user.Passhash, user.Id); err != nil {
		return err
	}

	return nil
}
