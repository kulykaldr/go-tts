package wellsaidlabs

import (
	"bytes"
	"context"
	"github.com/nochso/gomd/eol"
	"os"
	"regexp"
	"strings"
	"time"
)

var voiceSlice [][]byte

func (wl *Wellsaidlabs) GetVoice(ctx context.Context, textPath string, voice string) ([]byte, error) {
	text, err := os.ReadFile(textPath)
	if err != nil {
		return nil, err
	}

	err = enterToWellsaidlabs(ctx, voice)
	if err != nil {
		return nil, err
	}

	el, _ := eol.Detect(string(text))
	re, _ := regexp.Compile(`[\w\W\s]{0,900}[^.!? ]+[.!?]+`)
	tempStr := strings.Replace(string(text), el.String(), " ", -1)

	listenForNetworkEvent(ctx)
	for {
		textPart := re.FindString(tempStr)
		tempStr = strings.Replace(tempStr, textPart, "", 1)

		if len(strings.TrimSpace(textPart)) == 0 {
			break
		}

		err = inputTextToField(ctx, textPart)
		if err != nil {
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
