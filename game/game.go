package game

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	BrickRows                     = 5
	HeightOffset          float32 = .8
	LightHeight           float32 = 30
	GamePlaneHeight       float32 = 2
	CameraHeight          float32 = 80
	shadowBorderThickness float32 = 2
	CameraSpeed           float32 = 0.5
)

type Game struct {
	Width         int32
	Height        int32
	ball          Ball
	paddle        Paddle
	bricks        []*Brick
	lightPosition rl.Vector3
	camera        rl.Camera3D
	shadowColor   rl.Color

	ballModel        rl.Model
	ballShadowModel  rl.Model
	brickModel       rl.Model
	brickShadowModel rl.Model
	defaultShader    rl.Shader
	shadowShader     rl.Shader
}

func (g *Game) InitiModels() {
	defaultShader := rl.LoadShader("shader/defaultVertex.glsl", "shader/defaultFragment.glsl")
	g.defaultShader = defaultShader

	g.ballModel = rl.LoadModelFromMesh(rl.GenMeshSphere(BallRadius, 32, 32))
	for i := range int(g.ballModel.MaterialCount) {
		g.ballModel.GetMaterials()[i].Shader = defaultShader
	}

	g.brickModel = rl.LoadModel("objects/rounded_cube.glb")
	for i := range int(g.brickModel.MaterialCount) {
		g.brickModel.GetMaterials()[i].Shader = defaultShader
	}

	shadowShader := rl.LoadShader("shader/shadowVertex.glsl", "shader/shadowFragment.glsl")
	g.shadowShader = shadowShader

	g.ballShadowModel = rl.LoadModelFromMesh(rl.GenMeshSphere(BallRadius, 32, 32))
	for i := range int(g.ballShadowModel.MaterialCount) {
		g.ballShadowModel.GetMaterials()[i].Shader = shadowShader
	}

	g.brickShadowModel = rl.LoadModel("objects/rounded_cube.glb")
	for i := range int(g.brickShadowModel.MaterialCount) {
		g.brickShadowModel.GetMaterials()[i].Shader = shadowShader
	}
}

func (g *Game) Init() {
	g.InitiModels()

	g.lightPosition = rl.Vector3{X: float32(g.Width) / 2, Y: LightHeight, Z: float32(g.Height) + 2}
	// g.lightPosition = rl.Vector3{X: 0, Y: LightHeight, Z: float32(g.Height) / 2}

	// g.shadowMatrix = getShadowMatrix(g.lightPosition)
	g.shadowColor = rl.LightGray

	g.camera = rl.Camera3D{}
	g.camera.Position = rl.Vector3{X: float32(g.Width) / 2, Y: CameraHeight, Z: float32(g.Height) / 2}
	g.camera.Target = rl.Vector3{X: float32(g.Width) / 2, Y: 0, Z: float32(g.Height) / 2}
	g.camera.Up = rl.NewVector3(0.0, 0.0, 1.0) // Camera up vector (relative to target)
	g.camera.Fovy = 45.0                       // Camera field-of-view Y
	g.camera.Projection = rl.CameraPerspective // Camera projection type

	/*
		g.lightPosition = rl.Vector3{X: 0, Y: 0, Z: 0}
		g.camera.Position = rl.Vector3{X: 0, Y: 0, Z: 0}
	*/

	g.ball = Ball{
		Position: rl.Vector3{X: float32(g.Width) / 2, Y: GamePlaneHeight, Z: float32(g.Height) * 0.3},
		// Position: rl.Vector3{X: 0, Y: 0, Z: 0},
		Velocity: rl.Vector3{X: 1, Y: 0, Z: 1},
		// Velocity: rl.Vector3{X: 0, Y: 0, Z: 0},
	}
	g.paddle = Paddle{Position: rl.Vector3{X: float32(g.Width) / 2, Y: GamePlaneHeight, Z: PaddleHeight / 2}}
	g.bricks = []*Brick{}

	g.ball.Init(g)
	g.paddle.Init(g)

	if true {
		cols := int(g.Width / int32(BrickWidth))
		g.bricks = make([]*Brick, BrickRows*cols)
		for z := 0; z < BrickRows; z++ {
			for x := 0; x < cols; x++ {
				brick := Brick{
					Position: rl.Vector3{X: float32(x)*BrickWidth + BrickWidth/2, Y: GamePlaneHeight, Z: float32(g.Height)*HeightOffset - float32(z)*BrickHeight},
					Lives:    BrickRows - z,
				}
				brick.Init(g)
				g.bricks[z*cols+x] = &brick
			}
		}
	} else {
		g.bricks = make([]*Brick, 1)
		brick := Brick{
			Position: rl.Vector3{X: float32(g.Width) / 2, Y: GamePlaneHeight, Z: float32(g.Height) / 2},
			Lives:    5,
		}
		brick.Init(g)
		g.bricks[0] = &brick
	}

}

func (g *Game) updateCamera(time float32) {
	if rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl) {
		pos := g.camera.Position
		target := g.camera.Target
		if rl.IsKeyDown(rl.KeyUp) {
			g.camera.Position = Rotate(pos, target, rl.Vector3{X: 1, Y: 0, Z: 0}, time*CameraSpeed)
		}
		if rl.IsKeyDown(rl.KeyDown) {
			g.camera.Position = Rotate(pos, target, rl.Vector3{X: 1, Y: 0, Z: 0}, -time*CameraSpeed)
		}
		if rl.IsKeyDown(rl.KeyRight) {
			g.camera.Position = Rotate(pos, target, rl.Vector3{X: 0, Y: 0, Z: 1}, time*CameraSpeed)
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			g.camera.Position = Rotate(pos, target, rl.Vector3{X: 0, Y: 0, Z: 1}, -time*CameraSpeed)
		}
	}

}

func (g *Game) Update(time float32) {
	g.updateCamera(time)

	ballPosition := g.ball.Update(time)
	if ballPosition.X < g.ball.min.X || ballPosition.X > g.ball.max.X {
		g.ball.Velocity.X = -g.ball.Velocity.X
	}
	if ballPosition.Z < g.ball.min.Z || ballPosition.Z > g.ball.max.Z {
		g.ball.Velocity.Z = -g.ball.Velocity.Z
	}

	g.paddle.Update(time)
	g.ball.setVelocity(g.paddle.checkCollision(g.ball))

	for _, brick := range g.bricks {
		g.ball.setVelocity(brick.checkCollision(g.ball))

		brick.Update(time)
	}
}

func (g *Game) Draw() {
	// rl.UpdateCamera(&g.camera, rl.CameraOrbital)
	// g.camera.Target = g.ball.Position

	rl.BeginMode3D(g.camera) // Enter 3D mode

	shader := g.defaultShader
	viewPosLoc := rl.GetShaderLocation(shader, "viewPos")
	lightPosLoc := rl.GetShaderLocation(shader, "lightPos")
	lightColorLoc := rl.GetShaderLocation(shader, "lightColor")
	modelLoc := rl.GetShaderLocation(shader, "model")

	// Update uniforms
	rl.SetShaderValue(shader, viewPosLoc, []float32{g.camera.Position.X, g.camera.Position.Y, g.camera.Position.Z}, rl.ShaderUniformVec3)
	rl.SetShaderValue(shader, lightPosLoc, []float32{g.lightPosition.X, g.lightPosition.Y, g.lightPosition.Z}, rl.ShaderUniformVec3)
	rl.SetShaderValue(shader, lightColorLoc, []float32{1, 1, 1}, rl.ShaderUniformVec3)

	g.ball.Draw(shader, modelLoc)
	g.paddle.Draw(shader, modelLoc)
	for _, brick := range g.bricks {
		brick.Draw(shader, modelLoc)
	}

	// Draw a grid to visualize the 3D space (optional)
	// rl.DrawGrid(100, 1.0)

	rl.DrawSphere(g.lightPosition, 1, rl.Fade(rl.Red, .3))

	rl.EndMode3D() // Exit 3D mode
}

func (g *Game) DrawShadow() {
	// rl.UpdateCamera(&g.camera, rl.CameraOrbital)
	// g.camera.Target = g.ball.Position

	rl.BeginMode3D(g.camera) // Enter 3D mode

	shader := g.shadowShader
	lightPosLoc := rl.GetShaderLocation(shader, "lightPos")
	modelLoc := rl.GetShaderLocation(shader, "model")

	// Update uniforms
	rl.SetShaderValue(shader, lightPosLoc, []float32{g.lightPosition.X, g.lightPosition.Y, g.lightPosition.Z}, rl.ShaderUniformVec3)
	SetObjectColor(shader, g.shadowColor)

	g.ball.DrawShadow(shader, modelLoc)
	g.paddle.DrawShadow(shader, modelLoc)
	for _, brick := range g.bricks {
		brick.DrawShadow(shader, modelLoc)
	}

	rl.EndMode3D() // Exit 3D mode
}
