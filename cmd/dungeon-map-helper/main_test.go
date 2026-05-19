package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yoskeoka/dungeon-game-ai-arena/games/dungeon"
)

func TestRunPrintsMultipleRulesetSummaries(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := runWithIO([]string{
		"--rng-seed", dungeon.DefaultRNGSeed,
		"--ruleset", dungeon.RulesetSeededMazeV1,
		"--ruleset", dungeon.RulesetRogueRoomsV1,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("runWithIO: %v", err)
	}
	output := stdout.String()
	for _, ruleset := range []string{dungeon.RulesetSeededMazeV1, dungeon.RulesetRogueRoomsV1} {
		if !strings.Contains(output, "ruleset="+ruleset) {
			t.Fatalf("output missing ruleset %q:\n%s", ruleset, output)
		}
	}
	if !strings.Contains(output, "shape width=") {
		t.Fatalf("output missing shape summary:\n%s", output)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}
