package json

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	LangFolder = "./lang/json/"
	MaskJson   = "*.json"
)

type Message struct {
	defaultLang string
	lang        map[string]map[string]string
}

func NewMessage(defaultLang string) (*Message, error) {
	files, err := filepath.Glob(LangFolder + MaskJson)
	if err != nil {
		return nil, err
	}
	langs := make(map[string]map[string]string)
	for _, file := range files {
		filename := filepath.Base(file)
		dataJson, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		lang := strings.Replace(filename, filepath.Ext(filename), "", -1)
		var mapJson map[string]string
		if err = json.Unmarshal(dataJson, &mapJson); err != nil {
			return nil, err
		}
		langs[lang] = mapJson
	}

	return &Message{
		defaultLang: defaultLang,
		lang:        langs,
	}, nil
}

// Translate
/*
 * First argument is key
 * Second argument is language
 */
func (m *Message) Translate(args ...string) (mess string) {
	key := args[0]
	language := ""
	if len(args) > 1 {
		language = args[1]
	}
	if language == "" {
		language = m.defaultLang
	}
	mess, exist := m.lang[language][key]
	if !exist {
		mess, exist = m.lang[m.defaultLang][key]
		if !exist {
			mess = key
		}
	}
	return
}
