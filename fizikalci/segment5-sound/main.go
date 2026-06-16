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

	freqs := []float64{261.63, 293.67, 329.63, 349.23, 392.00, 440.00, 493.88, 523.25}

	for _, v := range freqs {
		playTone(otoCtx, v*0.66, 1, 500*time.Millisecond, "sine")
	}

	for _, v := range freqs {
		playTone(otoCtx, v*0.66, 1, 500*time.Millisecond, "piano")
	}

	for _, v := range freqs {
		playTone(otoCtx, v*0.66, 1, 500*time.Millisecond, "square")
	}
	for _, v := range freqs {
		playTone(otoCtx, v*0.66, 1, 500*time.Millisecond, "distortion")
	}

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

func playTone(ctx *oto.Context, freq float64, volume float64, duration time.Duration, sound string) {

	osc := &oscillator{
		freq:   freq,
		sample: 0,
		sound:  sound,
	}
	p := ctx.NewPlayer(osc)
	p.SetVolume(volume)
	p.Play()
	time.Sleep(duration)
	p.Pause()
}

type oscillator struct {
	freq   float64
	sample int
	sound  string
}

func getSine(freq float64, tStamp float64) float64 {
	return math.Sin(2 * math.Pi * freq * tStamp)
}

func getPiano(freq float64, tStamp float64) float64 {
	v := math.Sin(2*math.Pi*freq*tStamp) +
		0.5*math.Sin(2*2*math.Pi*freq*tStamp) +
		0.25*math.Sin(3*2*math.Pi*freq*tStamp) +
		0.125*math.Sin(4*2*math.Pi*freq*tStamp) +
		0.06*math.Sin(5*2*math.Pi*freq*tStamp)
	v /= 1.93
	amp := math.Exp(-3 * tStamp)
	v *= amp
	return v
}

func (o *oscillator) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p)/4; i++ {

		var v float64
		tStamp := float64(o.sample) / float64(sampleRate)

		switch o.sound {
		case "sine":
			v = getSine(o.freq, tStamp)
		case "piano":
			v = getPiano(o.freq, tStamp)
		case "square":
			v = getSine(o.freq, tStamp)
			v += 1.0 / 3.0 * getSine(3 * o.freq, tStamp)
			v += 1.0 / 5.0 * getSine(5 * o.freq, tStamp)
			v += 1.0 / 7.0 * getSine(7 * o.freq, tStamp)
			v += 1.0 / 9.0 * getSine(9 * o.freq, tStamp)
			v /= 1.7873
		case "distortion":
			v = getPiano(o.freq, tStamp)
			const gain = 3.0
			const limit = 0.7
			v *= gain
			if v > limit {
				v = limit
			} else if v < -limit {
				v = -limit
			}
		default:
			return len(p), fmt.Errorf("Type of sound doesn't match %s\n", o.sound)
		}

		o.sample++

		bits := math.Float32bits(float32(v))
		binary.LittleEndian.PutUint32(p[i*4:], bits)
	}
	return len(p), nil
}
