package game

import (
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
)

var BrickColors = []rl.Color{rl.Green, rl.Yellow, rl.Orange, rl.Purple, rl.Red, rl.Blue}

type Brick struct {
	Position rl.Vector3
	Lives    int
	game     *Game
	model    rl.Model
	size     rl.Vector3
	scale    rl.Vector3
	bbox     rl.BoundingBox
}

func (b *Brick) Init(g *Game) {
	b.game = g

	b.model = rl.LoadModel("objects/rounded_cube.glb")
	for i := range int(b.model.MaterialCount) {
		b.model.GetMaterials()[i].Shader = g.shader
	}

	bbox := rl.GetModelBoundingBox(b.model)

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
}

func (b *Brick) checkCollision(ball Ball) rl.Vector3 {
	if b.Lives <= 0 {
		return ball.Velocity
	}

	// Find the point on the box closest to the sphere center
	closest := rl.Vector3Clamp(ball.Position, b.bbox.Min, b.bbox.Max)

	// Calculate distance from sphere center to closest point
	distance := rl.Vector3Length(rl.Vector3Subtract(ball.Position, closest))

	// Check for collision
	collides := distance <= BallRadius

	resultVelocity := rl.Vector3{X: ball.Velocity.X, Y: ball.Velocity.Y, Z: ball.Velocity.Z}

	if collides {
		if closest.X == b.bbox.Min.X || closest.X == b.bbox.Max.X {
			resultVelocity.X = -resultVelocity.X
		}
		if closest.Z == b.bbox.Min.Z || closest.Z == b.bbox.Max.Z {
			resultVelocity.Z = -resultVelocity.Z
		}

		b.Lives--
	}
	return resultVelocity
}

func (b *Brick) DrawShadow() {
	/*
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
	*/
}

func (b *Brick) Draw() {
	if b.Lives <= 0 {
		return
	}
	baseColor := BrickColors[b.Lives-1]

	shader := b.game.shader
	modelLoc := rl.GetShaderLocation(shader, "model")
	viewPosLoc := rl.GetShaderLocation(shader, "viewPos")
	lightPosLoc := rl.GetShaderLocation(shader, "lightPos")
	lightColorLoc := rl.GetShaderLocation(shader, "lightColor")
	objectColorLoc := rl.GetShaderLocation(shader, "objectColor")

	rl.BeginShaderMode(shader)

	// Update uniforms
	rl.SetShaderValue(shader, viewPosLoc, []float32{b.game.camera.Position.X, b.game.camera.Position.Y, b.game.camera.Position.Z}, rl.ShaderUniformVec3)
	rl.SetShaderValue(shader, lightPosLoc, []float32{b.game.lightPosition.X, b.game.lightPosition.Y, b.game.lightPosition.Z}, rl.ShaderUniformVec3)

	rl.SetShaderValue(shader, lightColorLoc, []float32{1, 1, 1}, rl.ShaderUniformVec3)
	rl.SetShaderValue(shader, objectColorLoc, []float32{float32(baseColor.R) / 255, float32(baseColor.G) / 255, float32(baseColor.B) / 255}, rl.ShaderUniformVec3)

	transform := rl.MatrixTranslate(b.Position.X, b.Position.Y, b.Position.Z)
	// Or build full transform (translation + rotation + scale)

	rl.SetShaderValueMatrix(shader, modelLoc, transform)

	rl.DrawModelEx(b.model, b.Position, rl.Vector3{X: 0, Y: 1, Z: 0}, 0, b.scale, baseColor)

	rl.EndShaderMode()

	// rl.DrawCubeWires(b.Position, b.bbox.Max.X-b.bbox.Min.X, b.bbox.Max.Y-b.bbox.Min.Y, b.bbox.Max.Z-b.bbox.Min.Z, rl.Gray)
}
