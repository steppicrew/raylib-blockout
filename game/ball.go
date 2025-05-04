package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BallRadius      float32 = 1.5
	BallRadiusLight float32 = 5
	BallSpeed               = 20
)

type Ball struct {
	Position   rl.Vector3
	Velocity   rl.Vector3
	game       *Game
	min        rl.Vector3
	max        rl.Vector3
	color      rl.Color
	colorLight rl.Color

	transformation rl.Matrix
}

func (b *Ball) Init(g *Game) {
	b.game = g

	b.color = rl.Red
	b.colorLight = rl.ColorBrightness(b.color, 0.5)
	b.setVelocity(b.Velocity)
	b.min = rl.Vector3{X: BallRadius, Y: GamePlaneHeight, Z: BallRadius}
	b.max = rl.Vector3{X: float32(b.game.Width) - BallRadius, Y: GamePlaneHeight, Z: float32(b.game.Height) - BallRadius}
	b.min = rl.Vector3{X: 0, Y: 0, Z: 0}
	b.max = rl.Vector3{X: float32(b.game.Width) - BallRadius, Y: GamePlaneHeight, Z: float32(b.game.Height) - BallRadius}
}

func (b *Ball) setVelocity(v rl.Vector3) {
	b.Velocity = rl.Vector3Scale(rl.Vector3Normalize(v), BallSpeed)
	// b.Velocity = rl.Vector3{X: 0, Y: 0, Z: 0}
}

func (b *Ball) Update(time float32) rl.Vector3 {
	//return b.Position
	newPosition := rl.Vector3Add(b.Position, rl.Vector3Scale(b.Velocity, time))
	b.Position = rl.Vector3Clamp(newPosition, b.min, b.max)

	b.transformation = rl.MatrixTranslate(b.Position.X, b.Position.Y, b.Position.Z)

	// b.Position = newPosition
	return newPosition
}

func (b *Ball) DrawShadow(shader rl.Shader, modelLoc int32) {
	rl.BeginShaderMode(shader)

	rl.SetShaderValueMatrix(shader, modelLoc, b.transformation)

	// rl.DrawModelEx(b.game.ballModel, b.Position, rl.Vector3{X: 0, Y: 1, Z: 0}, 0, rl.Vector3{X: 1, Y: 1, Z: 1}, b.color)
	rl.DrawModel(b.game.ballShadowModel, rl.Vector3Zero(), 1, b.game.shadowColor)

	rl.EndShaderMode()
}

func (b *Ball) Draw(shader rl.Shader, modelLoc int32) {
	rl.BeginShaderMode(shader)

	// Update uniforms
	SetObjectColor(shader, b.color)

	rl.SetShaderValueMatrix(shader, modelLoc, b.transformation)

	// rl.DrawModelEx(b.game.ballModel, b.Position, rl.Vector3{X: 0, Y: 1, Z: 0}, 0, rl.Vector3{X: 1, Y: 1, Z: 1}, b.color)
	rl.DrawModel(b.game.ballModel, rl.Vector3Zero(), 1, b.color)

	// DrawNormals(b.game.ballModel.GetMeshes()[0], b.Position)

	rl.EndShaderMode()
}
