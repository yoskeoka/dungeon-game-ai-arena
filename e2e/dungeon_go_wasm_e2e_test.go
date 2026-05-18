package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/yoskeoka/dungeon-game-ai-arena/games/dungeon"
)

func TestDungeonGoWASMTaggedRunnerCompletes(t *testing.T) {
	requireWASME2E(t)

	entryPath := prepareGoWASMFixture(t)
	outputDir := t.TempDir()
	matchID := "dungeon-go-wasm-tagged-runner"

	cmd := exec.CommandContext(testContext(t),
		"go", "run", "github.com/yoskeoka/ai-arena/cmd/arena-runner@"+aiArenaVersion(),
		"--game", dungeon.GameID,
		"--game-version", dungeon.GameVersion,
		"--ruleset", dungeon.RulesetSeededMazeV1,
		"--rng-seed", seededRNGSeed,
		"--match-id", matchID,
		"--output-dir", outputDir,
		"--log-output", "none",
		"--player", "p1="+entryPath,
		"--player", "p2=./testdata/ai/dungeon/dungeon-bot-local-seeded",
	)
	cmd.Dir = repoRoot(t)
	cmd.Env = commandEnv()
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run tagged runner with dungeon Go-WASM player: %v\n%s", err, output)
	}

	summary := readResultSummary(t, filepath.Join(outputDir, matchID))
	if summary.Status != "completed" {
		t.Fatalf("status = %q, want completed", summary.Status)
	}
	if summary.GameID != dungeon.GameID {
		t.Fatalf("game_id = %q, want %q", summary.GameID, dungeon.GameID)
	}
	if summary.GameVersion != dungeon.GameVersion {
		t.Fatalf("game_version = %q, want %q", summary.GameVersion, dungeon.GameVersion)
	}
	if summary.RulesetVersion != dungeon.RulesetSeededMazeV1 {
		t.Fatalf("ruleset = %q, want %q", summary.RulesetVersion, dungeon.RulesetSeededMazeV1)
	}
	if summary.Dungeon == nil {
		t.Fatal("summary missing dungeon payload")
	}
}

func requireWASME2E(t *testing.T) {
	t.Helper()

	if os.Getenv("AI_ARENA_WASM_E2E") != "1" {
		t.Skip("set AI_ARENA_WASM_E2E=1 to enable WASM verification tests")
	}
}

func prepareGoWASMFixture(t *testing.T) string {
	t.Helper()

	repo := repoRoot(t)
	root, err := os.MkdirTemp(repo, ".tmp-dungeon-go-wasm-")
	if err != nil {
		t.Fatalf("create temp wasm fixture dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(root)
	})
	entryPath := filepath.Join(root, "dungeon-go-wasm-ai")
	outputPath := entryPath + ".wasm"
	manifestPath := entryPath + ".arena.json"
	manifest := map[string]any{
		"ai_id": "dungeon-go-wasm-ai",
		"protocol": map[string]any{
			"transport":       "stdio-jsonrpc-ndjson",
			"game_id":         dungeon.GameID,
			"game_version":    dungeon.GameVersion,
			"ruleset_version": dungeon.RulesetSeededMazeV1,
		},
		"runtime": map[string]any{
			"kind":               "wasm-wasi",
			"module":             "./dungeon-go-wasm-ai.wasm",
			"args":               []string{"./dungeon-go-wasm-ai.wasm"},
			"memory_limit_pages": 1024,
		},
	}
	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("marshal wasm manifest: %v", err)
	}
	manifestBytes = append(manifestBytes, '\n')
	if err := os.WriteFile(manifestPath, manifestBytes, 0o644); err != nil {
		t.Fatalf("write wasm manifest: %v", err)
	}

	cmd := exec.CommandContext(testContext(t), "go", "build", "-o", outputPath, "./testdata/ai/dungeon/dungeon-go-wasm-ai")
	cmd.Dir = repo
	cmd.Env = append(commandEnv(), "GOOS=wasip1", "GOARCH=wasm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build dungeon Go-WASM fixture: %v\n%s", err, output)
	}
	relativeEntry, err := filepath.Rel(repo, entryPath)
	if err != nil {
		t.Fatalf("compute relative wasm entry path: %v", err)
	}
	return "./" + filepath.ToSlash(relativeEntry)
}
