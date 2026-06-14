# Fizika – Audio

> *Momentalna struja*

Osnovni audio signal je samo **signal koji se menja u vremenu**.  
Filter određuje koje frekvencije propušta, a koje pušta.

**IMPEDANSA (Z)** = "otpor" koji komponenta pruža naizmeničnoj struji, zavisi od frekvencije.

---

## 1. Otpornik – R – "kočničar"

- R = const, R ≠ R(f)
- Podjednako gasi i niske i visoke frekvencije
- Partner kondenzatoru i kalemu
- Određuje oštrinu filtera i postavlja granicu jačine signala

---

## 2. Kondenzator – C

- = dve provodne ploče razdvojene izolatorom
- Otpornost (Xc, kapacitivna reaktansa) opada kako frekvencija raste:

$$X_C = \frac{1}{2\pi f C}$$

- niske f → ogroman Xc ⟹ kondenzator **blokira**
- visoke f → mali Xc ⟹ kondenzator kratko spaja → **propušta ih**

**Zašto?** Viša frekvencija → signal se brže menja, pa se kondenzator brže puni i prazni, ne stigne da pruži otpor signalu.

---

## 3. Kalem – L

- Namotan žica oko koje se stvara magnetno polje kad protiče struja
- Otpornost (induktivna reaktansa XL) raste kako frekvencija raste:

$$X_L = 2\pi f L$$

- niske f → mali XL → žica se samo žica, ima prolaz
- visoke f → ogroman XL → kalem **blokira** visoke tonove

**Zašto?** Brzi signal (visoke f) naglo menja smer/jačinu struje → kalem indukuje suprotnu struju koja se protivi toj promeni *(Lencov zakon)*. Brža promena (f) → veći otpor kalema.

$$\boxed{f_0 = \frac{1}{2\pi\sqrt{LC}}}$$

---

## U audio filterima

### 1. Low-pass – uklanjanje visokih

**A)** Kalem redno  
**B)** Kondenzator paralelno → voli visoke frekvencije → kratko ih spoji prema uzemljenju pre nego što stignu do zvučnika

### 2. High-pass – uklanjanje niskih

**A)** Kondenzator redno  
**B)** Kalem paralelno → kratko će spojiti niske frekvencije ka uzemljenju

---
---

# RLC kao fizički filter vs Digitalni filter koji se ponaša kao RLC

## Pre digitalnih elemenata → samo analogna obrada signala

- RADIO / TV prijemnici
- telekomunikacioni filteri → stvarni otpornici, kondenzatori i zavoji
- pr. Antena – hvatanje stanice na 100 MHz:
  1. antena hvata gomilu frekvencija
  2. RLC kolo je podešeno na 100 MHz → rezonuje baš na toj frekvenciji
  3. ostale f se potiskuju →  band-pass

> **\* digresija → rezonanca → guranje ljuljaške**  
> presporo → malo zavadi, pređe → opet loše  
> baš u ritmu ljuljaške → ogroman zamah  
> RLC = skroz isto → neke f jako pojačuju kolo

---

## RLC kolo

- RLC ima svoju "prirodnu frekvenciju" = rezonanciju  
  → f na kojoj kolo "najlakše osciluje"  
  - ako je f ≈ f₀ → signal prolazi
  - daleko od f₀ → signal se gasi

$$\text{RLC za } t \in \mathbb{R}: \quad x(t) = \frac{1}{C}i(t) + R\frac{di(t)}{dt} + L\frac{d^2i(t)}{dt^2}$$

*Klasični oscilator sa prigušenjem.*

---

## Diskretno vreme → izvodi postaju razlike

$$y[n] = Ax[n] + By[n-1] + Cy[n-2]$$


### Koraci:

**1° Diskretizacija:** `t = nT` → `n = 0, 1, 2, ...`
→ period T = 1 / sample rate

**2° Diferenciranje → razlike:**

$$\frac{dy}{dt} = \frac{y[n] - y[n-1]}{T} \qquad \frac{d^2y}{dt^2} = \frac{y[n] - 2y[n-1] + y[n-2]}{T^2}$$

**3° Ubacivanje:**

$$L\frac{y[n] - 2y[n-1] + y[n-2]}{T^2} + R\frac{y[n] - y[n-1]}{T} + \frac{1}{C}y[n] = x[n]$$

**4° Sređivanje:** `y[n] = Ax[n] + By[n-1] + Cy[n-2]`


- **A** – koliko direktno ulaz utiče
- **B** – koliko sistem "pamti trenutno stanje" → određuje koliko sistem želi da nastavi da osciluje
- **C** – koliko pamti "inerciju kretanja" → određuje stabilnost i pojačanje

---

## Digitalni filter: nema L, nema C → procesor samo računa

$$y[n] = Ax[n] + By[n-1] + Cy[n-2]$$

ali ako se A, B, C izaberu kako treba → **ponaša se identično pravom RLC kolu**.

- **analogni filter** → elektronika rešava jednačinu
- **digitalni filter** → CPU rešava jednačinu

---
---

# Biquad filter

Koristi se za implementaciju svakog ozbiljnog audio EQ-a, low-pass, high-pass filtera itd.

**Opšti oblik:**

$$\boxed{y[n] = a_0 x[n] + a_1 x[n-1] + a_2 x[n-2] + b_1 y[n-1] + b_2 y[n-2]}$$

---

## Za 44.1 kHz semplovanje i granicu od 440 Hz:

**Low-pass:**
$$y[n] = 0{,}059 \cdot x[n] + 0{,}941 \cdot y[n-1]$$

**High-pass:**
$$y[n] = 0{,}941 \cdot \bigl(x[n] - x[n-1]\bigr) + y[n-1]$$

---

## Izvođenje: RLC → biquad

### Za f₀ = 440 Hz:

$$\sqrt{LC} = \frac{1}{2\pi f_0} = 3{,}619 \cdot 10^{-4} \quad \Rightarrow \quad LC = 13{,}097 \cdot 10^{-8}$$

npr. $L = 10^{-2}$ H $\Rightarrow$ $C = 13{,}1 \cdot 10^{-6}$ F

...definišemo da je $R = 10\ \Omega$

**Sample Rate = 44100 Hz**

$$x(t) = V_R(t) + V_L(t) + V_C(t)$$

$$V_R = R \cdot i(t)$$
$$V_L = L\frac{di}{dt}$$
$$V_C = \frac{1}{C}\int i(t)\,dt$$

izlaz: $y(t) = V_R = Ri$ → $i = \frac{y}{R}$

Polazna jednačina:

$$x = y + L\frac{1}{R}\frac{dy}{dt} + \frac{1}{C}\int\frac{y}{R}\,dt$$

$$RC\frac{dx}{dt} = RC\frac{dy}{dt} + LC\frac{d^2y}{dt^2} + y$$

$$\Rightarrow \quad LC \cdot y'' + RC \cdot y' + y = RC \cdot x'$$

Ubacimo → $1{,}31 \cdot 10^{-7} y'' + 9{,}0015^{-4} y' + y = 9{,}0018^{-4} x'$

**Diskretizacija:**

$$y' \approx \frac{y[n] - y[n-1]}{T}, \qquad x' \approx \frac{x[n] - x[n-1]}{T}$$

$$y'' \approx \frac{y[n] - 2y[n-1] + y[n-2]}{T^2}, \qquad x'' \approx \text{ne treba nam}$$

Ubacujemo u diferencijalnu jednačinu i sređujemo:

$$\left(\frac{LC}{T^2} + \frac{RC}{T} + 1\right)y[n] = \frac{RC}{T}x[n] - \frac{RC}{T}x[n-1] + \left(2\frac{LC}{T^2} + \frac{RC}{T}\right)y[n-1] - \frac{LC}{T^2}y[n-2]$$

Nakon ubacivanja vrednosti:

$$y[n] = \frac{5{,}78}{1262}x[n] - \frac{5{,}78}{262}x[n-1] + \frac{2(n?)-5{,}78}{262}y[n-1] - \frac{285{,}8}{802}y[n-2]$$

**Finalna forma:**


$$y[n] = 0{,}02206\bigl(x[n] - x[n-1]\bigr) + 1{,}9691 \cdot y[n-1] - 0{,}9725 \cdot y[n-2]$$

---
