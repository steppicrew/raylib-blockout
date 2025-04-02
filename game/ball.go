package game

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	BallRadius = 10
	BallSpeed  = 200
)

type Ball struct {
	Position rl.Vector2
	Velocity rl.Vector2
	game     *Game
	radius   int32
	min      rl.Vector2
	max      rl.Vector2
}

func (b *Ball) Init(g *Game) {
	b.game = g
	b.radius = BallRadius
	b.setVelocity(b.Velocity)
	b.min = rl.Vector2{X: float32(b.radius), Y: float32(b.radius)}
	b.max = rl.Vector2{X: float32(b.game.Width) - float32(b.radius), Y: float32(b.game.Height) - float32(b.radius)}
}

func (b *Ball) setVelocity(v rl.Vector2) {
	b.Velocity = scale(normalize(v), BallSpeed)
}

func (b *Ball) Update(time float32) rl.Vector2 {
	newPosition := add(b.Position, scale(b.Velocity, time))
	b.Position = crop(newPosition, b.min, b.max)
	return newPosition
}

func (b *Ball) Draw() {
	rl.DrawCircle(int32(b.Position.X), int32(b.Position.Y), float32(b.radius), rl.Black)
}
