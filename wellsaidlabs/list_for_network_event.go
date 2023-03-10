package wellsaidlabs

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
)

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
