package lang

import (
	"github.com/sku4/alice-checklist/lang/json"
)

const (
	Ru = "ru"
	En = "en"
)

type Message interface {
	Translate(args ...string) string
}

type Localize struct {
	Message
}

func InitLocalize(defaultLang string) (*Localize, error) {
	newMessage, err := json.NewMessage(defaultLang)
	if err != nil {
		return nil, err
	}
	return &Localize{
		Message: newMessage,
	}, nil
}
