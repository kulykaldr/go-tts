package wellsaidlabs

import "github.com/chromedp/chromedp"

func (wl *Wellsaidlabs) Close() error {
	return chromedp.Cancel(wl.ctx)
}
