package wellsaidlabs

type Config struct {
	Login     string `json:"login,omitempty"`
	Password  string `json:"password,omitempty"`
	Voice     string `json:"voice,omitempty"`
	Headless  bool   `json:"headless,omitempty"`
	Debug     bool   `json:"debug,omitempty"`
	Timeout   int64  `json:"timeout,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}
