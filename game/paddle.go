package game

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	PaddleWidth                     = 8
	PaddleHeight                    = 3
	PaddleDepth                     = 3
	PaddleSpeed                     = 20
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
}

func (p *Paddle) Init(g *Game) {
	p.game = g
	p.color = rl.Maroon

	p.scale, p.baseBbox = GetScaleBoundingBox(g.brickModel, rl.NewVector3(0, 0, 0), rl.Vector3{X: PaddleWidth, Y: PaddleHeight, Z: PaddleDepth})
}

func (p *Paddle) updateX(x float32) {
	p.Position.X = x
	p.bbox = rl.BoundingBox{
		Min: rl.Vector3Add(p.Position, p.baseBbox.Min),
		Max: rl.Vector3Add(p.Position, p.baseBbox.Max),
	}
}

func (p *Paddle) Update(time float32) {
	if !rl.IsKeyDown(rl.KeyLeftShift) && !rl.IsKeyDown(rl.KeyRightShift) && !rl.IsKeyDown(rl.KeyLeftControl) && !rl.IsKeyDown(rl.KeyRightControl) && !rl.IsKeyDown(rl.KeyLeftAlt) && !rl.IsKeyDown(rl.KeyRightAlt) {
		if rl.IsKeyDown(rl.KeyLeft) {
			p.updateX(MinFloat(p.Position.X+PaddleSpeed*time, float32(p.game.Width)-PaddleWidth))
		}
		if rl.IsKeyDown(rl.KeyRight) {
			p.updateX(MaxFloat(p.Position.X-PaddleSpeed*time, 0))
		}
	}
}

func (p *Paddle) checkCollision(b Ball) rl.Vector3 {

	collides, closest := CheckCollisionSphereBox(b.Position, BallRadius, p.bbox)

	resultVelocity := rl.Vector3{X: b.Velocity.X, Y: b.Velocity.Y, Z: b.Velocity.Z}

	if collides {
		if closest.X == p.bbox.Min.X || closest.X == p.bbox.Max.X {
			resultVelocity.X = -resultVelocity.X
		}
		if closest.Z == p.bbox.Min.Z || closest.Z == p.bbox.Max.Z {
			resultVelocity.Z = -resultVelocity.Z
		}
	}
	return resultVelocity
}

func (p *Paddle) DrawShadow() {
}

func (p *Paddle) Draw() {
	shader := p.game.shader
	modelLoc := rl.GetShaderLocation(shader, "model")
	objectColorLoc := rl.GetShaderLocation(shader, "objectColor")

	rl.BeginShaderMode(shader)

	// Update uniforms
	rl.SetShaderValue(shader, objectColorLoc, []float32{float32(p.color.R) / 255, float32(p.color.G) / 255, float32(p.color.B) / 255}, rl.ShaderUniformVec3)

	transform := rl.MatrixTranslate(p.Position.X, p.Position.Y, p.Position.Z)
	// Or build full transform (translation + rotation + scale)

	rl.SetShaderValueMatrix(shader, modelLoc, transform)

	rl.DrawModelEx(p.game.brickModel, p.Position, rl.Vector3{X: 0, Y: 1, Z: 0}, 0, p.scale, p.color)

	rl.EndShaderMode()
}
