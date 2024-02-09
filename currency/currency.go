package currency

import (
	"encoding/xml"
	"errors"
	"time"
)

type Currency struct {
	ID    int64     `json:"id"`
	Title string    `json:"title"`
	Code  string    `json:"code"`
	Value float64   `json:"value"`
	ADate time.Time `json:"a_date"`
}

type CurrencyPointers struct {
	ID    *int64     `json:"id,omitempty"`
	Title *string    `json:"title,omitempty"`
	Code  *string    `json:"code,omitempty"`
	Value *float64   `json:"value,omitempty"`
	ADate *time.Time `json:"a_date,omitempty"`
}

type CurrencyService interface {
	SaveCurrency(date string) (map[string]bool, error)
	CreateCurrency(currency *Currency) (*Currency, error)
	FindCurrencies(pointers CurrencyPointers) ([]*Currency, error)
	FindCurrency(id int64) (*Currency, error)
	UpdateCurrency(pointers CurrencyPointers) error
	DeleteCurrency(id int64) error
}

type CurrencyRepository interface {
	Create(currency *Currency) (*Currency, error)
	FindAll(pointers CurrencyPointers) ([]*Currency, error)
	FindByID(id int64) (*Currency, error)
	Update(pointers CurrencyPointers) error
	Delete(id int64) error
}

var (
	ErrNothingToUpdate  = errors.New("nothing to update.")
	ErrCurrencyNotFound = errors.New("currency not found.")
)

type RatesXML struct {
	XMLName xml.Name  `xml:"rates"`
	Item    []ItemXML `xml:"item"`
}

type ItemXML struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"fullname"`
	Code    string   `xml:"title"`
	Value   float64  `xml:"description"`
}