package game

import (
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BallRadius      float32 = 2
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
}

func (b *Ball) Init(g *Game) {
	b.game = g

	b.color = rl.Red
	b.colorLight = rl.ColorBrightness(b.color, 0.5)
	b.setVelocity(b.Velocity)
	b.min = rl.Vector3{X: BallRadius, Y: GamePlaneHeight, Z: BallRadius}
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
	// b.Position = newPosition
	return newPosition
}

func (b *Ball) DrawShadow() {
	/*
		shadowCenter := b.game.ProjectCanvas(b.Position)
		shadowRadius := b.game.ProjectCanvas(rl.Vector2{X: b.Position.X + BallRadius, Y: b.Position.Y}).X - shadowCenter.X
		rl.DrawCircleV(shadowCenter, shadowRadius, b.game.shadowColorLight)
		rl.DrawCircleV(shadowCenter, shadowRadius-shadowBorderThickness, b.game.shadowColor)
	*/
}

func (b *Ball) DrawNormals() {
	mesh := b.game.ballModel.GetMeshes()[0]

	// Access mesh data
	vertices := (*[1 << 30]float32)(unsafe.Pointer(mesh.Vertices))[:mesh.VertexCount*3]
	normals := (*[1 << 30]float32)(unsafe.Pointer(mesh.Normals))[:mesh.VertexCount*3]

	for i := 0; i < int(mesh.VertexCount); i++ {
		// Vertex position
		vx := vertices[i*3+0]
		vy := vertices[i*3+1]
		vz := vertices[i*3+2]
		pos := rl.Vector3{X: vx, Y: vy, Z: vz}

		// Normal direction
		nx := normals[i*3+0]
		ny := normals[i*3+1]
		nz := normals[i*3+2]
		normal := rl.Vector3{X: nx, Y: ny, Z: nz}

		// Transform position by model matrix
		worldPos := rl.Vector3Add(b.Position, rl.Vector3Scale(pos, 1.0)) // Apply model position

		// End of normal (just visualize scaled normal direction)
		normalEnd := rl.Vector3Add(worldPos, rl.Vector3Scale(normal, 1)) // scale for visibility

		// Draw line
		rl.DrawLine3D(worldPos, normalEnd, rl.Blue)
	}
}

func (b *Ball) Draw() {
	shader := b.game.shader
	modelLoc := rl.GetShaderLocation(shader, "model")
	objectColorLoc := rl.GetShaderLocation(shader, "objectColor")

	rl.BeginShaderMode(shader)

	// Update uniforms
	rl.SetShaderValue(shader, objectColorLoc, []float32{float32(b.color.R) / 255, float32(b.color.G) / 255, float32(b.color.B) / 255}, rl.ShaderUniformVec3)

	transform := rl.MatrixTranslate(b.Position.X, b.Position.Y, b.Position.Z)
	// Or build full transform (translation + rotation + scale)

	rl.SetShaderValueMatrix(shader, modelLoc, transform)

	rl.DrawModelEx(b.game.ballModel, b.Position, rl.Vector3{X: 0, Y: 1, Z: 0}, 0, rl.Vector3{X: 1, Y: 1, Z: 1}, b.color)

	rl.EndShaderMode()
}
