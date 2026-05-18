/* ---------------------------
pustanje jednog obicnog talasa
--------------------------- */

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

	// argumenti: frekvencija, glasnoća (0.0 do 1.0), trajanje
	playTone(otoCtx, 440.0, 0.1, 3*time.Second)          // Tih A4
	playTone(otoCtx, 880.0, 0.5, 500*time.Millisecond) // Glasniji A5
}

func playTone(ctx *oto.Context, ffreq float64, volume float64, duration time.Duration) {

	// TODO: napraviti oscilator

	osc := &sineWave{
		
	}

	p := ctx.NewPlayer(osc)
	p.SetVolume(volume)
	p.Play()
	time.Sleep(duration)
	p.Pause()

}

type sineWave struct {
	// TODO: napisati strukturu koja opisuje sinusni talas

}

func (s *sineWave) Read(p []byte) (n int, err error) {

	// TODO: napisati funkciju koja puni buffer sa podacima

	// HINT:
	// bits := math.Float32bits(float32(v))
	// binary.LittleEndian.PutUint32(p[i*4:], bits)

}
