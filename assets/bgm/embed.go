package bgm

import (
	_ "embed"
)

var (
	//go:embed tetriiis.mp3
	Tetriiis []byte

	//go:embed theme.mp3
	Theme []byte
)
