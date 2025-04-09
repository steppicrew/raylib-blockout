package game

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	PaddleWidth                     = 80
	PaddleHeight                    = 20
	PaddleSpeed                     = 200
	PaddleRoundness         float32 = .4
	PaddleRoundnessSegments int32   = 10
)

type Paddle struct {
	Position rl.Vector2
	game     *Game
	color    rl.Color
}

func (p *Paddle) Init(g *Game) {
	p.game = g
	p.color = rl.Maroon
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
	if rl.CheckCollisionCircleLine(b.Position, BallRadius, topLeft, topRight) {
		resultVelocity.Y = -resultVelocity.Y
	}
	if rl.CheckCollisionCircleLine(b.Position, BallRadius, topRight, bottomRight) {
		resultVelocity.X = -resultVelocity.X
	}
	if rl.CheckCollisionCircleLine(b.Position, BallRadius, bottomLeft, topLeft) {
		resultVelocity.X = -resultVelocity.X
	}
	return resultVelocity
}

func (p *Paddle) DrawShadow() {
	shadowTopLeft := p.game.ProjectCanvas(p.Position)
	shadowBottomRight := p.game.ProjectCanvas(rl.Vector2Add(p.Position, rl.Vector2{X: PaddleWidth, Y: PaddleHeight}))
	rectangle := rl.Rectangle{X: shadowTopLeft.X, Y: shadowTopLeft.Y, Width: shadowBottomRight.X - shadowTopLeft.X, Height: shadowBottomRight.Y - shadowTopLeft.Y}
	rl.DrawRectangleRounded(rectangle, PaddleRoundness, segments, p.game.shadowColorLight)
	rectangle = rl.Rectangle{X: shadowTopLeft.X + shadowBorderThickness, Y: shadowTopLeft.Y + shadowBorderThickness, Width: shadowBottomRight.X - shadowTopLeft.X - 2*shadowBorderThickness, Height: shadowBottomRight.Y - shadowTopLeft.Y - 2*shadowBorderThickness}
	rl.DrawRectangleRounded(rectangle, PaddleRoundness, segments, p.game.shadowColor)
}

func (p *Paddle) Draw() {
	rectangle := rl.Rectangle{X: p.Position.X, Y: p.Position.Y, Width: PaddleWidth, Height: PaddleHeight}
	rl.DrawRectangleRounded(rectangle, PaddleRoundness, PaddleRoundnessSegments, p.color)
}
