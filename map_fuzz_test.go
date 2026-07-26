package main

import (
	"math/rand"
	"testing"
)

func TestGenerateMapReachabilityFuzz(t *testing.T) {
	mapSizes := [][2]int{
		{20, 10},
		{60, 30},
		{80, 40},
	}

	for _, size := range mapSizes {
		w, h := size[0], size[1]
		for seed := int64(1); seed <= 500; seed++ {
			rng := rand.New(rand.NewSource(seed))
			gameMap, rooms := GenerateMap(w, h, rng)

			if len(rooms) == 0 {
				t.Fatalf("seed %d (size %dx%d): generated 0 rooms", seed, w, h)
			}

			startPos := rooms[0].Center()

			// BFS to find all reachable non-wall tiles from startPos
			visited := make([][]bool, h)
			for y := range visited {
				visited[y] = make([]bool, w)
			}

			queue := []Position{startPos}
			visited[startPos.Y][startPos.X] = true
			reachableCount := 0

			dirs := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

			for len(queue) > 0 {
				curr := queue[0]
				queue = queue[1:]
				reachableCount++

				for _, d := range dirs {
					nx, ny := curr.X+d[0], curr.Y+d[1]
					if nx >= 0 && nx < w && ny >= 0 && ny < h {
						if !visited[ny][nx] && gameMap.Tiles[ny][nx] != TileWall {
							visited[ny][nx] = true
							queue = append(queue, Position{X: nx, Y: ny})
						}
					}
				}
			}

			// Verify every room center is reachable
			for rIdx, r := range rooms {
				center := r.Center()
				if !visited[center.Y][center.X] {
					t.Fatalf("seed %d (size %dx%d): room %d center %+v is unreachable from spawn %+v", seed, w, h, rIdx, center, startPos)
				}
			}

			// Verify stairs tile is reachable
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					if gameMap.Tiles[y][x] == TileStairsDown {
						if !visited[y][x] {
							t.Fatalf("seed %d (size %dx%d): stairs at (%d,%d) are unreachable from spawn %+v", seed, w, h, x, y, startPos)
						}
					}
				}
			}
		}
	}
}
