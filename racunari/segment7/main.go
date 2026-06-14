/*
---------------------------
sintisajzer sa efektima!!!!
---------------------------
*/
package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/ebitengine/oto/v3"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // Inicijalizuje drajvere za Windows/Mac/Linux
	// go get gitlab.com/gomidi/midi/v2
	// go get gitlab.com/gomidi/midi/v2/drivers/rtmididrv
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

	BuildMidiFreqTable()

	// po defaultu zvuk je klavir
	// sound := "piano"

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
				fmt.Printf("Nota: %d\n", note)
			} else {
				// TODO: dodati podrsku za efekte
				//go playTone(otoCtx, midiFreqTable[note], 1, 1000*time.Millisecond, sound)
				go play(otoCtx)
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

func applyEffect(note uint8) (string, error) {
	return "", nil
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

		amp := math.Exp(-3 * tStamp)
		v *= amp

		o.sample++

		bits := math.Float32bits(float32(v))
		binary.LittleEndian.PutUint32(p[i*4:], bits)
	}
	return len(p), nil
}

func play(otoCtx *oto.Context) {

	const (
		Sixteenth  = 120 * time.Millisecond
		Dotted16th = 180 * time.Millisecond
		Eighth     = 240 * time.Millisecond
		Dotted8th  = 360 * time.Millisecond
		Quarter    = 480 * time.Millisecond
		Half       = 960 * time.Millisecond
	)

	freqs := []float64{
		// Phrase 1: "Never gonna give you up,"
		415.30, 466.16, 554.37, 466.16, 698.46, 698.46, 622.25, 0.0, // rest

		// Phrase 2: "Never gonna let you down,"
		415.30, 466.16, 554.37, 466.16, 622.25, 622.25, 554.37, 523.25, 466.16, 0.0, // rest

		// Phrase 3: "Never gonna run around and desert you,"
		415.30, 466.16, 554.37, 466.16, 554.37, 622.25, 523.25, 466.16, 415.30, 415.30, 622.25, 554.37, 0.0, // rest

		// Phrase 4: "Never gonna make you cry,"
		415.30, 466.16, 554.37, 466.16, 698.46, 698.46, 622.25, 0.0, // rest

		// Phrase 5: "Never gonna say goodbye,"
		415.30, 466.16, 554.37, 466.16, 830.61, 698.46, 554.37, 523.25, 466.16, 0.0, // rest

		// Phrase 6: "Never gonna tell a lie and hurt you."
		415.30, 466.16, 554.37, 466.16, 554.37, 622.25, 523.25, 554.37, 466.16, 415.30, 554.37,
	}

	durations := []time.Duration{
		// Phrase 1: "Never gonna give you up,"
		Eighth, Eighth, Eighth, Eighth, Dotted8th, Dotted8th, Quarter, Eighth,

		// Phrase 2: "Never gonna let you down,"
		Eighth, Eighth, Eighth, Eighth, Dotted8th, Dotted8th, Eighth, Eighth, Eighth, Eighth,

		// Phrase 3: "Never gonna run around and desert you,"
		Eighth, Eighth, Eighth, Eighth, Quarter, Eighth, Eighth, Eighth, Eighth, Eighth, Dotted8th, (Quarter + Eighth), Eighth,

		// Phrase 4: "Never gonna make you cry,"
		Eighth, Eighth, Eighth, Eighth, Dotted8th, Dotted8th, Quarter, Eighth,

		// Phrase 5: "Never gonna say goodbye,"
		Eighth, Eighth, Eighth, Eighth, Dotted8th, Dotted8th, Eighth, Eighth, Eighth, Eighth,

		// Phrase 6: "Never gonna tell a lie and hurt you."
		Eighth, Eighth, Eighth, Eighth, Quarter, Eighth, Eighth, Eighth, Eighth, Eighth, (Half + Eighth),
	}

	_, _ = freqs, durations

	for i, v := range freqs {
		playTone(otoCtx, v, 1, durations[i], "piano")
	}
}
