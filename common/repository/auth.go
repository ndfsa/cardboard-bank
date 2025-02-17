package repository

import (
	"context"
	"database/sql"

	"github.com/ndfsa/cardboard-bank/common/model"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return AuthRepository{db}
}

func (repo *AuthRepository) Authenticate(
	ctx context.Context, username, password string,
) (model.User, error) {
	row := repo.db.QueryRowContext(ctx,
		`select u.id, u.clearance, u.username, u.password, u.fullname from users u
        where u.username = $1`,
		username)

	var user model.User
	if err := row.Scan(
		&user.Id,
		&user.Clearance,
		&user.Username,
		&user.Passhash,
		&user.Fullname); err != nil {

        // for constant time validation
        user.Validate(password)
        return model.User{}, err
	}

	if err := user.Validate(password); err != nil {
		return model.User{}, err
	}

	return user, nil
}
