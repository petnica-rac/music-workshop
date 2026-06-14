/* ---------------------------
-----     svi efekti    ------
--------------------------- */

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

}

func (o *oscillator) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p)/4; i++ {

		var v float64
		var tStamp float64
		tStamp = float64(o.sample) / float64(sampleRate)

		switch o.sound {
		case "sine":
			v = math.Sin(2 * math.Pi * o.freq * float64(o.sample) / float64(sampleRate))
		case "piano":
		case "square":
		case "distortion":
		default:
			return len(p), fmt.Errorf("Type of sound doesn't match %s\n", o.sound)
		}

		o.sample++

		bits := math.Float32bits(float32(v))
		binary.LittleEndian.PutUint32(p[i*4:], bits)
	}
	return len(p), nil
}
