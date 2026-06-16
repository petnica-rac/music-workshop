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

	freqs := []float64{250, 290, 340, 370}
	playChord(otoCtx, freqs, 1, 3*time.Second)
	// probati promeniti gain sa 1 na neki veci broj -> cuje se distorzija

	/*
		Akord			Intervali				Karakter
		Dur				1.0 · 1.25 · 1.5		Svetao
		Mol				1.0 · 1.20 · 1.5		Taman
		Power chord		1.0 · 1.5				Sirov
		Diminished		1.0 · 1.20 · 1.414		Napetost
		Sus4			1.0 · 1.333 · 1.5		Nedovršen, čeka razrešenje
	*/

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
