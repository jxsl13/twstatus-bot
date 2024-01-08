package model_test

import (
	"testing"

	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/require"
)

func TestRuneWidth(t *testing.T) {
	// test cases from
	weirdCharacters := "     Ƥ.I.Ƈ."
	normalCharacters := "[Syndicate]"

	c := runewidth.NewCondition()
	c.EastAsianWidth = false
	c.StrictEmojiNeutral = false

	wl := runewidth.StringWidth(weirdCharacters)
	nl := runewidth.StringWidth(normalCharacters)

	require.NotEqual(t, wl, nl)
}
