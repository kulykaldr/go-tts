package wellsaidlabs

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (wl *Wellsaidlabs) CreateContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctxInterr, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	// create a timeout
	timeCtx, cancel := context.WithTimeout(ctxInterr, time.Duration(wl.config.Timeout)*time.Minute)

	profileDir := CreateDirPath("wellsaid_profiles", wl.config.Login)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		//chromedp.DisableGPU,
		chromedp.Flag("headless", wl.config.Headless), // To display the browser
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("remote-debugging-port", "9222"),
		chromedp.WindowSize(1920, 1080), // init with a desktop view
		chromedp.UserDataDir(profileDir),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(timeCtx, opts...)

	var ctxOpts []chromedp.ContextOption
	if wl.config.Debug {
		ctxOpts = append(ctxOpts, chromedp.WithDebugf(log.Printf))
	}

	cdpCtx, cancel := chromedp.NewContext(
		allocCtx,
		ctxOpts...,
	)

	return cdpCtx, cancel
}
