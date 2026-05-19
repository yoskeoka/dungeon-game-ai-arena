package dungeon

import (
	"fmt"
	"testing"
)

const (
	testSeedAlpha = "00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"
	testSeedBeta  = "ffeeddccbbaa99887766554433221100ffeeddccbbaa99887766554433221100"
)

func TestSeededMazeGenerationIsDeterministic(t *testing.T) {
	for _, ruleset := range []string{RulesetSeededMazeV1, RulesetRogueRoomsV1} {
		cfg := Config{
			GameVersion: GameVersion,
			Ruleset:     ruleset,
			PlayerIDs:   []string{"p1", "p2"},
			RNGSeed:     testSeedAlpha,
		}
		first, err := New(cfg)
		if err != nil {
			t.Fatalf("%s first New: %v", ruleset, err)
		}
		second, err := New(cfg)
		if err != nil {
			t.Fatalf("%s second New: %v", ruleset, err)
		}
		if !equalStringSlices(first.Layout().Tiles, second.Layout().Tiles) {
			t.Fatalf("%s tiles differ for same seed", ruleset)
		}
		if !equalPositions(first.Layout().SpawnPoints, second.Layout().SpawnPoints) {
			t.Fatalf("%s spawn points differ for same seed", ruleset)
		}
		if first.Layout().Goal != second.Layout().Goal {
			t.Fatalf("%s goal differs for same seed", ruleset)
		}
		if !equalChests(first.Layout().InitialChests, second.Layout().InitialChests) {
			t.Fatalf("%s initial chests differ for same seed", ruleset)
		}
	}
}

func TestSeededMazeGenerationVariesAcrossSeeds(t *testing.T) {
	for _, ruleset := range []string{RulesetSeededMazeV1, RulesetRogueRoomsV1} {
		first, err := New(Config{
			GameVersion: GameVersion,
			Ruleset:     ruleset,
			PlayerIDs:   []string{"p1", "p2"},
			RNGSeed:     testSeedAlpha,
		})
		if err != nil {
			t.Fatalf("%s New alpha: %v", ruleset, err)
		}
		second, err := New(Config{
			GameVersion: GameVersion,
			Ruleset:     ruleset,
			PlayerIDs:   []string{"p1", "p2"},
			RNGSeed:     testSeedBeta,
		})
		if err != nil {
			t.Fatalf("%s New beta: %v", ruleset, err)
		}
		sameTiles := equalStringSlices(first.Layout().Tiles, second.Layout().Tiles)
		sameGoal := first.Layout().Goal == second.Layout().Goal
		sameChests := equalChests(first.Layout().InitialChests, second.Layout().InitialChests)
		if sameTiles && sameGoal && sameChests {
			t.Fatalf("%s expected different generated state for different seeds", ruleset)
		}
	}
}

func TestSeededMazeUsesFixedChestScoreSet(t *testing.T) {
	for _, ruleset := range []string{RulesetSeededMazeV1, RulesetRogueRoomsV1} {
		match, err := New(Config{
			GameVersion: GameVersion,
			Ruleset:     ruleset,
			PlayerIDs:   []string{"p1", "p2"},
			RNGSeed:     testSeedAlpha,
		})
		if err != nil {
			t.Fatalf("%s New: %v", ruleset, err)
		}
		total := 0
		got := append([]ChestState(nil), match.Layout().InitialChests...)
		for _, chest := range got {
			total += chest.Points
		}
		if total != 54 {
			t.Fatalf("%s total chest points = %d, want 54", ruleset, total)
		}
		expected := map[int]int{24: 1, 18: 1, 12: 1}
		for _, chest := range got {
			expected[chest.Points]--
		}
		for points, count := range expected {
			if count != 0 {
				t.Fatalf("%s score set mismatch for %d: remaining %d", ruleset, points, count)
			}
		}
	}
}

func TestSeededMazeUsesExpandedTurnBudget(t *testing.T) {
	for _, tc := range []struct {
		ruleset string
		turns   int
	}{
		{ruleset: RulesetSeededMazeV1, turns: seededMazeMaxTurns},
		{ruleset: RulesetRogueRoomsV1, turns: rogueRoomsMaxTurns},
	} {
		match, err := New(Config{
			GameVersion: GameVersion,
			Ruleset:     tc.ruleset,
			PlayerIDs:   []string{"p1", "p2"},
			RNGSeed:     testSeedAlpha,
		})
		if err != nil {
			t.Fatalf("%s New: %v", tc.ruleset, err)
		}
		if match.Ruleset().MaxTurns != tc.turns {
			t.Fatalf("%s max turns = %d, want %d", tc.ruleset, match.Ruleset().MaxTurns, tc.turns)
		}
		visible, err := match.CurrentVisibleState("p1")
		if err != nil {
			t.Fatalf("%s CurrentVisibleState: %v", tc.ruleset, err)
		}
		if visible.RemainingTurns != tc.turns {
			t.Fatalf("%s remaining turns = %d, want %d", tc.ruleset, visible.RemainingTurns, tc.turns)
		}
	}
}

func TestGeneratedRulesetsUseWallBoundedConnectedLayouts(t *testing.T) {
	for _, ruleset := range []string{RulesetSeededMazeV1, RulesetRogueRoomsV1} {
		match, err := New(Config{
			GameVersion: GameVersion,
			Ruleset:     ruleset,
			PlayerIDs:   []string{"p1", "p2"},
			RNGSeed:     testSeedAlpha,
		})
		if err != nil {
			t.Fatalf("%s New: %v", ruleset, err)
		}
		layout := match.Layout()
		assertWallBoundary(t, ruleset, layout.Tiles)
		assertReachable(t, ruleset, match, layout.SpawnPoints[0], layout.Goal)
		for _, spawn := range layout.SpawnPoints[1:] {
			assertReachable(t, ruleset, match, layout.SpawnPoints[0], spawn)
		}
		for _, chest := range layout.InitialChests {
			assertReachable(t, ruleset, match, layout.SpawnPoints[0], Position{X: chest.X, Y: chest.Y})
		}
		assertDistinctPlacements(t, ruleset, layout)
	}
}

func TestGenerateRogueRoomsLayoutRejectsTooSmallDimensions(t *testing.T) {
	rng, err := newSeededRand(testSeedAlpha)
	if err != nil {
		t.Fatalf("newSeededRand: %v", err)
	}
	if _, err := generateRogueRoomsLayout(rng, 7, 15); err == nil {
		t.Fatal("expected width validation error")
	}
	if _, err := generateRogueRoomsLayout(rng, 19, 7); err == nil {
		t.Fatal("expected height validation error")
	}
}

func TestNewGeneratesSeedWhenOmitted(t *testing.T) {
	match, err := New(Config{
		GameVersion: GameVersion,
		Ruleset:     RulesetSeededMazeV1,
		PlayerIDs:   []string{"p1", "p2"},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if len(match.FullState().RNGSeed) != 64 {
		t.Fatalf("generated rng seed length = %d, want 64", len(match.FullState().RNGSeed))
	}
}

func TestFixedMapRulesetRemainsResumable(t *testing.T) {
	match, err := New(Config{
		GameVersion: GameVersion,
		Ruleset:     RulesetFixedMapV1,
		PlayerIDs:   []string{"p1", "p2"},
		RNGSeed:     DefaultRNGSeed,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	state := match.FullState()
	restored, err := NewFromFullState(Config{
		GameVersion: GameVersion,
		Ruleset:     RulesetFixedMapV1,
		PlayerIDs:   []string{"p1", "p2"},
		RNGSeed:     DefaultRNGSeed,
	}, state)
	if err != nil {
		t.Fatalf("NewFromFullState: %v", err)
	}
	if !equalStringSlices(restored.Layout().Tiles, match.Layout().Tiles) {
		t.Fatal("restored tiles differ")
	}
	if !equalChests(restored.Layout().InitialChests, match.Layout().InitialChests) {
		t.Fatal("restored chests differ")
	}
}

func TestNewFromFullStateValidatesGeneratedSeed(t *testing.T) {
	match, err := New(Config{
		GameVersion: GameVersion,
		Ruleset:     RulesetSeededMazeV1,
		PlayerIDs:   []string{"p1", "p2"},
		RNGSeed:     testSeedAlpha,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	state := match.FullState()
	state.RNGSeed = testSeedBeta
	if _, err := NewFromFullState(Config{
		GameVersion: GameVersion,
		Ruleset:     RulesetSeededMazeV1,
		PlayerIDs:   []string{"p1", "p2"},
		RNGSeed:     testSeedAlpha,
	}, state); err == nil {
		t.Fatal("expected rng seed mismatch")
	}
}

func TestNewFromFullStateUsesSnapshotSeedWhenConfigOmitted(t *testing.T) {
	match, err := New(Config{
		GameVersion: GameVersion,
		Ruleset:     RulesetSeededMazeV1,
		PlayerIDs:   []string{"p1", "p2"},
		RNGSeed:     testSeedAlpha,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	state := match.FullState()
	restored, err := NewFromFullState(Config{
		GameVersion: GameVersion,
		Ruleset:     RulesetSeededMazeV1,
		PlayerIDs:   []string{"p1", "p2"},
	}, state)
	if err != nil {
		t.Fatalf("NewFromFullState: %v", err)
	}
	if restored.FullState().RNGSeed != state.RNGSeed {
		t.Fatalf("restored rng seed = %q, want %q", restored.FullState().RNGSeed, state.RNGSeed)
	}
}

func TestNewRejectsInvalidSeedFormat(t *testing.T) {
	if _, err := New(Config{
		GameVersion: GameVersion,
		Ruleset:     RulesetSeededMazeV1,
		PlayerIDs:   []string{"p1", "p2"},
		RNGSeed:     "alpha",
	}); err == nil {
		t.Fatal("expected invalid seed format error")
	}
}

func TestChestSplitAndGoalBonusesStillApply(t *testing.T) {
	match, err := New(Config{
		GameVersion: GameVersion,
		Ruleset:     RulesetFixedMapV1,
		PlayerIDs:   []string{"p1", "p2", "p3"},
		RNGSeed:     DefaultRNGSeed,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	layout := match.Layout()
	restored, err := NewFromFullState(Config{
		GameVersion: GameVersion,
		Ruleset:     RulesetFixedMapV1,
		PlayerIDs:   []string{"p1", "p2", "p3"},
		RNGSeed:     DefaultRNGSeed,
	}, FullState{
		MapID:         RulesetFixedMapV1,
		RNGSeed:       DefaultRNGSeed,
		Turn:          5,
		MaxTurns:      match.Ruleset().MaxTurns,
		Tiles:         append([]string(nil), layout.Tiles...),
		SpawnPoints:   append([]Position(nil), layout.SpawnPoints...),
		Goal:          layout.Goal,
		InitialChests: append([]ChestState(nil), layout.InitialChests...),
		Players: []PlayerState{
			{PlayerID: "p1", X: 1, Y: 6},
			{PlayerID: "p2", X: 3, Y: 6},
			{PlayerID: "p3", X: 5, Y: 6},
		},
		UncollectedChests: []ChestState{
			{X: 2, Y: 6, Points: 12},
		},
		Discovery: map[string]DiscoveryState{
			"p1": {},
			"p2": {},
			"p3": {},
		},
	})
	if err != nil {
		t.Fatalf("NewFromFullState: %v", err)
	}
	if err := restored.Apply(map[string]Action{
		"p1": {Action: "move", Direction: "right"},
		"p2": {Action: "move", Direction: "left"},
		"p3": {Action: "wait"},
	}); err != nil {
		t.Fatalf("Apply contested chest turn: %v", err)
	}
	if err := restored.Apply(map[string]Action{
		"p1": {Action: "move", Direction: "right"},
		"p2": {Action: "move", Direction: "right"},
		"p3": {Action: "move", Direction: "right"},
	}); err != nil {
		t.Fatalf("Apply first finish turn: %v", err)
	}
	if err := restored.Apply(map[string]Action{
		"p1": {Action: "move", Direction: "right"},
		"p2": {Action: "move", Direction: "right"},
	}); err != nil {
		t.Fatalf("Apply advance turn: %v", err)
	}
	if err := restored.Apply(map[string]Action{
		"p1": {Action: "move", Direction: "right"},
		"p2": {Action: "move", Direction: "right"},
	}); err != nil {
		t.Fatalf("Apply advance turn 2: %v", err)
	}
	if err := restored.Apply(map[string]Action{
		"p1": {Action: "move", Direction: "right"},
		"p2": {Action: "move", Direction: "right"},
	}); err != nil {
		t.Fatalf("Apply second finish turn: %v", err)
	}
	players := restored.scoreboardWithPositions()
	want := map[string]struct {
		chest int
		goal  int
		score int
	}{
		"p1": {chest: 6, goal: 50, score: 56},
		"p2": {chest: 6, goal: 50, score: 56},
		"p3": {chest: 0, goal: 100, score: 100},
	}
	for _, player := range players {
		expected := want[player.PlayerID]
		if player.ChestPoints != expected.chest || player.GoalBonus != expected.goal || player.Score != expected.score {
			t.Fatalf("%s = chest:%d goal:%d score:%d, want chest:%d goal:%d score:%d",
				player.PlayerID, player.ChestPoints, player.GoalBonus, player.Score,
				expected.chest, expected.goal, expected.score)
		}
	}
	placements := restored.Placements()
	if placements[0].Place != 1 || placements[1].Place != 2 || placements[2].Place != 2 {
		t.Fatalf("placements = %+v, want competition ranking 1,2,2", placements)
	}
}

func TestThirdPlaceWithMajorityChestPointsBeatsChestlessWinner(t *testing.T) {
	match, err := New(Config{
		GameVersion: GameVersion,
		Ruleset:     RulesetSeededMazeV1,
		PlayerIDs:   []string{"p1", "p2", "p3"},
		RNGSeed:     testSeedAlpha,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	layout := match.Layout()
	path, ok := match.ShortestPath(layout.SpawnPoints[0], layout.Goal)
	if !ok || len(path) < 4 {
		t.Fatalf("shortest path length = %d, want >= 4", len(path))
	}
	p1Start := path[len(path)-2]
	p2Start := path[len(path)-3]
	p3Start := path[len(path)-4]
	directionTo := func(from, to Position) string {
		switch {
		case to.X == from.X && to.Y == from.Y-1:
			return "up"
		case to.X == from.X && to.Y == from.Y+1:
			return "down"
		case to.X == from.X-1 && to.Y == from.Y:
			return "left"
		case to.X == from.X+1 && to.Y == from.Y:
			return "right"
		default:
			t.Fatalf("non-adjacent move from %+v to %+v", from, to)
			return ""
		}
	}
	restored, err := NewFromFullState(Config{
		GameVersion: GameVersion,
		Ruleset:     RulesetSeededMazeV1,
		PlayerIDs:   []string{"p1", "p2", "p3"},
		RNGSeed:     testSeedAlpha,
	}, FullState{
		MapID:         RulesetSeededMazeV1,
		RNGSeed:       testSeedAlpha,
		Turn:          5,
		MaxTurns:      match.Ruleset().MaxTurns,
		Tiles:         append([]string(nil), layout.Tiles...),
		SpawnPoints:   append([]Position(nil), layout.SpawnPoints...),
		Goal:          layout.Goal,
		InitialChests: append([]ChestState(nil), layout.InitialChests...),
		Players: []PlayerState{
			{PlayerID: "p1", X: p1Start.X, Y: p1Start.Y},
			{PlayerID: "p2", X: p2Start.X, Y: p2Start.Y},
			{PlayerID: "p3", X: p3Start.X, Y: p3Start.Y, Score: 30, ChestPoints: 30},
		},
		UncollectedChests: []ChestState{},
		Discovery: map[string]DiscoveryState{
			"p1": {},
			"p2": {},
			"p3": {},
		},
	})
	if err != nil {
		t.Fatalf("NewFromFullState: %v", err)
	}
	layout = match.Layout()
	if err := restored.Apply(map[string]Action{
		"p1": {Action: "move", Direction: directionTo(p1Start, layout.Goal)},
		"p2": {Action: "move", Direction: directionTo(p2Start, p1Start)},
		"p3": {Action: "move", Direction: directionTo(p3Start, p2Start)},
	}); err != nil {
		t.Fatalf("Apply first finish turn: %v", err)
	}
	if err := restored.Apply(map[string]Action{
		"p2": {Action: "move", Direction: directionTo(p1Start, layout.Goal)},
		"p3": {Action: "move", Direction: directionTo(p2Start, p1Start)},
	}); err != nil {
		t.Fatalf("Apply second finish turn: %v", err)
	}
	if err := restored.Apply(map[string]Action{
		"p3": {Action: "move", Direction: directionTo(p1Start, layout.Goal)},
	}); err != nil {
		t.Fatalf("Apply third finish turn: %v", err)
	}
	players := restored.scoreboardWithPositions()
	want := map[string]struct {
		chest int
		goal  int
		score int
	}{
		"p1": {chest: 0, goal: 42, score: 42},
		"p2": {chest: 0, goal: 28, score: 28},
		"p3": {chest: 30, goal: 14, score: 44},
	}
	for _, player := range players {
		expected := want[player.PlayerID]
		if player.ChestPoints != expected.chest || player.GoalBonus != expected.goal || player.Score != expected.score {
			t.Fatalf("%s = chest:%d goal:%d score:%d, want chest:%d goal:%d score:%d",
				player.PlayerID, player.ChestPoints, player.GoalBonus, player.Score,
				expected.chest, expected.goal, expected.score)
		}
	}
	placements := restored.Placements()
	if placements[0].PlayerID != "p3" || placements[0].Place != 1 {
		t.Fatalf("placements[0] = %+v, want p3 first", placements[0])
	}
	if placements[1].PlayerID != "p1" || placements[1].Place != 2 {
		t.Fatalf("placements[1] = %+v, want p1 second", placements[1])
	}
	if placements[2].PlayerID != "p2" || placements[2].Place != 3 {
		t.Fatalf("placements[2] = %+v, want p2 third", placements[2])
	}
}

func TestCurrentVisibleStateClampsTerminalTurn(t *testing.T) {
	match, err := New(Config{
		GameVersion: GameVersion,
		Ruleset:     RulesetSeededMazeV1,
		PlayerIDs:   []string{"p1", "p2"},
		RNGSeed:     testSeedAlpha,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for !match.Terminal() {
		if err := match.Apply(map[string]Action{
			"p1": {Action: "wait"},
			"p2": {Action: "wait"},
		}); err != nil {
			t.Fatalf("Apply: %v", err)
		}
	}
	visible, err := match.CurrentVisibleState("p1")
	if err != nil {
		t.Fatalf("CurrentVisibleState: %v", err)
	}
	if visible.Turn != match.Turn() {
		t.Fatalf("visible turn = %d, want %d", visible.Turn, match.Turn())
	}
	if visible.RemainingTurns != 0 {
		t.Fatalf("remaining turns = %d, want 0", visible.RemainingTurns)
	}
}

func assertWallBoundary(t *testing.T, ruleset string, tiles []string) {
	t.Helper()
	if len(tiles) == 0 {
		t.Fatalf("%s produced empty tile rows", ruleset)
	}
	if len(tiles[0]) == 0 {
		t.Fatalf("%s produced an empty first row", ruleset)
	}
	lastRow := len(tiles) - 1
	lastCol := len(tiles[0]) - 1
	for y, row := range tiles {
		if len(row) != lastCol+1 {
			t.Fatalf("%s produced jagged row %d: got width %d, want %d", ruleset, y, len(row), lastCol+1)
		}
	}
	for x := 0; x <= lastCol; x++ {
		if tiles[0][x] != '#' || tiles[lastRow][x] != '#' {
			t.Fatalf("%s has non-wall boundary row", ruleset)
		}
	}
	for y := 0; y <= lastRow; y++ {
		if tiles[y][0] != '#' || tiles[y][lastCol] != '#' {
			t.Fatalf("%s has non-wall boundary col", ruleset)
		}
	}
}

func assertReachable(t *testing.T, ruleset string, match *Match, from, to Position) {
	t.Helper()
	if path, ok := match.ShortestPath(from, to); !ok || len(path) == 0 {
		t.Fatalf("%s no path from %+v to %+v", ruleset, from, to)
	}
}

func assertDistinctPlacements(t *testing.T, ruleset string, layout GeneratedLayout) {
	t.Helper()
	seen := map[string]string{}
	mark := func(kind string, pos Position) {
		key := posKey(pos)
		if prev, ok := seen[key]; ok {
			t.Fatalf("%s %s overlaps with %s at %+v", ruleset, kind, prev, pos)
		}
		seen[key] = kind
	}
	mark("goal", layout.Goal)
	for i, spawn := range layout.SpawnPoints {
		mark(fmt.Sprintf("spawn_%d", i+1), spawn)
	}
	for i, chest := range layout.InitialChests {
		mark(fmt.Sprintf("chest_%d", i+1), Position{X: chest.X, Y: chest.Y})
	}
}
