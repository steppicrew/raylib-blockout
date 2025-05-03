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
	half := rl.Vector3{X: objWidth / 2, Y: objHeight / 2, Z: objDepth / 2}
	b.bbox = rl.BoundingBox{
		Min: rl.Vector3Subtract(b.Position, half),
		Max: rl.Vector3Add(b.Position, half),
	}

	b.size = rl.Vector3{X: BrickWidth, Y: BrickHeight, Z: BrickDepth}
}

func (b *Brick) Update(time float32) {
}

func (b *Brick) checkCollision(ball Ball) rl.Vector3 {
	topLeft := rl.Vector2{X: b.bbox.Max.X, Y: b.bbox.Max.Z}
	topRight := rl.Vector2{X: b.bbox.Min.X, Y: b.bbox.Max.Z}
	bottomLeft := rl.Vector2{X: b.bbox.Max.X, Y: b.bbox.Min.Z}
	bottomRight := rl.Vector2{X: b.bbox.Min.X, Y: b.bbox.Min.Z}
	resultVelocity := rl.Vector3{X: ball.Velocity.X, Y: ball.Velocity.Y, Z: ball.Velocity.Z}
	ballPosition := rl.Vector2{X: ball.Position.X, Y: ball.Position.Z}
	bounced := false
	if rl.CheckCollisionCircleLine(ballPosition, BallRadius, topLeft, topRight) {
		resultVelocity.Z = AbsFloat(resultVelocity.Z)
		bounced = true
	}
	if rl.CheckCollisionCircleLine(ballPosition, BallRadius, bottomLeft, bottomRight) {
		resultVelocity.Z = -AbsFloat(resultVelocity.Z)
		bounced = true
	}
	if rl.CheckCollisionCircleLine(ballPosition, BallRadius, topRight, bottomRight) {
		resultVelocity.X = AbsFloat(resultVelocity.X)
		bounced = true
	}
	if rl.CheckCollisionCircleLine(ballPosition, BallRadius, bottomLeft, topLeft) {
		resultVelocity.X = -AbsFloat(resultVelocity.X)
		bounced = true
	}
	if bounced && b.Lives > 0 {
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
	// highlightColor := rl.ColorBrightness(baseColor, 0.3)
	// Draw the 3D brick
	// Parameters: position, width, height, depth, color

	/*
		rl.DrawCubeV(b.Position, b.size, baseColor)
		rl.DrawCubeWiresV(b.Position, b.size, highlightColor)
	*/

	/*
		start := rl.Vector3{X: b.Position.X + b.size.Y/2, Y: b.Position.Y + b.size.Y/2, Z: b.Position.Z + b.size.Z/2}
		end := rl.Vector3Add(start, rl.Vector3{X: b.size.X - 2*+b.size.Y/2, Y: 0, Z: 0})
		rl.DrawCapsule(start, end, b.size.Y/2, 8, 10, baseColor)
		rl.DrawCapsuleWires(start, end, b.size.Y/2, 8, 10, highlightColor)
	*/

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

}
