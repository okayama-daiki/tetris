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
	errors := make([]error, 0)

	load := func(b []byte) *audio.Player {
		s, err := mp3.DecodeWithoutResampling(bytes.NewReader(b))
		if err != nil {
			errors = append(errors, err)
			return nil
		}
		player, err := audioContext.NewPlayer(s)
		if err != nil {
			errors = append(errors, err)
			return nil
		}
		return player
	}

	player := &Player{
		audioContext:        audioContext,
		audioPlayer:         load(bgm.Theme),
		hardDropAudioPlayer: load(se.HardDrop),
		clearAudioPlayer:    load(se.Clear),
		rotateAudioPlayer:   load(se.Rotate),
		moveAudioPlayer:     load(se.Move),
		holdAudioPlayer:     load(se.Hold),
	}

	if len(errors) > 0 {
		return nil, errors[0]
	}

	player.audioPlayer.Play()
	return player, nil
}

func _play(player *audio.Player) {
	err := player.Rewind()
	if err != nil {
		panic(err)
	}
	player.Play()
}

func (p *Player) PlayMove() {
	_play(p.moveAudioPlayer)
}

func (p *Player) PlayRotate() {
	_play(p.rotateAudioPlayer)
}

func (p *Player) PlayHold() {
	_play(p.holdAudioPlayer)
}

func (p *Player) PlayClear() {
	_play(p.clearAudioPlayer)
}

func (p *Player) PlayHardDrop() {
	_play(p.hardDropAudioPlayer)
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
