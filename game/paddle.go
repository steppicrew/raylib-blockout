package game

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	PaddleWidth  = 80
	PaddleHeight = 20
	PaddleSpeed  = 200
)

type Paddle struct {
	Position rl.Vector2
	game     *Game
}

func (p *Paddle) Init(g *Game) {
	p.game = g
}

func (p *Paddle) Update(time float32) {
	if rl.IsKeyDown(rl.KeyRight) {
		p.Position.X = MinFloat(p.Position.X+PaddleSpeed*time, float32(p.game.Width)-PaddleWidth)
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		p.Position.X = MaxFloat(p.Position.X-PaddleSpeed*time, 0)
	}
}

func (p *Paddle) checkCollision(b Ball) rl.Vector2 {
	topLeft := rl.Vector2{X: p.Position.X, Y: p.Position.Y}
	topRight := rl.Vector2{X: p.Position.X + PaddleWidth, Y: p.Position.Y}
	bottomLeft := rl.Vector2{X: p.Position.X, Y: p.Position.Y + PaddleHeight}
	bottomRight := rl.Vector2{X: p.Position.X + PaddleWidth, Y: p.Position.Y + PaddleHeight}
	resultVelocity := rl.Vector2{X: b.Velocity.X, Y: b.Velocity.Y}
	if rl.CheckCollisionCircleLine(b.Position, float32(b.radius), topLeft, topRight) {
		resultVelocity.Y = -resultVelocity.Y
	}
	if rl.CheckCollisionCircleLine(b.Position, float32(b.radius), topRight, bottomRight) {
		resultVelocity.X = -resultVelocity.X
	}
	if rl.CheckCollisionCircleLine(b.Position, float32(b.radius), bottomLeft, topLeft) {
		resultVelocity.X = -resultVelocity.X
	}
	return resultVelocity
}

func (p *Paddle) Draw() {
	rl.DrawRectangle(int32(p.Position.X), int32(p.Position.Y), int32(PaddleWidth), int32(PaddleHeight), rl.Black)
}
