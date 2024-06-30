package audio

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/okayama-daiki/tetris/assets/bgm"
	"github.com/okayama-daiki/tetris/assets/se"
)

type Player struct {
	audioContext        *audio.Context
	audioPlayer         *audio.Player
	hardDropAudioPlayer *audio.Player
	clearAudioPlayer    *audio.Player
	rotateAudioPlayer   *audio.Player
	moveAudioPlayer     *audio.Player
	holdAudioPlayer     *audio.Player
}

func NewPlayer(audioContext *audio.Context) (*Player, error) {
	s, err := mp3.DecodeWithoutResampling(bytes.NewReader(bgm.Tetriiis))
	if err != nil {
		return nil, err
	}
	mainAudioPlayer, err := audioContext.NewPlayer(s)
	if err != nil {
		return nil, err
	}

	s, err = mp3.DecodeWithoutResampling(bytes.NewReader(se.HardDrop))
	if err != nil {
		return nil, err
	}
	hardDropAudioPlayer, err := audioContext.NewPlayer(s)
	if err != nil {
		return nil, err
	}

	s, err = mp3.DecodeWithoutResampling(bytes.NewReader(se.Clear))
	if err != nil {
		return nil, err
	}
	clearAudioPlayer, err := audioContext.NewPlayer(s)
	if err != nil {
		return nil, err
	}

	s, err = mp3.DecodeWithoutResampling(bytes.NewReader(se.Rotate))
	if err != nil {
		return nil, err
	}
	rotateAudioPlayer, err := audioContext.NewPlayer(s)
	if err != nil {
		return nil, err
	}

	s, err = mp3.DecodeWithoutResampling(bytes.NewReader(se.Move))
	if err != nil {
		return nil, err
	}
	moveAudioPlayer, err := audioContext.NewPlayer(s)
	if err != nil {
		return nil, err
	}

	s, err = mp3.DecodeWithoutResampling(bytes.NewReader(se.Hold))
	if err != nil {
		return nil, err
	}
	holdAudioPlayer, err := audioContext.NewPlayer(s)
	if err != nil {
		return nil, err
	}

	player := &Player{
		audioContext:        audioContext,
		audioPlayer:         mainAudioPlayer,
		hardDropAudioPlayer: hardDropAudioPlayer,
		clearAudioPlayer:    clearAudioPlayer,
		rotateAudioPlayer:   rotateAudioPlayer,
		moveAudioPlayer:     moveAudioPlayer,
		holdAudioPlayer:     holdAudioPlayer,
	}

	player.audioPlayer.Play()
	return player, nil
}

func (p *Player) PlayMove() {
	err := p.moveAudioPlayer.Rewind()
	if err != nil {
		panic(err)
	}
	p.moveAudioPlayer.Play()
}

func (p *Player) PlayRotate() {
	err := p.rotateAudioPlayer.Rewind()
	if err != nil {
		panic(err)
	}
	p.rotateAudioPlayer.Play()
}

func (p *Player) PlayHold() {
	err := p.holdAudioPlayer.Rewind()
	if err != nil {
		panic(err)
	}
	p.holdAudioPlayer.Play()
}

func (p *Player) PlayClear() {
	err := p.clearAudioPlayer.Rewind()
	if err != nil {
		panic(err)
	}
	p.clearAudioPlayer.Play()
}

func (p *Player) PlayHardDrop() {
	err := p.hardDropAudioPlayer.Rewind()
	if err != nil {
		panic(err)
	}
	p.hardDropAudioPlayer.Play()
}

func (p *Player) Update() {
	if p.audioPlayer.IsPlaying() {
		return
	}
	err := p.audioPlayer.Rewind()
	if err != nil {
		panic(err)
	}
	p.audioPlayer.Play()
}
