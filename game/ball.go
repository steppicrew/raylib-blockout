package game

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BallRadius      float32 = 10
	BallRadiusLight float32 = 5
	BallSpeed               = 200
)

type Ball struct {
	Position   rl.Vector2
	Velocity   rl.Vector2
	game       *Game
	min        rl.Vector2
	max        rl.Vector2
	color      rl.Color
	colorLight rl.Color
}

func (b *Ball) Init(g *Game) {
	b.game = g
	b.color = rl.Red
	b.colorLight = rl.ColorBrightness(b.color, 0.3)
	b.setVelocity(b.Velocity)
	b.min = rl.Vector2{X: BallRadius, Y: BallRadius}
	b.max = rl.Vector2{X: float32(b.game.Width) - BallRadius, Y: float32(b.game.Height) - BallRadius}
}

func (b *Ball) setVelocity(v rl.Vector2) {
	b.Velocity = scale(normalize(v), BallSpeed)
}

func (b *Ball) Update(time float32) rl.Vector2 {
	newPosition := add(b.Position, scale(b.Velocity, time))
	b.Position = crop(newPosition, b.min, b.max)
	return newPosition
}

func (b *Ball) DrawShadow() {
	shadowCenter := b.game.ProjectCanvas(b.Position)
	shadowRadius := b.game.ProjectCanvas(rl.Vector2{X: b.Position.X + BallRadius, Y: b.Position.Y}).X - shadowCenter.X
	rl.DrawCircleV(shadowCenter, shadowRadius, b.game.shadowColorLight)
	rl.DrawCircleV(shadowCenter, shadowRadius-shadowBorderThickness, b.game.shadowColor)
}

func (b *Ball) Draw() {
	camera := rl.Camera3D{}
	camera.Position = rl.Vector3{X: float32(b.game.Width) / 2, Y: float32(b.game.Height), Z: 100}
	camera.Target = rl.Vector3{X: float32(b.game.Width) / 2, Y: float32(b.game.Height) / 2, Z: CanvasZ}
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective
	rl.BeginMode3D(camera)
	rl.DrawSphere(rl.Vector3{X: b.Position.X, Y: b.Position.Y, Z: CanvasZ}, BallRadius, b.color)
	rl.EndMode3D()

	// rl.DrawCircleV(b.Position, BallRadius, b.color)
	// rl.DrawCircleV(rl.Vector2{X: b.Position.X - 1, Y: b.Position.Y - 2}, BallRadiusLight, b.colorLight)
	lightCenter := b.game.ProjectZ(b.Position, CanvasZ+BallRadius)
	distance := rl.Vector2{X: b.Position.X - lightCenter.X, Y: b.Position.Y - lightCenter.Y}
	distanceLength := float32(math.Sqrt(float64(distance.X*distance.X + distance.Y*distance.Y)))
	var lightOffset rl.Vector2
	if distanceLength == 0 {
		lightOffset = rl.Vector2{X: 0, Y: 0}
	} else {
		realDist := BallRadius - (BallRadius / (distanceLength + 1))
		lightOffset = rl.Vector2{
			X: distance.X / distanceLength * realDist,
			Y: distance.Y / distanceLength * realDist,
		}
	}
	rl.DrawEllipse(
		int32(b.Position.X-lightOffset.X),
		int32(b.Position.Y-lightOffset.Y),
		(BallRadius-float32(math.Abs(float64(lightOffset.X))))/BallRadius*BallRadiusLight,
		(BallRadius-float32(math.Abs(float64(lightOffset.Y))))/BallRadius*BallRadiusLight,
		b.colorLight,
	)
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
