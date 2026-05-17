package e2e

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const seededRNGSeed = "00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"

type resultSummaryArtifact struct {
	MatchID        string `json:"match_id"`
	GameID         string `json:"game_id"`
	GameVersion    string `json:"game_version"`
	RulesetVersion string `json:"ruleset_version"`
	Status         string `json:"status"`
	Turn           int    `json:"turn"`
	Placements     []struct {
		PlayerID string `json:"player_id"`
		Place    int    `json:"place"`
	} `json:"placements"`
	Dungeon *struct {
		MapID                string `json:"map_id"`
		Turn                 int    `json:"turn"`
		MaxTurns             int    `json:"max_turns"`
		RemainingChestCount  int    `json:"remaining_chest_count"`
		RemainingChestPoints int    `json:"remaining_chest_points"`
		Players              []struct {
			PlayerID     string `json:"player_id"`
			Place        int    `json:"place"`
			Score        int    `json:"score"`
			GoalBonus    int    `json:"goal_bonus"`
			ChestPoints  int    `json:"chest_points"`
			FinishedTurn *int   `json:"finished_turn"`
		} `json:"players"`
		RemainingChests []struct {
			X      int `json:"x"`
			Y      int `json:"y"`
			Points int `json:"points"`
		} `json:"remaining_chests"`
	} `json:"dungeon"`
}

type normalizedDungeonResult struct {
	MapID               string                    `json:"map_id"`
	Turn                int                       `json:"turn"`
	MaxTurns            int                       `json:"max_turns"`
	Placements          []normalizedPlacement     `json:"placements"`
	Players             []normalizedDungeonPlayer `json:"players"`
	RemainingChests     []normalizedChest         `json:"remaining_chests"`
	RemainingChestCount int                       `json:"remaining_chest_count"`
	RemainingChestTotal int                       `json:"remaining_chest_total"`
}

type normalizedPlacement struct {
	PlayerID string `json:"player_id"`
	Place    int    `json:"place"`
}

type normalizedDungeonPlayer struct {
	PlayerID     string `json:"player_id"`
	Place        int    `json:"place"`
	Score        int    `json:"score"`
	GoalBonus    int    `json:"goal_bonus"`
	ChestPoints  int    `json:"chest_points"`
	FinishedTurn int    `json:"finished_turn"`
}

type normalizedChest struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Points int `json:"points"`
}

func TestDungeonLocalRunCompletes(t *testing.T) {
	summary := runSeededDungeonMatch(t, "dungeon-local-completion")
	if summary.Status != "completed" {
		t.Fatalf("status = %q, want completed", summary.Status)
	}
	if summary.GameID != "dungeon" {
		t.Fatalf("game_id = %q, want dungeon", summary.GameID)
	}
	if summary.RulesetVersion != "seeded-maze-v1" {
		t.Fatalf("ruleset = %q, want seeded-maze-v1", summary.RulesetVersion)
	}
	if summary.Dungeon == nil {
		t.Fatal("summary missing dungeon payload")
	}
}

func TestDungeonSeededDeterministicGoldenParity(t *testing.T) {
	first := normalizeDungeonResultPayload(t, runSeededDungeonMatch(t, "dungeon-seeded-deterministic-a"))
	second := normalizeDungeonResultPayload(t, runSeededDungeonMatch(t, "dungeon-seeded-deterministic-b"))

	assertCanonicalJSONEqual(t, first, second, "normalized result mismatch across reruns")
	assertGoldenJSON(t, filepath.Join(repoRoot(t), "e2e", "golden", "normalized-dungeon-result.json"), first)
}

func TestDungeonDeterministicGoldenIsCanonicalJSON(t *testing.T) {
	goldenPath := filepath.Join(repoRoot(t), "e2e", "golden", "normalized-dungeon-result.json")
	data, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	var payload normalizedDungeonResult
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("decode golden: %v", err)
	}
	if string(data) != string(mustIndentedJSON(payload)) {
		t.Fatalf("golden %s is not canonical pretty JSON", goldenPath)
	}
}

func runSeededDungeonMatch(t *testing.T, matchID string) resultSummaryArtifact {
	t.Helper()

	outputDir := t.TempDir()
	cmd := exec.CommandContext(testContext(t),
		"go", "run", "github.com/yoskeoka/ai-arena/cmd/arena-runner",
		"--game", "dungeon",
		"--game-version", "1.0.0",
		"--ruleset", "seeded-maze-v1",
		"--rng-seed", seededRNGSeed,
		"--match-id", matchID,
		"--output-dir", outputDir,
		"--log-output", "none",
		"--player", "p1=./testdata/ai/dungeon/dungeon-bot-local-seeded",
		"--player", "p2=./testdata/ai/dungeon/dungeon-bot-local-seeded",
	)
	cmd.Dir = repoRoot(t)
	cmd.Env = childEnvWithLocalWorkspace(t)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run arena-runner: %v\n%s", err, output)
	}
	return readResultSummary(t, filepath.Join(outputDir, matchID))
}

func readResultSummary(t *testing.T, matchDir string) resultSummaryArtifact {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(matchDir, "result-summary.json"))
	if err != nil {
		t.Fatalf("read result summary: %v", err)
	}
	var summary resultSummaryArtifact
	if err := json.Unmarshal(data, &summary); err != nil {
		t.Fatalf("decode result summary: %v\nsummary=%s", err, data)
	}
	return summary
}

func normalizeDungeonResultPayload(t *testing.T, summary resultSummaryArtifact) normalizedDungeonResult {
	t.Helper()

	if summary.Dungeon == nil {
		t.Fatal("summary missing dungeon payload")
	}

	got := normalizedDungeonResult{
		MapID:               summary.Dungeon.MapID,
		Turn:                summary.Dungeon.Turn,
		MaxTurns:            summary.Dungeon.MaxTurns,
		RemainingChestCount: summary.Dungeon.RemainingChestCount,
		RemainingChestTotal: summary.Dungeon.RemainingChestPoints,
	}
	for _, placement := range summary.Placements {
		got.Placements = append(got.Placements, normalizedPlacement{
			PlayerID: placement.PlayerID,
			Place:    placement.Place,
		})
	}
	for _, player := range summary.Dungeon.Players {
		got.Players = append(got.Players, normalizedDungeonPlayer{
			PlayerID:     player.PlayerID,
			Place:        player.Place,
			Score:        player.Score,
			GoalBonus:    player.GoalBonus,
			ChestPoints:  player.ChestPoints,
			FinishedTurn: finishedTurnValue(player.FinishedTurn),
		})
	}
	for _, chest := range summary.Dungeon.RemainingChests {
		got.RemainingChests = append(got.RemainingChests, normalizedChest{
			X:      chest.X,
			Y:      chest.Y,
			Points: chest.Points,
		})
	}
	sortRemainingChests(got.RemainingChests)
	return got
}

func sortRemainingChests(chests []normalizedChest) {
	for i := 0; i < len(chests); i++ {
		for j := i + 1; j < len(chests); j++ {
			if chests[j].X < chests[i].X ||
				(chests[j].X == chests[i].X && chests[j].Y < chests[i].Y) ||
				(chests[j].X == chests[i].X && chests[j].Y == chests[i].Y && chests[j].Points < chests[i].Points) {
				chests[i], chests[j] = chests[j], chests[i]
			}
		}
	}
}

func finishedTurnValue(v *int) int {
	if v == nil {
		return -1
	}
	return *v
}

func assertGoldenJSON(t *testing.T, goldenPath string, got any) {
	t.Helper()

	wantData, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden %s: %v", goldenPath, err)
	}
	var want normalizedDungeonResult
	if err := json.Unmarshal(wantData, &want); err != nil {
		t.Fatalf("decode golden %s: %v", goldenPath, err)
	}
	assertCanonicalJSONEqual(t, want, got, "golden mismatch for "+goldenPath)
}

func assertCanonicalJSONEqual(t *testing.T, left, right any, message string) {
	t.Helper()

	leftJSON := mustIndentedJSON(left)
	rightJSON := mustIndentedJSON(right)
	if string(leftJSON) != string(rightJSON) {
		t.Fatalf("%s\nleft:\n%s\nright:\n%s", message, leftJSON, rightJSON)
	}
}

func mustIndentedJSON(v any) []byte {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return append(data, '\n')
}

func repoRoot(t *testing.T) string {
	t.Helper()

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	return strings.TrimSpace(string(output))
}

func testContext(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)
	return ctx
}

func childEnvWithLocalWorkspace(t *testing.T) []string {
	t.Helper()

	env := append([]string(nil), os.Environ()...)
	repo := repoRoot(t)
	if aiArenaDir := detectAIArenaDir(repo); aiArenaDir != "" {
		workDir := t.TempDir()
		workFile := filepath.Join(workDir, "go.work")
		data := "go 1.26\n\nuse (\n  " + repo + "\n  " + aiArenaDir + "\n)\n"
		if err := os.WriteFile(workFile, []byte(data), 0o644); err != nil {
			t.Fatalf("write go.work: %v", err)
		}
		env = append(env, "GOWORK="+workFile)
	}
	return env
}

func detectAIArenaDir(repo string) string {
	candidates := []string{}
	if env := strings.TrimSpace(os.Getenv("AI_ARENA_DIR")); env != "" {
		candidates = append(candidates, env)
	}
	candidates = append(candidates,
		filepath.Join(repo, "ai-arena"),
		filepath.Join(repo, "..", "ai-arena"),
		filepath.Join(repo, "..", "..", "ai-arena"),
	)
	for _, candidate := range candidates {
		if _, err := os.Stat(filepath.Join(candidate, "go.mod")); err == nil {
			return candidate
		}
	}
	return ""
}
