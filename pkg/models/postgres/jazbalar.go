package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"kazybek.net/finalproject/pkg/models"
	"strconv"
	"time"
)

const (
	insertSql                 = "INSERT INTO jazbalar (takyryp,mazmuny,created,expires) VALUES ($1,$2,$3,$4) RETURNING id"
	getJazbaById              = "SELECT id, takyryp, mazmuny, created, expires FROM jazbalar where id=$1 AND expires > now()"
	getLastTenCreatedJazbalar = "SELECT id, takyryp, mazmuny, created, expires FROM jazbalar WHERE expires > now() ORDER BY created DESC LIMIT 10"
)

type JazbaModel struct {
	Pool *pgxpool.Pool
}

func (m *JazbaModel) Insert(takyryp, mazmuny, expires string) (int, error) {
	var id uint64
	interval, err := strconv.Atoi(expires)

	row := m.Pool.QueryRow(context.Background(), insertSql, takyryp, mazmuny, time.Now(), time.Now().AddDate(0, 0, interval))
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *JazbaModel) Get(id int) (*models.Jazba, error) {
	s := &models.Jazba{}
	err := m.Pool.QueryRow(context.Background(), getJazbaById, id).
		Scan(&s.ID, &s.Takyryp, &s.Mazmuny, &s.Created, &s.Expires)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *JazbaModel) Latest() ([]*models.Jazba, error) {
	jazbalar := []*models.Jazba{}
	rows, err := m.Pool.Query(context.Background(), getLastTenCreatedJazbalar)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := &models.Jazba{}
		err = rows.Scan(&s.ID, &s.Takyryp, &s.Mazmuny, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		jazbalar = append(jazbalar, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return jazbalar, nil
}
