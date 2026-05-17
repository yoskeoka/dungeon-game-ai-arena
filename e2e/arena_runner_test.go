package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/yoskeoka/dungeon-game-ai-arena/games/dungeon"
)

const aiArenaVersion = "v0.1.0"

type resultSummary struct {
	RulesetVersion string `json:"ruleset_version"`
	Status         string `json:"status"`
	Turn           int    `json:"turn"`
	Placements     []struct {
		PlayerID string `json:"player_id"`
		Place    int    `json:"place"`
	} `json:"placements"`
	Dungeon *struct {
		MapID           string `json:"map_id"`
		Turn            int    `json:"turn"`
		MaxTurns        int    `json:"max_turns"`
		RemainingChests []struct {
			X      int `json:"x"`
			Y      int `json:"y"`
			Points int `json:"points"`
		} `json:"remaining_chests"`
		RemainingChestCount  int `json:"remaining_chest_count"`
		RemainingChestPoints int `json:"remaining_chest_points"`
		Players              []struct {
			PlayerID     string `json:"player_id"`
			Place        int    `json:"place"`
			Score        int    `json:"score"`
			GoalBonus    int    `json:"goal_bonus"`
			ChestPoints  int    `json:"chest_points"`
			FinishedTurn *int   `json:"finished_turn"`
		} `json:"players"`
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

func TestArenaRunnerDungeonSeededReferenceBotDeterministicResultRegression(t *testing.T) {
	first := runSeededDungeonDeterministicRegression(t, "dungeon-seeded-deterministic-a")
	second := runSeededDungeonDeterministicRegression(t, "dungeon-seeded-deterministic-b")

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("normalized result mismatch across reruns\nfirst=%s\nsecond=%s",
			mustIndentedJSON(first), mustIndentedJSON(second))
	}
	assertGoldenJSON(t, filepath.Join("golden", "normalized-dungeon-result.json"), first)
}

func TestArenaRunnerDungeonDeterministicGoldenIsCanonicalJSON(t *testing.T) {
	goldenPath := filepath.Join("golden", "normalized-dungeon-result.json")
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

func runSeededDungeonDeterministicRegression(t *testing.T, matchID string) normalizedDungeonResult {
	t.Helper()

	repoRoot := repoRoot(t)
	outputDir := t.TempDir()
	args := []string{
		"run", "github.com/yoskeoka/ai-arena/cmd/arena-runner@" + aiArenaVersion,
		"--game", dungeon.GameID,
		"--game-version", dungeon.GameVersion,
		"--ruleset", dungeon.RulesetSeededMazeV1,
		"--rng-seed", "00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff",
		"--match-id", matchID,
		"--output-dir", outputDir,
		"--player", "p1=./testdata/ai/dungeon/dungeon-bot-local-seeded",
		"--player", "p2=./testdata/ai/dungeon/dungeon-bot-local-seeded",
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "GOWORK=off")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("arena-runner failed: %v\noutput=%s", err, output)
	}

	matchDir := filepath.Join(outputDir, matchID)
	data, err := os.ReadFile(filepath.Join(matchDir, "result-summary.json"))
	if err != nil {
		t.Fatalf("read result summary: %v", err)
	}
	var summary resultSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		t.Fatalf("decode result summary: %v\nsummary=%s", err, data)
	}
	if summary.Status != "completed" {
		t.Fatalf("status = %q, want completed", summary.Status)
	}
	if summary.RulesetVersion != dungeon.RulesetSeededMazeV1 {
		t.Fatalf("ruleset = %q, want %q", summary.RulesetVersion, dungeon.RulesetSeededMazeV1)
	}
	return normalizeDungeonResult(t, summary)
}

func normalizeDungeonResult(t *testing.T, summary resultSummary) normalizedDungeonResult {
	t.Helper()

	if summary.Dungeon == nil {
		t.Fatal("summary missing dungeon payload")
	}
	if summary.Turn != summary.Dungeon.Turn {
		t.Fatalf("summary turn mismatch: top-level=%d dungeon=%d", summary.Turn, summary.Dungeon.Turn)
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
	sort.Slice(got.RemainingChests, func(i, j int) bool {
		if got.RemainingChests[i].X != got.RemainingChests[j].X {
			return got.RemainingChests[i].X < got.RemainingChests[j].X
		}
		if got.RemainingChests[i].Y != got.RemainingChests[j].Y {
			return got.RemainingChests[i].Y < got.RemainingChests[j].Y
		}
		return got.RemainingChests[i].Points < got.RemainingChests[j].Points
	})
	return got
}

func assertGoldenJSON(t *testing.T, goldenPath string, got normalizedDungeonResult) {
	t.Helper()

	wantData, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden %s: %v", goldenPath, err)
	}
	var want normalizedDungeonResult
	if err := json.Unmarshal(wantData, &want); err != nil {
		t.Fatalf("decode golden %s: %v", goldenPath, err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("golden mismatch for %s\nwant=%s\ngot=%s", goldenPath, mustIndentedJSON(want), mustIndentedJSON(got))
	}
}

func finishedTurnValue(v *int) int {
	if v == nil {
		return -1
	}
	return *v
}

func repoRoot(t *testing.T) string {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if filepath.Base(cwd) == "e2e" {
		return filepath.Dir(cwd)
	}
	if strings.HasSuffix(cwd, string(filepath.Separator)+"e2e") {
		return filepath.Dir(cwd)
	}
	return cwd
}

func mustIndentedJSON(v any) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(append(data, '\n'))
}
