package currency

import (
	"bytes"
	"encoding/xml"
	"github.com/kas2000/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	CurrencyService
}

type service struct {
	log          logger.Logger
	currencyRepo CurrencyRepository
}

func NewService(
	log logger.Logger,
	currencyRepo CurrencyRepository,
) Service {
	return &service{
		log:          log,
		currencyRepo: currencyRepo,
	}
}

func (svc *service) SaveCurrency(date string) (map[string]bool, error) {
	resp, err := http.Get("https://nationalbank.kz/rss/get_rates.cfm?fdate="+date)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	xmlReader := bytes.NewReader(body)
	data := new(RatesXML)
	if err := xml.NewDecoder(xmlReader).Decode(data); err != nil {
		return nil, err
	}
	go func() {
		//TODO: сделать нормальное преобразование
		d := strings.Split(date, ".")
		dt, err := time.Parse("2006-01-02 15:04:05", d[2] + "-" + d[1] + "-" + d[0] + " 00:00:00")
		if err != nil {

		}
		for _, item := range data.Item {
			_, err = svc.currencyRepo.Create(&Currency{
				Title: item.Title,
				Code:  item.Code,
				Value: item.Value,
				ADate: dt,
			})
			if err != nil {
				svc.log.Warn("Error saving currency. Couldn't create currency: " + err.Error(), zap.Any("currency_title", item.Title), zap.Any("currency_code", item.Code))
			}
		}
	}()

	return map[string]bool{
		"success": true,
	}, nil
}

func (svc *service) CreateCurrency(currency *Currency) (*Currency, error) {
	return svc.currencyRepo.Create(currency)
}

func (svc *service) FindCurrencies(pointers CurrencyPointers) ([]*Currency, error) {
	return svc.currencyRepo.FindAll(pointers)
}

func (svc *service) FindCurrency(id int64) (*Currency, error) {
	return svc.currencyRepo.FindByID(id)
}

func (svc *service) UpdateCurrency(pointers CurrencyPointers) error {
	return svc.currencyRepo.Update(pointers)
}

func (svc *service) DeleteCurrency(id int64) error {
	return svc.currencyRepo.Delete(id)
}