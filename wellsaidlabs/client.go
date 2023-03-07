package wellsaidlabs

import (
	"context"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/kulykaldr/go-tts/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Wellsaidlabs struct {
	ctx         context.Context
	config      Config
	uaEmulation *emulation.SetUserAgentOverrideParams
}

func NewClient(cfg *Config) (*Wellsaidlabs, error) {
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

	profileDir := utils.CreateDirPath("wellsaid_profiles", cfg.Login)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		//chromedp.DisableGPU,
		chromedp.Flag("headless", cfg.Headless), // To display the browser
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("remote-debugging-port", "9222"),
		chromedp.WindowSize(1920, 1080), // init with a desktop view
		chromedp.UserDataDir(profileDir),
	)

	ctxInterr, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// create a timeout
	timeCtx, cancel := context.WithTimeout(ctxInterr, time.Duration(cfg.Timeout)*time.Minute)
	defer cancel()

	allocCtx, cancel := chromedp.NewExecAllocator(timeCtx, opts...)
	defer cancel()

	var ctxOpts []chromedp.ContextOption
	if cfg.Debug {
		ctxOpts = append(ctxOpts, chromedp.WithDebugf(log.Printf))
	}

	cdpCtx, cancel := chromedp.NewContext(
		allocCtx,
		ctxOpts...,
	)
	defer cancel()

	var currUserAgent string
	if cfg.UserAgent == "" {
		currUserAgent = userAgent
	} else {
		currUserAgent = cfg.UserAgent
	}

	wl := &Wellsaidlabs{
		ctx:         cdpCtx,
		config:      *cfg,
		uaEmulation: emulation.SetUserAgentOverride(currUserAgent),
	}

	err := wl.signIn(cdpCtx, cfg.Login, cfg.Password)
	if err != nil {
		return nil, err
	}
	log.Print("Signin is successful")

	return wl, nil
}
