# Petnica Music Workshop - Muzika na struju

## Prerequisites & Installation

Before getting started, make sure you have **Go (Golang)** installed on your system.

### Troubleshooting (Linux Users)
If you run into issues compilation or audio playback errors related to the `oto` library on Linux, you need to install the missing development dependencies. Run the following command in your terminal:

```bash
sudo apt update && sudo apt install pkg-config libasound2-dev
```

## Getting Started
Clone this repository to your local machine:

```bash
git clone [https://github.com/petnica-rac/music-workshop.git](https://github.com/petnica-rac/music-workshop.git)
```

Navigate into the project directory:
```bash
cd music-workshop
```

---


Run the main workshop entry point from the correct directory:

```bash
go mod init <name>
go run main.go
```
or, if in the project dir
```bash
go run ./segmentN/main.go
```

Run go mod tidy frequently to cleanly synchronize your go.mod and go.sum files with the actual source code imports in your .go files
```bash
go mod tidy
```