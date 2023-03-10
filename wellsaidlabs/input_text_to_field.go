package wellsaidlabs

import (
	"context"
	"github.com/chromedp/chromedp"
)

func inputTextToField(ctx context.Context, text string) error {
	textAreaSel := `textarea[data-e2e="project-editor"]`
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.SetValue(textAreaSel, ""),
		chromedp.SendKeys(textAreaSel, text, chromedp.NodeVisible),
		//chromedp.Sleep(2 * time.Second),
		chromedp.Click(`button[data-e2e="project-editor-submit"]`, chromedp.NodeVisible),
		chromedp.WaitReady(`div[data-rbd-draggable-context-id="0"]`),
		//chromedp.Sleep(2 * time.Second),
		chromedp.Click(`button[aria-label="delete clip"]`, chromedp.NodeVisible),
	}); err != nil {
		return err
	}

	return nil
}
