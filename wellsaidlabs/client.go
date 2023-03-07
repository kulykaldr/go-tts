package wellsaidlabs

import (
	"github.com/chromedp/cdproto/emulation"
)

type Wellsaidlabs struct {
	config      Config
	uaEmulation *emulation.SetUserAgentOverrideParams
}

func NewClient(cfg *Config) *Wellsaidlabs {
	const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"

	if cfg == nil {
		cfg = &Config{
			Login:     "",
			Password:  "",
			Voice:     "Alana B.",
			Headless:  true,
			Debug:     true,
			Timeout:   15,
			UserAgent: userAgent,
		}
	}

	var currUserAgent string
	if cfg.UserAgent == "" {
		currUserAgent = userAgent
	} else {
		currUserAgent = cfg.UserAgent
	}

	return &Wellsaidlabs{
		config:      *cfg,
		uaEmulation: emulation.SetUserAgentOverride(currUserAgent),
	}
}
