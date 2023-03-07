package wellsaidlabs

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/nochso/gomd/eol"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var voiceSlice [][]byte

func (wl *Wellsaidlabs) GetVoice(textPath string, voice string) ([]byte, error) {
	text, err := os.ReadFile(textPath)
	if err != nil {
		return nil, err
	}

	err = chromedp.Run(wl.ctx, chromedp.Tasks{
		network.Enable(),
		chromedp.Sleep(5 * time.Second),
		chromedp.Click(`div[data-e2e="project-card"]`, chromedp.NodeVisible),
		chromedp.Sleep(5 * time.Second),
		chromedp.Click(`//*[@id="page-studio"]//img[@class="MuiAvatar-img"]/ancestor::button`, chromedp.NodeVisible),
		chromedp.Click(fmt.Sprintf(`//p[text()='%s']`, voice), chromedp.NodeVisible),
		chromedp.Sleep(2 * time.Second),
	})
	if err != nil {
		return nil, err
	}

	el, _ := eol.Detect(string(text))
	re, _ := regexp.Compile(`[\w\W\s]{0,900}[^.!? ]+[.!?]+`)
	tempStr := strings.Replace(string(text), el.String(), " ", -1)

	listenForNetworkEvent(wl.ctx)
	for {
		textPart := re.FindString(tempStr)
		tempStr = strings.Replace(tempStr, textPart, "", 1)

		textAreaSel := `textarea[data-e2e="project-editor"]`
		if err = chromedp.Run(wl.ctx, chromedp.Tasks{
			chromedp.SetValue(textAreaSel, ""),
			chromedp.SendKeys(textAreaSel, textPart, chromedp.NodeVisible),
			//chromedp.Sleep(2 * time.Second),
			chromedp.Click(`button[data-e2e="project-editor-submit"]`, chromedp.NodeVisible),
			chromedp.WaitReady(`div[data-rbd-draggable-context-id="0"]`),
			//chromedp.Sleep(2 * time.Second),
			chromedp.Click(`button[aria-label="delete clip"]`, chromedp.NodeVisible),
		}); err != nil {
			return nil, err
		}

		if len(strings.TrimSpace(tempStr)) == 0 {
			break
		}

		time.Sleep(2 * time.Second)
	}

	time.Sleep(5 * time.Second)
	bVoice := bytes.Join(voiceSlice, []byte(""))

	return bVoice, nil
}

// this will be used to capture the request id for matching network events
var requestID network.RequestID

// set up a listener to watch the network events and close the channel when
// complete the request id matching is important both to filter out
// unwanted network events and to reference the downloaded file later
func listenForNetworkEvent(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventResponseReceived:
			if strings.Contains(ev.Response.URL, "text_to_speech/stream") {
				requestID = ev.RequestID
			}
		case *network.EventLoadingFinished:
			if ev.RequestID == requestID {
				go func() {
					c := chromedp.FromContext(ctx)
					rbp := network.GetResponseBody(ev.RequestID)
					body, err := rbp.Do(cdp.WithExecutor(ctx, c.Target))
					if err != nil {
						log.Fatalf("body err: %s", err)
					}

					voiceSlice = append(voiceSlice, body)
				}()
			}
		}
	})
}
