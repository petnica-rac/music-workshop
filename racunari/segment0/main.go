package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
	// go get github.com/ebitengine/oto/v3
	// go get github.com/hajimehoshi/go-mp3
	// go mod init package-name, ako nisu jos
)

func main() {
	const (
		sampleRate   = 44100
		channelCount = 2
	)

	file, err := os.Open("./segment0/song.mp3")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	decoded, err := mp3.NewDecoder(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	op := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: channelCount,
		Format:       oto.FormatSignedInt16LE,
	}

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		fmt.Println(err)
	}
	<-readyChan

	player := otoCtx.NewPlayer(decoded)
	player.Play()

	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
}
