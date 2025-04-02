package game

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	BrickWidth  = 80
	BrickHeight = 30
)

var BrickColors = []rl.Color{rl.Green, rl.Yellow, rl.Orange, rl.Purple, rl.Red, rl.Blue}

type Brick struct {
	Position  rl.Vector2
	Lives     int
	rectangle rl.Rectangle
}

func (b *Brick) Init() {
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
	if rl.CheckCollisionCircleLine(b.Position, float32(b.radius), topLeft, topRight) {
		resultVelocity.Y = -resultVelocity.Y
		bounced = true
	}
	if rl.CheckCollisionCircleLine(b.Position, float32(b.radius), bottomLeft, bottomRight) {
		resultVelocity.Y = -resultVelocity.Y
		bounced = true
	}
	if rl.CheckCollisionCircleLine(b.Position, float32(b.radius), topRight, bottomRight) {
		resultVelocity.X = -resultVelocity.X
		bounced = true
	}
	if rl.CheckCollisionCircleLine(b.Position, float32(b.radius), bottomLeft, topLeft) {
		resultVelocity.X = -resultVelocity.X
		bounced = true
	}
	if bounced && brick.Lives > 0 {
		brick.Lives--
	}
	return resultVelocity
}

func (b *Brick) Draw() {
	if b.Lives > 0 {
		rl.DrawRectangle(int32(b.Position.X), int32(b.Position.Y), int32(BrickWidth), int32(BrickHeight), BrickColors[b.Lives-1])
		rl.DrawRectangleLines(int32(b.Position.X), int32(b.Position.Y), int32(BrickWidth), int32(BrickHeight), rl.Black)
	}
}
