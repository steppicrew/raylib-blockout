package main

import (
	"breakout/game"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	screenWidth := int32(800)
	screenHeight := int32(600)

	rl.InitWindow(screenWidth, screenHeight, "Breakout in Go")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	board := game.Game{
		Width:  screenWidth,
		Height: screenHeight,
	}

	board.Init()

	for !rl.WindowShouldClose() {
		board.Update(rl.GetFrameTime())

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		board.Draw()

		rl.EndDrawing()
	}
}
