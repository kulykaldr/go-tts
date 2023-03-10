package wellsaidlabs

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"time"
)

func enterToWellsaidlabs(ctx context.Context, voice string) error {
	err := chromedp.Run(ctx, chromedp.Tasks{
		network.Enable(),
		chromedp.Sleep(5 * time.Second),
		chromedp.Click(`div[data-e2e="project-card"]`, chromedp.NodeVisible),
		chromedp.Sleep(5 * time.Second),
		chromedp.Click(`//*[@id="page-studio"]//img[@class="MuiAvatar-img"]/ancestor::button`, chromedp.NodeVisible),
		chromedp.Click(fmt.Sprintf(`//p[text()='%s']`, voice), chromedp.NodeVisible),
		chromedp.Sleep(2 * time.Second),
	})
	if err != nil {
		return err
	}

	return nil
}
