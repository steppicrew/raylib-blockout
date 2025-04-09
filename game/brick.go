package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BrickWidth  float32 = 80
	BrickHeight float32 = 30

	segments  int32   = 0
	roundness float32 = .7
	// shadowThickness float32 = 2
	padding float32 = 2
)

var BrickColors = []rl.Color{rl.Green, rl.Yellow, rl.Orange, rl.Purple, rl.Red, rl.Blue}

type Brick struct {
	Position rl.Vector2
	Lives    int
	game     *Game

	rectangle   rl.Rectangle
	shadowColor rl.Color
}

func (b *Brick) Init(g *Game) {
	b.game = g
	b.rectangle = rl.Rectangle{
		X:      b.Position.X,
		Y:      b.Position.Y,
		Width:  PaddleWidth,
		Height: PaddleHeight,
	}
}

func (b *Brick) Update(time float32) {
}

func (brick *Brick) checkCollision(b Ball) rl.Vector2 {
	topLeft := rl.Vector2{X: brick.Position.X, Y: brick.Position.Y}
	topRight := rl.Vector2{X: brick.Position.X + PaddleWidth, Y: brick.Position.Y}
	bottomLeft := rl.Vector2{X: brick.Position.X, Y: brick.Position.Y + BrickHeight}
	bottomRight := rl.Vector2{X: brick.Position.X + BrickWidth, Y: brick.Position.Y + BrickHeight}
	resultVelocity := rl.Vector2{X: b.Velocity.X, Y: b.Velocity.Y}
	bounced := false
	if rl.CheckCollisionCircleLine(b.Position, BallRadius, topLeft, topRight) {
		resultVelocity.Y = -resultVelocity.Y
		bounced = true
	}
	if rl.CheckCollisionCircleLine(b.Position, BallRadius, bottomLeft, bottomRight) {
		resultVelocity.Y = -resultVelocity.Y
		bounced = true
	}
	if rl.CheckCollisionCircleLine(b.Position, BallRadius, topRight, bottomRight) {
		resultVelocity.X = -resultVelocity.X
		bounced = true
	}
	if rl.CheckCollisionCircleLine(b.Position, BallRadius, bottomLeft, topLeft) {
		resultVelocity.X = -resultVelocity.X
		bounced = true
	}
	if bounced && brick.Lives > 0 {
		brick.Lives--
	}
	return resultVelocity
}

func (b *Brick) DrawShadow() {
	if b.Lives > 0 {
		x := b.Position.X + padding/2
		y := b.Position.Y + padding/2

		width := BrickWidth - padding
		height := BrickHeight - padding

		var rectangle rl.Rectangle

		shadowTopLeft := b.game.ProjectCanvas(rl.NewVector2(x, y))
		shadowBottomRight := b.game.ProjectCanvas(rl.NewVector2(x+width, y+height))
		rectangle = rl.Rectangle{X: shadowTopLeft.X, Y: shadowTopLeft.Y, Width: shadowBottomRight.X - shadowTopLeft.X, Height: shadowBottomRight.Y - shadowTopLeft.Y}
		rl.DrawRectangleRounded(rectangle, roundness, segments, b.game.shadowColorLight)
		rectangle = rl.Rectangle{X: shadowTopLeft.X + shadowBorderThickness, Y: shadowTopLeft.Y + shadowBorderThickness, Width: shadowBottomRight.X - shadowTopLeft.X - 2*shadowBorderThickness, Height: shadowBottomRight.Y - shadowTopLeft.Y - 2*shadowBorderThickness}
		rl.DrawRectangleRounded(rectangle, roundness, segments, b.game.shadowColor)
	}
}

func (b *Brick) Draw() {
	if b.Lives > 0 {
		x := b.Position.X + padding/2
		y := b.Position.Y + padding/2

		width := BrickWidth - padding
		height := BrickHeight - padding

		var rectangle rl.Rectangle

		baseColor := BrickColors[b.Lives-1]
		highlightColor := rl.ColorBrightness(baseColor, 0.3)

		rectangle = rl.Rectangle{X: x, Y: y, Width: width, Height: height}
		rl.DrawRectangleRounded(rectangle, roundness, segments, baseColor)

		var highlightOffsetX float32 = .15 * width
		var highlightWidth float32 = .6 * width
		var highlightOffsetY float32 = .1 * height
		var highlightHeight float32 = .3 * height
		rectangle = rl.Rectangle{X: x + highlightOffsetX, Y: y + highlightOffsetY, Width: highlightWidth, Height: highlightHeight}
		rl.DrawRectangleRounded(rectangle, roundness, segments, highlightColor)

	}
}
