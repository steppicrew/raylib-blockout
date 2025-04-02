package game

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	BrickRows    = 5
	HeightOffset = 100
)

type Game struct {
	Width  int32
	Height int32
	ball   Ball
	paddle Paddle
	bricks []*Brick
}

func (g *Game) Init() {
	g.ball = Ball{
		Position: rl.Vector2{X: float32(g.Width) / 2, Y: float32(g.Height) * 0.8},
		Velocity: rl.Vector2{X: 1, Y: -1},
	}
	g.paddle = Paddle{Position: rl.Vector2{X: float32(g.Width) / 2, Y: float32(g.Height) - PaddleHeight - 30}}
	g.bricks = []*Brick{}

	g.ball.Init(g)
	g.paddle.Init(g)

	cols := int(g.Width / BrickWidth)

	g.bricks = make([]*Brick, BrickRows*cols)
	for y := 0; y < BrickRows; y++ {
		for x := 0; x < cols; x++ {
			brick := Brick{
				Position: rl.Vector2{X: float32(x) * BrickWidth, Y: float32(y)*BrickHeight + HeightOffset},
				Lives:    BrickRows - y,
			}
			brick.Init()
			g.bricks[y*cols+x] = &brick
		}
	}
}

func (g *Game) Update(time float32) {
	ballPosition := g.ball.Update(time)
	if ballPosition.X < g.ball.min.X || ballPosition.X > g.ball.max.X {
		g.ball.Velocity.X = -g.ball.Velocity.X
	}
	if ballPosition.Y < g.ball.min.Y || ballPosition.Y > g.ball.max.Y {
		g.ball.Velocity.Y = -g.ball.Velocity.Y
	}

	g.paddle.Update(time)
	g.ball.setVelocity(g.paddle.checkCollision(g.ball))

	for _, brick := range g.bricks {
		if brick.Lives <= 0 {
			continue
		}
		g.ball.setVelocity(brick.checkCollision(g.ball))

		brick.Update(time)
	}

}

func (g *Game) Draw() {
	g.ball.Draw()
	g.paddle.Draw()
	for _, brick := range g.bricks {
		brick.Draw()
	}

}
