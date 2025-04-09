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

	game := game.Game{
		Width:  screenWidth,
		Height: screenHeight,
	}

	game.Init()

	for !rl.WindowShouldClose() {
		game.Update(rl.GetFrameTime())

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		game.DrawShadow()
		game.Draw()

		rl.EndDrawing()
	}
}
