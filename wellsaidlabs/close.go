package wellsaidlabs

import (
	"context"
	"github.com/chromedp/chromedp"
)

func (wl *Wellsaidlabs) Close(ctx context.Context) error {
	return chromedp.Cancel(ctx)
}
