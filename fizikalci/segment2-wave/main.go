package main

import (
	"encoding/binary"
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

	playTone(otoCtx, 440.0, 0.6, 3*time.Second)
}

func playTone(ctx *oto.Context, freq float64, volume float64, duration time.Duration) {

	osc := &sqrWave{
		freq:   freq,
		sample: 0,
	}
	p := ctx.NewPlayer(osc)
	p.SetVolume(volume)
	p.Play()
	time.Sleep(duration)
	p.Pause()
}

type sqrWave struct {
	freq   float64
	sample int
}

func (s *sqrWave) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p)/4; i++ {

		v := math.Sin(2 * math.Pi * s.freq * float64(s.sample) / float64(sampleRate))
		v += 1.0 / 3.0 * math.Sin(3*2*math.Pi*s.freq*float64(s.sample)/float64(sampleRate))
		v += 1.0 / 5.0 * math.Sin(5*2*math.Pi*s.freq*float64(s.sample)/float64(sampleRate))
		v += 1.0 / 7.0 * math.Sin(7*2*math.Pi*s.freq*float64(s.sample)/float64(sampleRate))
		v += 1.0 / 9.0 * math.Sin(9*2*math.Pi*s.freq*float64(s.sample)/float64(sampleRate))

		v /= 1.7873

		s.sample++

		bits := math.Float32bits(float32(v))
		binary.LittleEndian.PutUint32(p[i*4:], bits)
	}
	return len(p), nil
}
