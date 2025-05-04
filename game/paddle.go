package game

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	PaddleWidth             float32 = 8
	PaddleHeight            float32 = 3
	PaddleDepth             float32 = 3
	PaddleSpeed             float32 = 20
	PaddleRoundness         float32 = .4
	PaddleRoundnessSegments int32   = 15
)

type Paddle struct {
	Position rl.Vector3
	game     *Game
	color    rl.Color
	scale    rl.Vector3
	baseBbox rl.BoundingBox
	bbox     rl.BoundingBox

	transformation rl.Matrix
}

func (p *Paddle) Init(g *Game) {
	p.game = g
	p.color = rl.Maroon

	p.scale, p.baseBbox = GetScaleBoundingBox(g.brickModel, rl.NewVector3(0, 0, 0), rl.Vector3{X: PaddleWidth, Y: PaddleHeight, Z: PaddleDepth})
	p.updateX(float32(g.Width) / 2)
}

func (p *Paddle) updateX(x float32) {
	p.Position.X = x
	p.bbox = rl.BoundingBox{
		Min: rl.Vector3Add(p.Position, p.baseBbox.Min),
		Max: rl.Vector3Add(p.Position, p.baseBbox.Max),
	}

	scaleMatrix := rl.MatrixScale(p.scale.X, p.scale.Y, p.scale.Z)
	rotationMatrix := rl.MatrixRotateY(0)
	translationMatrix := rl.MatrixTranslate(p.Position.X, p.Position.Y, p.Position.Z)
	p.transformation = rl.MatrixMultiply(rl.MatrixMultiply(scaleMatrix, rotationMatrix), translationMatrix)
}

func (p *Paddle) Update(time float32) {
	if !rl.IsKeyDown(rl.KeyLeftShift) && !rl.IsKeyDown(rl.KeyRightShift) && !rl.IsKeyDown(rl.KeyLeftControl) && !rl.IsKeyDown(rl.KeyRightControl) && !rl.IsKeyDown(rl.KeyLeftAlt) && !rl.IsKeyDown(rl.KeyRightAlt) {
		if rl.IsKeyDown(rl.KeyLeft) {
			p.updateX(MinFloat(p.Position.X+PaddleSpeed*time, float32(p.game.Width)-PaddleWidth/2))
		}
		if rl.IsKeyDown(rl.KeyRight) {
			p.updateX(MaxFloat(p.Position.X-PaddleSpeed*time, PaddleWidth/2))
		}
	}
}

func (p *Paddle) checkCollision(b Ball) rl.Vector3 {

	collides, closest := CheckCollisionSphereBox(b.Position, BallRadius, p.bbox)

	resultVelocity := rl.Vector3{X: b.Velocity.X, Y: b.Velocity.Y, Z: b.Velocity.Z}

	if collides {
		if closest.X == p.bbox.Max.X {
			resultVelocity.X = -resultVelocity.X
		}
		if closest.Z == p.bbox.Min.Z {
			resultVelocity.Z = -AbsFloat(resultVelocity.Z)
		}
		if closest.Z == p.bbox.Max.Z {
			resultVelocity.Z = AbsFloat(resultVelocity.Z)
		}
	}
	return resultVelocity
}

func (p *Paddle) DrawShadow(shader rl.Shader, modelLoc int32) {
	rl.BeginShaderMode(shader)

	rl.SetShaderValueMatrix(shader, modelLoc, p.transformation)

	// rl.DrawModelEx(b.game.ballModel, b.Position, rl.Vector3{X: 0, Y: 1, Z: 0}, 0, rl.Vector3{X: 1, Y: 1, Z: 1}, b.color)
	rl.DrawModel(p.game.brickShadowModel, rl.Vector3Zero(), 1, p.game.shadowColor)

	rl.EndShaderMode()
}

func (p *Paddle) Draw(shader rl.Shader, modelLoc int32) {
	rl.BeginShaderMode(shader)

	// Update uniforms
	SetObjectColor(shader, p.color)

	rl.SetShaderValueMatrix(shader, modelLoc, p.transformation)

	// rl.DrawModelEx(p.game.brickModel, p.Position, rl.Vector3{X: 0, Y: 1, Z: 0}, 0, p.scale, p.color)
	rl.DrawModel(p.game.brickModel, rl.Vector3Zero(), 1, p.color)

	rl.EndShaderMode()
}
