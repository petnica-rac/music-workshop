package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/ebitengine/oto/v3"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // Inicijalizuje drajvere za Windows/Mac/Linux
)

const sampleRate int = 44100

var midiFreqTable [128]float64

func BuildMidiFreqTable() {
	for i := 0; i < 128; i++ {
		midiFreqTable[i] = 440.0 * math.Pow(2.0, float64(i-69)/12.0)
	}
}

func main() {
	op := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 1,
		Format:       oto.FormatFloat32LE,
	}
	otoCtx, ready, _ := oto.NewContext(op)
	<-ready

	defer midi.CloseDriver()

	// napravimo tabelu koja mapira midi notu sa frekkvencijom
	BuildMidiFreqTable()

	inPorts := midi.GetInPorts()

	if len(inPorts) == 0 {
		fmt.Println("Nije pronađena nijedna klavijatura!")
		return
	}

	in := inPorts[0]

	fmt.Println("Pritisni dirke na klavijaturi")

	_, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var channel, note, velocity uint8

		if msg.GetNoteOn(&channel, &note, &velocity) {
			if velocity == 0 {
				fmt.Printf("[Note OFF] Nota: %d\n", note)
			} else {
				go playTone(otoCtx, midiFreqTable[note], 1, 1000*time.Millisecond, "piano")
			}
		}
	})

	if err != nil {
		fmt.Printf("Greška pri slušanju: %v\n", err)
		return
	}

	for {
		time.Sleep(time.Second)
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
			if v >= 0 {
				v = 1.0
			} else {
				v = -1.0
			}
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
