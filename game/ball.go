package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BallRadius      float32 = 10
	BallRadiusLight float32 = 5
	BallSpeed               = 200
)

type Ball struct {
	Position   rl.Vector2
	Velocity   rl.Vector2
	game       *Game
	min        rl.Vector2
	max        rl.Vector2
	color      rl.Color
	colorLight rl.Color
}

func (b *Ball) Init(g *Game) {
	b.game = g
	b.color = rl.Red
	b.colorLight = rl.ColorBrightness(b.color, 0.3)
	b.setVelocity(b.Velocity)
	b.min = rl.Vector2{X: BallRadius, Y: BallRadius}
	b.max = rl.Vector2{X: float32(b.game.Width) - BallRadius, Y: float32(b.game.Height) - BallRadius}
}

func (b *Ball) setVelocity(v rl.Vector2) {
	b.Velocity = scale(normalize(v), BallSpeed)
}

func (b *Ball) Update(time float32) rl.Vector2 {
	newPosition := add(b.Position, scale(b.Velocity, time))
	b.Position = crop(newPosition, b.min, b.max)
	return newPosition
}

func (b *Ball) DrawShadow() {
	shadowCenter := b.game.ProjectCanvas(b.Position)
	shadowRadius := b.game.ProjectCanvas(rl.Vector2{X: b.Position.X + BallRadius, Y: b.Position.Y}).X - shadowCenter.X
	rl.DrawCircleV(shadowCenter, shadowRadius, b.game.shadowColorLight)
	rl.DrawCircleV(shadowCenter, shadowRadius-shadowBorderThickness, b.game.shadowColor)
}

func (b *Ball) Draw() {
	rl.DrawCircleV(b.Position, BallRadius, b.color)
	rl.DrawCircleV(rl.Vector2{X: b.Position.X - 1, Y: b.Position.Y - 2}, BallRadiusLight, b.colorLight)
}
