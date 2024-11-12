package game

import (
	"testing"
)

func TestRotate(t *testing.T) {
	got := Rotate([][]int{
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 1},
	})
	want := [][]int{
		{0, 0, 0},
		{1, 1, 1},
		{1, 0, 0},
	}

	for y := range len(got) {
		for x := range len(got[y]) {
			if got[y][x] != want[y][x] {
				t.Errorf("got %v, want %v", got, want)
			}
		}
	}

}
