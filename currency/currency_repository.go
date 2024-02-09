package currency

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type currencyRow struct {
	ID    int64           `json:"id"`
	Title sql.NullString  `json:"title"`
	Code  sql.NullString  `json:"code"`
	Value sql.NullFloat64 `json:"value"`
	ADate *time.Time      `json:"a_date"`
}

func (r *currencyRow) toCurrency() *Currency {
	return &Currency{
		ID:    r.ID,
		Title: r.Title.String,
		Code:  r.Code.String,
		Value: r.Value.Float64,
		ADate: *r.ADate,
	}
}

type currencyRepo struct {
	db        *sql.DB
	tableName string
}

func NewPostgresCurrencyRepo(db *sql.DB) (CurrencyRepository, error) {
	var tableName = "currencies"
	var queries = []string{
		`create table if not exists currencies
(
    id                serial not null constraint currencies_pk primary key,
    title        	  varchar(60) not null,
    code          	  varchar(3) not null,
	value 			  numeric(18,2) not null,
	a_date 			  timestamp with time zone not null
);`,
	}
	var err error
	for _, q := range queries {
		_, err = db.Exec(q)
		if err != nil {
			return nil, err
		}
	}
	return &currencyRepo{db, tableName}, nil
}

func (repository *currencyRepo) Create(currency *Currency) (*Currency, error) {
	err := repository.db.QueryRow(
		"INSERT INTO "+repository.tableName+" (title, code, value, a_date) VALUES ($1, $2, $3, $4) RETURNING id",
		currency.Title,
		currency.Code,
		currency.Value,
		currency.ADate,
	).Scan(
		&currency.ID,
	)
	if err != nil {
		return nil, err
	}
	return currency, nil
}

func (repository *currencyRepo) FindAll(pointers CurrencyPointers) ([]*Currency, error) {
	items := []*Currency{}
	q := "SELECT id, title, code, value, a_date FROM " + repository.tableName
	parts := []string{}
	values := []interface{}{}
	cnt := 0
	if pointers.Code != nil {
		cnt++
		parts = append(parts, "code = $"+strconv.Itoa(cnt))
		values = append(values, *pointers.Code)
	}
	if pointers.ADate != nil {
		cnt++
		parts = append(parts, "a_date = $"+strconv.Itoa(cnt))
		values = append(values, *pointers.ADate)
	}
	if len(values) > 0 {
		q = q + " WHERE "
	}
	q = q + strings.Join(parts, " AND ")
	rows, err := repository.db.Query(q, values...)

	defer rows.Close()

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := &currencyRow{}
		err = rows.Scan(
			&item.ID,
			&item.Title,
			&item.Code,
			&item.Value,
			&item.ADate,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item.toCurrency())
	}
	return items, nil
}

func (repository *currencyRepo) FindByID(id int64) (*Currency, error) {
	item := &currencyRow{}
	err := repository.db.QueryRow(
		"SELECT id, title, code, value, a_date FROM "+repository.tableName+" WHERE id = $1 LIMIT 1",
		id,
	).Scan(
		&item.ID,
		&item.Title,
		&item.Code,
		&item.Value,
		&item.ADate,
	)
	//TODO: switch-case
	if err == sql.ErrNoRows {
		return nil, ErrCurrencyNotFound
	} else if err != nil {
		return nil, err
	}
	return item.toCurrency(), nil
}

func (repository *currencyRepo) Update(upd CurrencyPointers) error {
	q := "UPDATE " + repository.tableName + " SET "
	parts := []string{}
	values := []interface{}{}
	cnt := 0
	if upd.Title != nil {
		cnt++
		parts = append(parts, "title = $"+strconv.Itoa(cnt))
		values = append(values, *upd.Title)
	}
	if upd.Code != nil {
		cnt++
		parts = append(parts, "code = $"+strconv.Itoa(cnt))
		values = append(values, *upd.Code)
	}
	if upd.Value != nil {
		cnt++
		parts = append(parts, "value = $"+strconv.Itoa(cnt))
		values = append(values, *upd.Value)
	}
	if upd.ADate != nil {
		cnt++
		parts = append(parts, "a_date = $"+strconv.Itoa(cnt))
		values = append(values, *upd.ADate)
	}
	if len(parts) <= 0 {
		return ErrNothingToUpdate
	}
	cnt++
	q = q + strings.Join(parts, " , ") + " WHERE id = $" + strconv.Itoa(cnt)
	values = append(values, upd.ID)
	ret, err := repository.db.Exec(
		q,
		values...,
	)
	if err != nil {
		return err
	}
	n, err := ret.RowsAffected()
	if err != nil {
		return err
	}
	if n <= 0 {
		return ErrCurrencyNotFound
	}
	return nil
}

func (repository *currencyRepo) Delete(id int64) error {
	ret, err := repository.db.Exec("DELETE FROM "+repository.tableName+" WHERE id = $1", id)
	if err != nil {
		return err
	}
	n, err := ret.RowsAffected()
	if err != nil {
		return err
	}
	if n <= 0 {
		return ErrCurrencyNotFound
	}
	return nil
}
