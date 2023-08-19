package repository

import (
	"database/sql"

	"github.com/ndfsa/backend-test/internal/model"
)

func GetServices(db *sql.DB, userId uint64) ([]model.Service, error) {
	rows, err := db.Query(`SELECT * FROM GET_USER_SERVICES($1)`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []model.Service
	for rows.Next() {
		var srv model.Service
		if err := rows.Scan(
			&srv.Id,
			&srv.Type,
			&srv.State,
			&srv.InitBalance,
			&srv.DebitBalance,
			&srv.CreditBalance); err != nil {
			return services, err
		}

		services = append(services, srv)
	}

    if err := rows.Err(); err != nil {
		return services, err
	}

	return services, nil
}

func GetService(db *sql.DB, userId uint64, serviceId uint64) (model.Service, error) {
	rows := db.QueryRow(`SELECT * FROM GET_USER_SERVICES($1) WHERE id = $2`, userId, serviceId)

	var service model.Service
	if err := rows.Scan(
		&service.Id,
		&service.Type,
		&service.State,
		&service.InitBalance,
		&service.DebitBalance,
		&service.CreditBalance); err != nil {
		return service, err
	}

    if err := rows.Err(); err != nil {
		return service, err
	}

	return service, nil
}
