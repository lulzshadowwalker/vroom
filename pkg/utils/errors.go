package utils

type AppErr struct {
	message string
}

func (a *AppErr) Error() string {
	return a.message
}

func NewAppErr(message string) *AppErr {
	return &AppErr{
		message: message,
	}
}
