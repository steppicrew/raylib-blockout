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
	model      rl.Model
	min        rl.Vector3
	max        rl.Vector3
	color      rl.Color
	colorLight rl.Color
}

func (b *Ball) Init(g *Game) {
	b.game = g
	b.model = rl.LoadModelFromMesh(rl.GenMeshSphere(BallRadius, 32, 32))
	for i := range int(b.model.MaterialCount) {
		b.model.GetMaterials()[i].Shader = g.shader
	}

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
	mesh := b.model.GetMeshes()[0]

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
	viewPosLoc := rl.GetShaderLocation(shader, "viewPos")
	lightPosLoc := rl.GetShaderLocation(shader, "lightPos")
	lightColorLoc := rl.GetShaderLocation(shader, "lightColor")
	objectColorLoc := rl.GetShaderLocation(shader, "objectColor")

	rl.BeginShaderMode(shader)

	// Update uniforms
	rl.SetShaderValue(shader, viewPosLoc, []float32{b.game.camera.Position.X, b.game.camera.Position.Y, b.game.camera.Position.Z}, rl.ShaderUniformVec3)
	rl.SetShaderValue(shader, lightPosLoc, []float32{b.game.lightPosition.X, b.game.lightPosition.Y, b.game.lightPosition.Z}, rl.ShaderUniformVec3)

	rl.SetShaderValue(shader, lightColorLoc, []float32{1, 1, 1}, rl.ShaderUniformVec3)
	rl.SetShaderValue(shader, objectColorLoc, []float32{float32(b.color.R) / 255, float32(b.color.G) / 255, float32(b.color.B) / 255}, rl.ShaderUniformVec3)

	transform := rl.MatrixTranslate(b.Position.X, b.Position.Y, b.Position.Z)
	// Or build full transform (translation + rotation + scale)

	rl.SetShaderValueMatrix(shader, modelLoc, transform)

	rl.DrawModelEx(b.model, b.Position, rl.Vector3{X: 0, Y: 1, Z: 0}, 0, rl.Vector3{X: 1, Y: 1, Z: 1}, b.color)

	rl.EndShaderMode()
}

/*
package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"math"
)

func main() {
	rl.InitWindow(800, 600, "Ball with Reflection")
	rl.SetTargetFPS(60)

	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(0.0, 10.0, 10.0)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	ballPos := rl.NewVector3(0, 0, 0)
	ballRadius := float32(2.0)
	lightPos := rl.NewVector3(5, 5, 5)

	for !rl.WindowShouldClose() {
		// Update light position (optional)
		// lightPos.X = float32(math.Sin(float64(rl.GetTime()))) * 5
		// lightPos.Z = float32(math.Cos(float64(rl.GetTime()))) * 5

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		// Draw the ball
		rl.DrawSphere(ballPos, ballRadius, rl.Blue)

		// Calculate reflection circle position
		lightDir := rl.Vector3Subtract(lightPos, ballPos)
		lightDir = rl.Vector3Normalize(lightDir)

		// The reflection point is in the direction of the light's reflection
		reflectionPoint := rl.Vector3Add(ballPos, rl.Vector3Scale(lightDir, ballRadius))

		// Draw the reflection circle
		reflectionRadius := ballRadius * 0.3 // Size of the reflection
		drawReflectionCircle(ballPos, reflectionPoint, reflectionRadius, ballRadius, rl.White)

		// Draw light source for reference
		rl.DrawSphere(lightPos, 0.2, rl.Yellow)

		rl.EndMode3D()

		rl.DrawText("Ball with Reflection Highlight", 10, 10, 20, rl.DarkGray)
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func drawReflectionCircle(ballCenter, reflectionPoint rl.Vector3, reflectionRadius, ballRadius float32, color rl.Color) {
	// Calculate the normal direction from ball center to reflection point
	normal := rl.Vector3Subtract(reflectionPoint, ballCenter)
	normal = rl.Vector3Normalize(normal)

	// Create two perpendicular vectors to the normal
	var tangent1 rl.Vector3
	if math.Abs(float64(normal.X)) > math.Abs(float64(normal.Z)) {
		tangent1 = rl.NewVector3(-normal.Y, normal.X, 0)
	} else {
		tangent1 = rl.NewVector3(0, -normal.Z, normal.Y)
	}
	tangent1 = rl.Vector3Normalize(tangent1)
	tangent2 := rl.Vector3CrossProduct(normal, tangent1)

	// Draw the circle using line segments
	segments := 32
	prevPoint := rl.Vector3Add(reflectionPoint, rl.Vector3Scale(tangent1, reflectionRadius))

	for i := 1; i <= segments; i++ {
		angle := float32(i) * 2 * math.Pi / float32(segments)
		offset := rl.Vector3Add(
			rl.Vector3Scale(tangent1, reflectionRadius*float32(math.Cos(float64(angle)))),
			rl.Vector3Scale(tangent2, reflectionRadius*float32(math.Sin(float64(angle)))),
		)
		currentPoint := rl.Vector3Add(reflectionPoint, offset)

		rl.DrawLine3D(prevPoint, currentPoint, color)
		prevPoint = currentPoint
	}
}

*/
