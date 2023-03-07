package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kulykaldr/go-tts/wellsaidlabs"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	var cfg *wellsaidlabs.Config
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	log.Println("Successfully Opened config.json")
	defer configFile.Close()
	byteValue, _ := io.ReadAll(configFile)
	json.Unmarshal(byteValue, &cfg)

	if cfg.Login == "" || cfg.Password == "" {
		log.Fatal("Login or Password not provided")
	}

	wl := wellsaidlabs.NewClient(cfg)
	ctx, cancel := wl.CreateContext(context.Background())
	defer cancel()

	err = wl.SignIn(ctx)
	if err != nil {
		log.Fatalf("signin error: %v\n", err)
	}

	dirPath := "articles/"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("get files error: %v\n", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		textPath := path.Join(dirPath, file.Name())
		voice, err := wl.GetVoice(ctx, textPath, cfg.Voice)
		if err != nil {
			log.Fatalf("get voice error: %v\n", err)
		}

		_, file := path.Split(textPath)
		outFileName := strings.Replace(file, path.Ext(file), "", -1)
		outFilePath := path.Join("voices", outFileName+".mp3")

		wellsaidlabs.CreateDirPath("voices")
		if err = os.WriteFile(outFilePath, voice, 0644); err != nil {
			log.Fatalf("write voice error: %v\n", err)
		}

		err = os.Remove(textPath)
		if err != nil {
			return
		}
	}

	err = wl.Close(ctx)
	if err != nil {
		log.Fatalf("close context error: %v\n", err)
	}
}
