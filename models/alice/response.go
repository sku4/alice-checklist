package alice

const (
	Version = "1.0"
)

type Response struct {
	Version  string `json:"version"`
	RespData `json:"response"`
}

func NewResponse() *Response {
	return &Response{
		Version: Version,
	}
}

type RespData struct {
	Text       string `json:"text"`
	Tts        string `json:"tts"`
	EndSession bool   `json:"end_session"`
}
