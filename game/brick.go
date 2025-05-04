package game

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BrickWidth  float32 = 6
	BrickHeight float32 = 3
	BrickDepth  float32 = 3
	segments    int32   = 0
	roundness   float32 = .7
	// shadowThickness float32 = 2
	padding float32 = 0.2

	MaxWiggleTime float32 = 2
)

func ease(x float64) float64 {
	return 1 - math.Sqrt(1-math.Pow(x-1, 4))
}

func wiggle(time float32) float32 {
	t := float64(time)
	wiggle := math.Sin(t*30) * ease(t/float64(MaxWiggleTime))
	return float32(wiggle)
}

var BrickColors = []rl.Color{rl.Green, rl.Yellow, rl.Orange, rl.Purple, rl.Red, rl.Blue}

type Brick struct {
	Position   rl.Vector3
	angle      float32
	Lives      int
	reallyDead bool
	game       *Game
	size       rl.Vector3
	scale      rl.Vector3
	bbox       rl.BoundingBox

	wiggleTime      float32
	wiggleDirection float32
	transformation  rl.Matrix
}

func (b *Brick) Init(g *Game) {
	b.game = g
	b.angle = 0
	b.wiggleTime = MaxWiggleTime
	b.reallyDead = false

	bbox := rl.GetModelBoundingBox(g.brickModel)

	objWidth := bbox.Max.X - bbox.Min.X
	objHeight := bbox.Max.Y - bbox.Min.Y
	objDepth := bbox.Max.Z - bbox.Min.Z

	b.scale = rl.NewVector3((BrickWidth-padding)/objWidth, BrickHeight/objHeight, BrickDepth/objDepth)
	half := rl.Vector3{X: (BrickWidth - padding) / 2, Y: BrickHeight / 2, Z: BrickDepth / 2}
	b.bbox = rl.BoundingBox{
		Min: rl.Vector3Subtract(b.Position, half),
		Max: rl.Vector3Add(b.Position, half),
	}

	b.size = rl.Vector3{X: BrickWidth, Y: BrickHeight, Z: BrickDepth}
}

func (b *Brick) Update(time float32) {
	if b.wiggleTime < MaxWiggleTime {
		b.wiggleTime += time
		b.angle = wiggle(b.wiggleTime) * .2 * b.wiggleDirection
	} else {
		b.reallyDead = b.Lives <= 0
	}

	scaleMatrix := rl.MatrixScale(b.scale.X, b.scale.Y, b.scale.Z)
	rotationMatrix := rl.MatrixRotateY(b.angle)
	translationMatrix := rl.MatrixTranslate(b.Position.X, b.Position.Y, b.Position.Z)
	b.transformation = rl.MatrixMultiply(rl.MatrixMultiply(scaleMatrix, rotationMatrix), translationMatrix)
}

func (b *Brick) checkCollision(ball Ball) rl.Vector3 {
	if b.Lives <= 0 {
		return ball.Velocity
	}

	collides, closest := CheckCollisionSphereBox(ball.Position, BallRadius, b.bbox)

	resultVelocity := rl.Vector3{X: ball.Velocity.X, Y: ball.Velocity.Y, Z: ball.Velocity.Z}

	if collides {
		b.wiggleTime = 0

		closestVector := rl.Vector3Subtract(b.Position, closest)
		if closestVector.X*closestVector.Z > 0 {
			b.wiggleDirection = -1
		} else {
			b.wiggleDirection = 1
		}

		if closest.X == b.bbox.Min.X {
			resultVelocity.X = -AbsFloat(resultVelocity.X)
		}
		if closest.X == b.bbox.Max.X {
			resultVelocity.X = AbsFloat(resultVelocity.X)
		}
		if closest.Z == b.bbox.Min.Z {
			resultVelocity.Z = -AbsFloat(resultVelocity.Z)
		}
		if closest.Z == b.bbox.Max.Z {
			resultVelocity.Z = AbsFloat(resultVelocity.Z)
		}

		b.Lives--
	}
	return resultVelocity
}

func (b *Brick) DrawShadow(shader rl.Shader, modelLoc int32) {
	if b.reallyDead {
		return
	}
	shadowColor := b.game.shadowColor
	if b.Lives <= 0 {
		shadowColor = rl.Fade(shadowColor, float32(ease(float64(b.wiggleTime/MaxWiggleTime))))
	}

	rl.BeginShaderMode(shader)

	rl.SetShaderValueMatrix(shader, modelLoc, b.transformation)
	SetObjectColor(shader, shadowColor)

	// rl.DrawModelEx(b.game.ballModel, b.Position, rl.Vector3{X: 0, Y: 1, Z: 0}, 0, rl.Vector3{X: 1, Y: 1, Z: 1}, b.color)
	rl.DrawModel(b.game.brickShadowModel, rl.Vector3Zero(), 1, b.game.shadowColor)

	rl.EndShaderMode()
}

func (b *Brick) Draw(shader rl.Shader, modelLoc int32) {
	if b.reallyDead {
		return
	}
	baseColor := BrickColors[MaxInt(b.Lives-1, 0)]
	if b.Lives <= 0 {
		baseColor = rl.Fade(baseColor, float32(ease(float64(b.wiggleTime/MaxWiggleTime))))
	}

	rl.BeginShaderMode(shader)

	// Update uniforms
	SetObjectColor(shader, baseColor)

	// Or build full transform (translation + rotation + scale)

	rl.SetShaderValueMatrix(shader, modelLoc, b.transformation)

	// b.game.brickModel.Transform = rl.MatrixMultiply(transformation, rl.MatrixScale(b.scale.X, b.scale.Y, b.scale.Z))

	// rl.DrawModelEx(b.game.brickModel, b.Position, rl.Vector3{X: 0, Y: 1, Z: 0}, b.angle, b.scale, baseColor)
	rl.DrawModel(b.game.brickModel, rl.Vector3Zero(), 1, baseColor)

	// DrawNormals(b.game.brickModel.GetMeshes()[0], b.Position)

	rl.EndShaderMode()

	// rl.DrawCubeWires(b.Position, b.bbox.Max.X-b.bbox.Min.X, b.bbox.Max.Y-b.bbox.Min.Y, b.bbox.Max.Z-b.bbox.Min.Z, rl.Gray)
}
