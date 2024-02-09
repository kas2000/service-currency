package currency

type SaveCurrencyCommand struct {
	Date string
}

func (cmd *SaveCurrencyCommand) Execute(svc interface{}) (interface{}, error) {
	return svc.(Service).SaveCurrency(cmd.Date)
}

type FindCurrenciesCommand struct {
	CurrencyPointers
}

func (cmd *FindCurrenciesCommand) Execute(svc interface{}) (interface{}, error) {
	return svc.(Service).FindCurrencies(cmd.CurrencyPointers)
}

type CreateCurrencyCommand struct {
	*Currency
}

func (cmd *CreateCurrencyCommand) Execute(svc interface{}) (interface{}, error) {
	return svc.(Service).CreateCurrency(cmd.Currency)
}

type FindCurrencyCommand struct {
	ID int64
}

func (cmd *FindCurrencyCommand) Execute(svc interface{}) (interface{}, error) {
	return svc.(Service).FindCurrency(cmd.ID)
}

type DeleteCurrencyCommand struct {
	ID int64
}

func (cmd *DeleteCurrencyCommand) Execute(svc interface{}) (interface{}, error) {
	err := svc.(Service).DeleteCurrency(cmd.ID)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

type UpdateCurrencyCommand struct {
	CurrencyPointers
}

func (cmd *UpdateCurrencyCommand) Execute(svc interface{}) (interface{}, error) {
	err := svc.(Service).UpdateCurrency(cmd.CurrencyPointers)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}