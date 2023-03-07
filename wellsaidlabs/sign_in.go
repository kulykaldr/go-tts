package wellsaidlabs

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"time"
)

func (wl *Wellsaidlabs) signIn(ctx context.Context, login string, password string) error {
	if login == "" || password == "" {
		return fmt.Errorf("login or Password not provided")
	}

	uploadUrl := `https://wellsaidlabs.com/auth/sign_in?redirect_to=%2Fdashboard`

	loginInputSel := `input[name="email"]`

	var isLoginExists bool
	err := chromedp.Run(ctx, chromedp.Tasks{
		wl.uaEmulation,
		chromedp.Navigate(uploadUrl),
		chromedp.Sleep(10 * time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			var nodes []*cdp.Node
			if err := chromedp.Nodes(loginInputSel, &nodes, chromedp.AtLeast(0)).Do(ctx); err != nil {
				return err
			}
			if len(nodes) == 0 {
				isLoginExists = false
			} else {
				isLoginExists = true
			}
			return nil
		}),
	})
	if err != nil {
		return err
	}

	if !isLoginExists {
		return nil
	}

	signin := chromedp.Tasks{
		chromedp.SendKeys(loginInputSel, login, chromedp.NodeVisible),
		chromedp.Sleep(5 * time.Second),

		chromedp.SendKeys(`input[type="password"]`, password, chromedp.NodeVisible),
		chromedp.Click(`button[type="submit"]`, chromedp.NodeVisible),
		chromedp.Sleep(5 * time.Second),

		chromedp.WaitReady(`button[data-e2e="action-project-create"]`),
	}

	err = chromedp.Run(ctx, signin)
	if err != nil {
		return err
	}

	return nil
}
