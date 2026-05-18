# Audio FX Laboratorija
## Kako "iskriviti" matematiku i napraviti muziku koja zaboli u stomak

Naučili smo da zvuk = talas = matematička funkcija:

```
y(t) = A · sin(2π · f · t)
```

`A` kontroliše glasnoću, `f` kontroliše visinu tona, `t` je vreme. Svaki efekat koji postoji dotiče jednu ili više ovih varijabli. To je sve.

---

## Distorzija i kvadratni talas

Najbrži način da uništiš čistotu sinusnog talasa: ograniči amplitudu silom.

```
signal = sin(2π · f · t) × 5     ← namerno preglasnjen

if signal > 0.8:  signal = 0.8
if signal < -0.8: signal = -0.8
```

Talas koji je bio lepa kriva sada ima odrezane vrhove. Ekstremni slučaj ovoga je kvadratni talas — limit je praktično nula, signal skače između `+1` i `-1` bez prelaza:

```
  +1  ──────┐      ┌──────┐      ┌──────
            │      │      │      │
  -1        └──────┘      └──────┘
```

Zašto to zvuči "prljavo"? Fourier-ova teorema kaže da kvadratni talas matematički sadrži beskonačno harmonika (3f, 5f, 7f...). Uho čuje sve to odjednom i percipira bogatstvo koje zovemo distorzijom.

Smanjivanjem limita na `0.1` dobijaš agresivniji zvuk. Na `0.01` skoro potpun kvadratni talas.

---

## Vibrato

Frekvencija tona nije statična — polako se ljulja gore-dole pod uticajem sporog oscilatora.

```
f(t) = 440 + sin(2π · 6 · t) × 10
```

Ovde `6 Hz` je brzina ljuljanja, `10` je dubina u hercima. Talas se skuplja i širi ritmički:

```
Normalan:  ∿∿∿∿∿∿∿∿∿∿∿∿∿∿∿∿∿∿∿∿   (ravnomerni razmak)

Vibrato:   ∿∿∿∿∿ ∿∿∿∿ ∿∿∿∿∿ ∿∿∿∿   (skuplja se i širi)
```

Frekvencija sporog oscilatora određuje karakter:

| LFO | Efekat |
|-----|--------|
| 1–3 Hz | Opersko, sporo vibrato |
| 5–8 Hz | Klasično — gitara, violina |
| 20–50 Hz | Tremolo, nestabilno |
| 80+ Hz | Više nije vibrato — to je FM sinteza, potpuno novi timbar |

Taj poslednji slučaj je posebno zanimljiv. Isti kod, samo brži broj, i dobijaš zvuk koji nema nikakve veze sa originalnim tonom. FM sinteza je osnova DX7 synthesizera iz 80ih — zvuk koji čuješ u skoro svakoj pop pesmi te decenije.

---

## Ring Modulation

Formula je banalna. Zvuk koji izlazi nije.

```
output(t) = sin(2π · f_carrier · t) × sin(2π · f_modulator · t)
```

Množenje dva sinusa matematički razbija oba originalna tona i na izlazu dobijaš samo zbir i razliku:

```
carrier   = 300 Hz
modulator =  50 Hz

Čuješ:     350 Hz  i  250 Hz
Ne čuješ:  300 Hz  i   50 Hz  ← nestaju
```

Ovo je razlog zašto zvuči robotski i metalno — nijedan od originalnih tonova više ne postoji u signalu. Čuješ nešto što nisi poslao unutra.

Daleci iz Doctor Who su ring modulator na ljudskom glasu. Daft Punk su ga koristili za vokale na "Da Funk". Rezultat uvek ima tu hladnoću koja ne pripada nijednom prirodnom instrumentu.

---

## Akord umesto tona

Akord je samo više sinusa u isto vreme. Odnosi frekvencija su precizni razlomci:

```
Osnova (C4):   f × 1.000  =  261.63 Hz
Terca (E4):    f × 1.250  =  327.03 Hz
Kvinta (G4):   f × 1.500  =  392.44 Hz
Oktava (C5):   f × 2.000  =  523.26 Hz
```

Razlika između dur i mol je samo jedan multiplikator:

| Akord | Intervali | Karakter |
|-------|-----------|---------|
| Dur | 1.0 · 1.25 · 1.5 | Svetao |
| Mol | 1.0 · 1.20 · 1.5 | Taman |
| Power chord | 1.0 · 1.5 | Sirov, bez karaktera |
| Diminished | 1.0 · 1.20 · 1.414 | Napetost |
| Sus4 | 1.0 · 1.333 · 1.5 | Napet, nedovršen, čeka razrešenje |

Sus4 je posebno interesantan za slušanje — uho ga doživljava kao pitanje bez odgovora. Svi The Edge gitarski uvodi iz U2 su u suštini sus4 koji se razrešava.

---

## ADSR — životni ciklus zvuka

Čim pustiš dirku, zvuk nestane trenutno. To zvuči neprirodno jer u stvarnom svetu nijedan zvuk ne radi tako. Klavirska žica odzvanja, violinska nota se utišava, udarac bubnja se raspada postepeno.

```
Glasnoća
  │
  │    ╭─╮
  │   ╱   ╲_________
  │  ╱               ╲
  │ ╱                  ╲────
  │╱                        ╲
  └──────────────────────────────  Vreme
      A    D      S         R
```

**Attack** — koliko brzo se glasnoća diže od nule. Klavir je brz, gudački instrument spor.  
**Decay** — pad do normalnog nivoa posle inicijalnog udarca.  
**Sustain** — nivo dok je dirka pritisnuta.  
**Release** — koliko polako zvuk nestaje posle otpuštanja.

```go
go func() {
    for vol := 0.8; vol > 0; vol -= 0.02 {
        player.SetVolume(vol)
        time.Sleep(15 * time.Millisecond)
    }
    player.Pause()
}()
```

Isti ton sa attack od 5ms vs 500ms zvuči kao potpuno drugačiji instrument.

---

## Delay i Echo

Uzmeš signal, sačekaš, pustiš ga ponovo — ali tiše. Ponavljaš.

```
output(t) = signal(t) + 0.6 × signal(t - Δt) + 0.36 × signal(t - 2Δt) + ...
```

Kašnjenje `Δt` određuje karakter:

```
Original:  ████░░░░░░░░░░░░░░░░░░░░░
Echo 1:    ░░░░░░░███░░░░░░░░░░░░░░░
Echo 2:    ░░░░░░░░░░░░░██░░░░░░░░░░
Echo 3:    ░░░░░░░░░░░░░░░░░░░█░░░░░
```

Kratko kašnjenje (ispod 30ms) i ne čuješ ga kao eho — čuješ ga kao promenu prostora, kao da si u manjoj sobi. Duže od 100ms i postaje prepoznatljivi eho. Sinhronizovano sa tempom pesme i postaje ritmički element sam po sebi — tako Pink Floyd pravi gitarsku liniju koja zvuči kao da se udvostručuje.

---

## Bitcrusher

Namerno smanjivanje kvaliteta zvuka. CD je 16-bit (65536 nivoa amplitude). Spusti na 4-bit (16 nivoa):

```
Originalni sample:  0.73847281

4-bit quantization:
  0.73847281 × 16 = 11.8...
  round(11.8) = 12
  12 / 16 = 0.75       ← malo pogrešno, ali to je poenta
```

Talas koji je bio glatka kriva postaje stepenice:

```
Originalni:  ╭──────────────╮   (glatko)

4-bit:       ╭──┐           ┌╮  (stepenice)
                 └──┐    ┌──┘
                     └───┘
```

Ta "greška" je lo-fi estetika. Chiptune, vaporwave, GameBoy melodije — sve je to namerno degradiran zvuk koji je negde između nostalgije i glitch arta.

---

## Reverb

Reverb je akustika prostorije simulirana u kodu. Zvuk odskače od zidova i stiže do tebe sa malim kašnjenjem iz svakog pravca. Prosta implementacija: hiljade kratkih delay-eva sa nasumičnim kašnjenjima i amplitudama.

```
reverb(t) = Σ  signal(t - τᵢ) × aᵢ
```

gde su `τᵢ` nasumična kašnjenja i `aᵢ` opadajuće amplitude.

Razlika između tipova reverba je samo u parametrima te raspodele:

| Tip | Decay time | Karakter |
|-----|-----------|---------|
| Room | 0.3–0.8s | Mala prostorija, intimno |
| Hall | 1–3s | Koncertna dvorana |
| Cathedral | 3–8s | Ogromno, atmosferično |
| Spring | 0.5–2s | Opružni reverb, gitarsko pojačalo, vintage |
| Plate | 0.5–2s | Studio klasik, glatki reverb |

Bez reverba zvuk je "suv" — čuješ izvor tona bez konteksta. Sa reverberacijom zvuk se smešta u prostor. Mozak automatski čita te refleksije i konstruiše osećaj veličine prostorije.

---

## Signal chain

Redosled efekata menja sve. Isti zvuk, drugačiji redosled:

```
Instrument → EQ → Distorzija → Chorus → Delay → Reverb → Izlaz
```

Distorzija pre reverba: distorzuješ čist signal, reverb ozvučava rezultat.  
Reverb pre distorzije: distorzuješ reverb zajedno sa signalom — muljavo, haotično, ponekad odlično.

Nema pogrešnog redosleda. Ima samo različitih zvukova.

---

## Šta probati

Promeni jedan parametar u jednom efektu i zapamti šta se desilo:

- Vibratu ubrzaj LFO sa `6` na `80` Hz i posluša kako vibrato postaje FM sinteza
- Distorziji spusti limit na `0.05` i vidi gde nestaje tonski identitet
- Ring modulatoru postavi modulator na istu frekvenciju kao carrier — šta izlazi?
- ADSR-u postavi attack na `2000ms` i isti ton postaje pad umesto nota
- Delay sinhronizuj sa tempom: ako je 120 BPM, jedan takt = 500ms

---

*Srećno sa sintisajzovanjem:)*

*P.S. ako vam se ovo svidja bacite pogled na programski jezik ChucK*