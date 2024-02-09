package currency

import "github.com/kas2000/http"

type currencyController struct {
	server      *http.Server
	httpFactory *CurrencyHttp
	prefix      string
}

func NewCurrencyController(server *http.Server, httpFactory *CurrencyHttp, prefix string) *currencyController {
	return &currencyController{
		server:      server,
		httpFactory: httpFactory,
		prefix:      prefix,
	}
}

func (tc *currencyController) Bind() {
	srvr := *tc.server
	srvr.Handle("GET", tc.prefix+"/currency/save/{date}", tc.httpFactory.SaveCurrency("date"))
	srvr.Handle("GET", tc.prefix+"/currency", tc.httpFactory.FindCurrencies())
	srvr.Handle("POST", tc.prefix+"/currency", tc.httpFactory.CreateCurrency())
	srvr.Handle("PUT", tc.prefix+"/currency/{id}", tc.httpFactory.UpdateCurrency("id"))
	srvr.Handle("DELETE", tc.prefix+"/currency/{id}", tc.httpFactory.DeleteCurrency("id"))
}