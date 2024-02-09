package currency

import (
	"encoding/json"
	"github.com/gorilla/mux"
	command "github.com/kas2000/commandlib"
	"github.com/kas2000/logger"
	"io/ioutil"
	"net/http"
	httpLib "github.com/kas2000/http"
	"strconv"
	"strings"
	"time"
)

type CurrencyHttp struct {
	log        logger.Logger
	ch         command.CommandHandler
	systemName string
}

func NewCurrencyHttpHandler(log logger.Logger, ch command.CommandHandler, systemName string) *CurrencyHttp {
	return &CurrencyHttp{
		log:        log,
		ch:         ch,
		systemName: systemName,
	}
}

func (factory *CurrencyHttp) SaveCurrency(dateString string) httpLib.Endpoint {
	return func(w http.ResponseWriter, r *http.Request) httpLib.Response {
		vars := mux.Vars(r)
		date, found := vars[dateString]
		if !found {
			return httpLib.BadRequest(310, "No subject id", factory.systemName)
		}

		cmd := SaveCurrencyCommand{Date: date}

		resp, err := factory.ch.ExecuteCommand(&cmd)
		if err != nil {
			switch err {
			case ErrCurrencyNotFound:
				return httpLib.NotFound(330, err.Error(), factory.systemName)
			}
			return httpLib.InternalServer(340, err.Error(), factory.systemName)
		}
		return httpLib.NewResponse(http.StatusOK, resp, nil)
	}
}

func (factory *CurrencyHttp) FindCurrencies() httpLib.Endpoint {
	return func(w http.ResponseWriter, r *http.Request) httpLib.Response {
		//В тз указано, что дата и код - параметры
		var pointers CurrencyPointers
		if r.URL.Query().Has("date") {
			q := r.URL.Query().Get("date")
			//TODO: сделать нормальное преобразование
			d := strings.Split(q, ".")
			dt, err := time.Parse("2006-01-02 15:04:05", d[2] + "-" + d[1] + "-" + d[0] + " 00:00:00")
			if err != nil {

			}
			pointers.ADate = &dt
		}
		if r.URL.Query().Has("code") {
			code := r.URL.Query().Get("code")
			pointers.Code = &code
		}

		cmd := FindCurrenciesCommand{
			CurrencyPointers: pointers,
		}

		resp, err := factory.ch.ExecuteCommand(&cmd)
		if err != nil {
			return httpLib.InternalServer(260, err.Error(), factory.systemName)
		}
		return httpLib.NewResponse(http.StatusOK, resp, nil)
	}
}

func (factory *CurrencyHttp) FindCurrency(idParam string) httpLib.Endpoint {
	return func(w http.ResponseWriter, r *http.Request) httpLib.Response {
		var id string
		vars := mux.Vars(r)
		var ok bool
		id, ok = vars[idParam]
		if !ok {
			return httpLib.BadRequest(390, "No subject id", factory.systemName)
		}
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return httpLib.BadRequest(160, err.Error(), factory.systemName)
		}
		cmd := FindCurrencyCommand{
			ID: idInt,
		}
		resp, err := factory.ch.ExecuteCommand(&cmd)
		if err != nil {
			return httpLib.InternalServer(260, err.Error(), factory.systemName)
		}
		return httpLib.NewResponse(http.StatusOK, resp, nil)
	}
}

func (factory *CurrencyHttp) UpdateCurrency(idParameter string) httpLib.Endpoint {
	return func(w http.ResponseWriter, r *http.Request) httpLib.Response {
		var id string
		vars := mux.Vars(r)
		var ok bool
		id, ok = vars[idParameter]
		if !ok {
			return httpLib.BadRequest(390, "No subject id", factory.systemName)
		}
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return httpLib.BadRequest(200, err.Error(), factory.systemName)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return httpLib.BadRequest(410, "Error reading request body: "+err.Error(), factory.systemName)
		}
		var upd CurrencyPointers
		err = json.Unmarshal(body, &upd)
		if err != nil {
			return httpLib.BadRequest(420, "Error unmarshalling: "+err.Error(), factory.systemName)
		}
		upd.ID = &idInt

		cmd := UpdateCurrencyCommand{
			CurrencyPointers: upd,
		}

		resp, err := factory.ch.ExecuteCommand(&cmd)
		if err != nil {
			switch err {
			case ErrCurrencyNotFound:
				return httpLib.NotFound(440, err.Error(), factory.systemName)
			}
			return httpLib.InternalServer(450, err.Error(), factory.systemName)
		}
		return httpLib.NewResponse(http.StatusOK, resp, nil)
	}
}

func (factory *CurrencyHttp) CreateCurrency() httpLib.Endpoint {
	return func(w http.ResponseWriter, r *http.Request) httpLib.Response {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return httpLib.BadRequest(140, "Error reading request body: "+err.Error(), factory.systemName)
		}
		var currency Currency
		err = json.Unmarshal(body, &currency)
		if err != nil {
			return httpLib.BadRequest(150, "Error unmarshalling: "+err.Error(), factory.systemName)
		}
		cmd := CreateCurrencyCommand{Currency: &currency}
		resp, err := factory.ch.ExecuteCommand(&cmd)
		if err != nil {
			return httpLib.InternalServer(180, err.Error(), factory.systemName)
		}
		return httpLib.NewResponse(http.StatusOK, resp, nil)
	}
}

func (factory *CurrencyHttp) DeleteCurrency(idParameter string) httpLib.Endpoint {
	return func(w http.ResponseWriter, r *http.Request) httpLib.Response {
		var id string
		vars := mux.Vars(r)
		var ok bool
		id, ok = vars[idParameter]
		if !ok {
			return httpLib.BadRequest(390, "No subject id", factory.systemName)
		}
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return httpLib.BadRequest(200, err.Error(), factory.systemName)
		}

		cmd := DeleteCurrencyCommand{ID: idInt}

		resp, err := factory.ch.ExecuteCommand(&cmd)
		if err != nil {
			switch err {
			case ErrCurrencyNotFound:
				return httpLib.NotFound(170, err.Error(), factory.systemName)
			}
			return httpLib.InternalServer(180, err.Error(), factory.systemName)
		}
		return httpLib.NewResponse(http.StatusOK, resp, nil)
	}
}