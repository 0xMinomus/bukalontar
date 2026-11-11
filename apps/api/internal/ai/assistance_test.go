package ai

import (
	"testing"

	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
)

func TestSanitizeBalineseDropsLookalikeScripts(t *testing.T) {
	paleography := model.Paleography{BalineseTranscription: "ᬅᬓ᭄ꦲꦥᬤ᭄ᬤ 3 ᬲᬂ"}
	if removed := sanitizeBalinese(&paleography); removed != 2 {
		t.Fatalf("removed %d lookalike characters, want 2", removed)
	}
	if paleography.BalineseTranscription != "ᬅᬓ᭄ᬤ᭄ᬤ 3 ᬲᬂ" {
		t.Fatalf("unexpected cleaned transcription: %q", paleography.BalineseTranscription)
	}
}
