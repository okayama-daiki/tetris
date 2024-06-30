package se

import (
	_ "embed"
)

var (
	//go:embed clear.mp3
	Clear []byte

	//go:embed hard-drop.mp3
	HardDrop []byte

	//go:embed hold.mp3
	Hold []byte

	//go:embed move.mp3
	Move []byte

	//go:embed rotate.mp3
	Rotate []byte
)
