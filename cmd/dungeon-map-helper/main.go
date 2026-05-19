// Command dungeon-map-helper prints fixed-map layout and shortest-path helpers.
//
// Keep the helper on public dungeon APIs and avoid ai-arena internal
// dependencies.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/yoskeoka/dungeon-game-ai-arena/games/dungeon"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	return runWithIO(args, os.Stdout, os.Stderr)
}

func runWithIO(args []string, stdout, stderr io.Writer) error {
	var (
		rulesets rulesetListFlag
		rngSeed  = dungeon.DefaultRNGSeed
	)
	fs := flag.NewFlagSet("dungeon-map-helper", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Var(&rulesets, "ruleset", "dungeon ruleset (repeatable or comma-separated)")
	fs.StringVar(&rngSeed, "rng-seed", rngSeed, "deterministic seed")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(rulesets) == 0 {
		rulesets = append(rulesets, dungeon.RulesetSeededMazeV1)
	}

	for i, ruleset := range rulesets {
		if i > 0 {
			fmt.Fprintln(stdout)
		}
		if err := printRulesetSummary(stdout, rngSeed, ruleset); err != nil {
			return err
		}
	}
	return nil
}

type rulesetListFlag []string

func (f *rulesetListFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *rulesetListFlag) Set(value string) error {
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		*f = append(*f, part)
	}
	return nil
}

func printRulesetSummary(stdout io.Writer, rngSeed, ruleset string) error {
	world, err := dungeon.New(dungeon.Config{
		GameVersion: dungeon.GameVersion,
		Ruleset:     ruleset,
		PlayerIDs:   []string{"p1", "p2"},
		RNGSeed:     rngSeed,
	})
	if err != nil {
		return err
	}
	layout := world.Layout()
	rulesetDef := world.Ruleset()
	height := len(layout.Tiles)
	width := 0
	walkable := 0
	if height > 0 {
		width = len(layout.Tiles[0])
	}
	for _, row := range layout.Tiles {
		for x := 0; x < len(row); x++ {
			if row[x] != '#' {
				walkable++
			}
		}
	}

	fmt.Fprintf(stdout, "=== ruleset=%s map_id=%s rng_seed=%s ===\n",
		ruleset,
		rulesetDef.MapID,
		world.PublicState().RNGSeed,
	)
	fmt.Fprintf(stdout, "shape width=%d height=%d walkable=%d max_turns=%d view_radius=%d goal_bonuses=%v\n",
		width,
		height,
		walkable,
		rulesetDef.MaxTurns,
		rulesetDef.ViewRadius,
		rulesetDef.GoalBonuses,
	)
	fmt.Fprintf(stdout, "placements goal=%+v spawns=%v chests=%v\n", layout.Goal, layout.SpawnPoints, layout.InitialChests)

	chestTotal := 0
	for _, chest := range layout.InitialChests {
		chestTotal += chest.Points
	}
	majorityFloor := chestTotal/2 + 1
	thirdPlaceBonus := 0
	if len(rulesetDef.GoalBonuses) >= 3 {
		thirdPlaceBonus = rulesetDef.GoalBonuses[2]
	}
	chestValues := make([]int, 0, len(layout.InitialChests))
	for _, chest := range layout.InitialChests {
		chestValues = append(chestValues, chest.Points)
	}
	slices.Sort(chestValues)
	majorityChestMin := chestTotal
	for mask := 1; mask < 1<<len(chestValues); mask++ {
		sum := 0
		for i, points := range chestValues {
			if mask&(1<<i) != 0 {
				sum += points
			}
		}
		if sum >= majorityFloor && sum < majorityChestMin {
			majorityChestMin = sum
		}
	}
	fmt.Fprintf(stdout, "balance chest_total=%d majority_threshold=%d first_no_chest=%d third_with_min_majority=%d\n",
		chestTotal,
		majorityFloor,
		rulesetDef.GoalBonuses[0],
		thirdPlaceBonus+majorityChestMin,
	)
	for _, row := range world.PublicState().Tiles {
		fmt.Fprintln(stdout, row)
	}
	for i, spawn := range world.SpawnPoints()[:2] {
		path, ok := world.ShortestPath(spawn, layout.Goal)
		if !ok {
			return fmt.Errorf("no path from spawn %d to goal", i+1)
		}
		fmt.Fprintf(stdout, "spawn_%d_to_goal steps=%d route=%v\n", i+1, len(path)-1, path)
	}
	for i, chest := range layout.InitialChests {
		path, ok := world.ShortestPath(world.SpawnPoints()[0], dungeon.Position{X: chest.X, Y: chest.Y})
		if !ok {
			return fmt.Errorf("no path from spawn 1 to chest %d", i+1)
		}
		fmt.Fprintf(stdout, "spawn_1_to_chest_%d points=%d steps=%d route=%v\n", i+1, chest.Points, len(path)-1, path)
	}
	return nil
}
