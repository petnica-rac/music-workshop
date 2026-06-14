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

	chord := []float64{400.0, 500.0, 600.0}
	playChord(otoCtx, chord, 0.6, 3*time.Second)
}

func playChord(ctx *oto.Context, freq []float64, volume float64, duration time.Duration) {
	osc := &multipleWave{
		freqs:  freq,
		sample: 0,
	}
	p := ctx.NewPlayer(osc)
	p.SetVolume(volume)
	p.Play()
	time.Sleep(duration)
	p.Pause()
}

type multipleWave struct {
	freqs  []float64
	sample int
}

func (m *multipleWave) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p)/4; i++ {

		v := 0.0

		for _, freq := range m.freqs {
			v += math.Sin(2 * math.Pi * freq * float64(m.sample) / float64(sampleRate))
		}

		v /= float64(len(m.freqs))

		m.sample++

		bits := math.Float32bits(float32(v))
		binary.LittleEndian.PutUint32(p[i*4:], bits)
	}
	return len(p), nil
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
}

func (s *sineWave) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p)/4; i++ {
		v := math.Sin(2 * math.Pi * s.freq * float64(s.sample) / float64(sampleRate))
		s.sample++

		bits := math.Float32bits(float32(v))
		binary.LittleEndian.PutUint32(p[i*4:], bits)
	}
	return len(p), nil
}
