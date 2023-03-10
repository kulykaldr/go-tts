package wellsaidlabs

import (
	"bytes"
	"context"
	"github.com/nochso/gomd/eol"
	"os"
	"strings"
	"time"
)

func (wl *Wellsaidlabs) GetVoiceByParagraph(ctx context.Context, textPath string, voice string) ([]byte, error) {
	text, err := os.ReadFile(textPath)
	if err != nil {
		return nil, err
	}

	err = enterToWellsaidlabs(ctx, voice)
	if err != nil {
		return nil, err
	}

	el, _ := eol.Detect(string(text))
	list := strings.Split(string(text), el.String())

	listenForNetworkEvent(ctx)

	for _, str := range list {
		if len(strings.TrimSpace(str)) == 0 {
			continue
		}

		err = inputTextToField(ctx, str)
		if err != nil {
			return nil, err
		}

		time.Sleep(2 * time.Second)
	}

	time.Sleep(5 * time.Second)
	bVoice := bytes.Join(voiceSlice, []byte(""))

	return bVoice, nil
}
