/* -------------------------
	     	MIDI		
------------------------- */

package main

import (
	"fmt"
	"time"
	// "math"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // Inicijalizuje drajvere za Windows/Mac/Linux
	// go get gitlab.com/gomidi/midi/v2
	// go get gitlab.com/gomidi/midi/v2/drivers/rtmididrv
)


/* ---------------------

TODO: pored MIDI broja note ispisati i frekvenciju
i tom logikom popuniti midiFreqTable

--------------------- */


var midiFreqTable [128]float64

func BuildMidiFreqTable() {
}


func main() {
	// 1. Zatvori MIDI drajvere na kraju programa
	defer midi.CloseDriver()
	
	// TODO: napraviti tabelu koja mapira midi notu sa frekkvencijom
	BuildMidiFreqTable()

	// 2. Izlistaj sve dostupne MIDI ulaze - sta je sve detektovano
	fmt.Println("Dostupni MIDI ulazi:")
	inPorts := midi.GetInPorts()
	for _, port := range inPorts {
		fmt.Printf("[%d] -> % s\n", port.Number(), port.String())
	}

	if len(inPorts) == 0 {
		fmt.Println("Nije pronađena nijedna klavijatura!")
		return
	}

	// 3. Odaberi prvi slobodan port (obično tvoja klavijatura)
	in := inPorts[0]
	fmt.Printf("\nSlušam na portu: %s\n", in.String())

	fmt.Println("Pritisni dirke na klavijaturi")

	// 4. Pokreni slušanje
	// TODO: pored MIDI broja note ispisati i frekvenciju
	_, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var channel, note, velocity uint8
		if (msg.GetNoteOn(&channel, &note, &velocity)){ 
			if velocity == 0 {
				fmt.Printf("[Note OFF] Nota: %d\n", note)
			} else {
				fmt.Printf("[Note ON] Nota: %d\n", note)
			}
		}
	})

	if err != nil {
		fmt.Printf("Greška pri slušanju: %v\n", err)
		return
	}

	// 5. Drži program budnim
	for {
		time.Sleep(time.Second)
	}
}
