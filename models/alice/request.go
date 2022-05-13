package alice

type Request struct {
	Meta    `json:"meta"`
	Session `json:"session"`
	ReqData `json:"request"`
	State   `json:"state"`
	Version string `json:"version"`
}

type Meta struct {
	Locale     string `json:"locale"`
	Timezone   string `json:"timezone"`
	ClientId   string `json:"client_id"`
	Interfaces struct {
		Screen struct {
		} `json:"screen"`
		Payments struct {
		} `json:"payments"`
		AccountLinking struct {
		} `json:"account_linking"`
	} `json:"interfaces"`
}

type Session struct {
	MessageId int    `json:"message_id"`
	SessionId string `json:"session_id"`
	SkillId   string `json:"skill_id"`
	User      struct {
		UserId string `json:"user_id"`
	} `json:"user"`
	Application struct {
		ApplicationId string `json:"application_id"`
	} `json:"application"`
	New    bool   `json:"new"`
	UserId string `json:"user_id"`
}

type ReqData struct {
	Command           string `json:"command"`
	OriginalUtterance string `json:"original_utterance"`
	Nlu               struct {
		Tokens   []interface{} `json:"tokens"`
		Entities []struct {
			Type   string `json:"type"`
			Tokens struct {
				Start int `json:"start"`
				End   int `json:"end"`
			} `json:"tokens"`
			Value interface{} `json:"value"`
		} `json:"entities"`
		Intents map[string]Intent `json:"intents"`
	} `json:"nlu"`
	Markup struct {
		DangerousContext bool `json:"dangerous_context"`
	} `json:"markup"`
	Type string `json:"type"`
}

type Intent struct {
	Slots map[string]Slot `json:"slots"`
}

type Slot struct {
	Type   string `json:"type"`
	Tokens struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"tokens"`
	Value string `json:"value"`
}

type State struct {
	Session struct {
	} `json:"session"`
	User struct {
	} `json:"user"`
	Application struct {
	} `json:"application"`
}
