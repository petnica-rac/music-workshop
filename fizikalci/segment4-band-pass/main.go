package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/ebitengine/oto/v3"
)

const sampleRate int = 44100

func main() {
	op := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 1,
		Format:       oto.FormatFloat32LE,
	}
	otoCtx, ready, _ := oto.NewContext(op)
	<-ready

	for i := 40; i < 110; i++ {
		freq := 440.0 * math.Exp2((float64(i)-69.0)/12.0)
		playTone(otoCtx, freq, 1, 300*time.Millisecond)
		fmt.Printf("Freq: %.2f\n", freq)
	}
}

func playTone(ctx *oto.Context, freq float64, volume float64, duration time.Duration) {
	osc := &sineWave{
		freq:   freq,
		sample: 0,
	}
	p := ctx.NewPlayer(osc)
	p.SetVolume(volume)
	p.Play()
	time.Sleep(duration)
	p.Pause()
}

type sineWave struct {
	freq   float64
	sample int

	prevPrevX float64
	prevX     float64
	prevY     float64
	prevPrevY float64

	filter string
}

func bandPass(x float64, prevX *float64, prevPrevX *float64, prevY *float64, prevPrevY *float64) float64 {
	y := 0.02206*x - 0.02206*(*prevX) + 1.9671*(*prevY) - 0.9725*(*prevPrevY)

	*prevPrevX = *prevX
	*prevX = x
	*prevPrevY = *prevY
	*prevY = y
	return y
}

func (s *sineWave) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p)/4; i++ {
		v := math.Sin(2 * math.Pi * s.freq * float64(s.sample) / float64(sampleRate))
		s.sample++

		v = bandPass(v, &s.prevX, &s.prevPrevX, &s.prevY, &s.prevPrevY)

		v *= math.Exp(-3 * float64(s.sample) / float64(sampleRate))

		bits := math.Float32bits(float32(v))
		binary.LittleEndian.PutUint32(p[i*4:], bits)
	}
	return len(p), nil
}
