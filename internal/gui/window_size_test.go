package gui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"

	"github.com/sumit-poudel/gine/internal/config"
	"github.com/sumit-poudel/gine/internal/mimetype"
)

func TestResizeSwitchesTooSmallScreen(t *testing.T) {
	model := New(config.Config{}, nil, mimetype.Catalog{}, mimetype.Associations{})

	before := model.screen()
	if strings.Contains(before, windowTooSmall) {
		t.Fatal("notice shown before a window size arrives")
	}

	narrow := model.resize(minWindowWidth-1, 40)
	if narrow.hasWindowSize {
		t.Fatal("narrow window marked usable")
	}
	narrowScreen := narrow.screen()
	if !strings.Contains(narrowScreen, windowTooSmall) {
		t.Fatalf("narrow screen = %q", narrowScreen)
	}
	if got := noticeColumn(narrowScreen); got != (minWindowWidth-1-lipgloss.Width(windowTooSmall))/2 {
		t.Fatalf("notice column = %d", got)
	}
	if lipgloss.Height(narrowScreen) != 40 {
		t.Fatalf("narrow screen height = %d", lipgloss.Height(narrowScreen))
	}

	short := model.resize(80, 3)
	if short.hasWindowSize {
		t.Fatal("short window marked usable")
	}
	if !strings.Contains(short.screen(), windowTooSmall) {
		t.Fatal("short window kept the full layout")
	}

	fitted := model.resize(80, 40)
	if !fitted.hasWindowSize {
		t.Fatal("fitted window marked too small")
	}
	if strings.Contains(fitted.screen(), windowTooSmall) {
		t.Fatal("fitted window still shows the notice")
	}
}

func noticeColumn(screen string) int {
	for _, line := range strings.Split(screen, "\n") {
		if strings.Contains(line, windowTooSmall) {
			return lipgloss.Width(line) - lipgloss.Width(windowTooSmall)
		}
	}
	return -1
}
