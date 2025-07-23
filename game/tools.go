package game

import (
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MinFloat(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaxFloat(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func AbsFloat(a float32) float32 {
	if a < 0 {
		return -a
	}
	return a
}

func Rotate(pos, center, axis rl.Vector3, angle float32) rl.Vector3 {
	return rl.Vector3Add(rl.Vector3RotateByAxisAngle(rl.Vector3Subtract(pos, center), axis, angle), center)
}

func CheckCollisionSphereBox(sphere rl.Vector3, radius float32, box rl.BoundingBox) (bool, rl.Vector3) {
	// Find the point on the box closest to the sphere center
	closest := rl.Vector3Clamp(sphere, box.Min, box.Max)

	// Calculate distance from sphere center to closest point
	distance := rl.Vector3Length(rl.Vector3Subtract(sphere, closest))

	// Sphere and box collide if the distance is less than the sphere radius
	return distance < radius, closest
}

func GetScaleBoundingBox(model rl.Model, position, dimensions rl.Vector3) (rl.Vector3, rl.BoundingBox) {
	bbox := rl.GetModelBoundingBox(model)

	objWidth := bbox.Max.X - bbox.Min.X
	objHeight := bbox.Max.Y - bbox.Min.Y
	objDepth := bbox.Max.Z - bbox.Min.Z

	scale := rl.NewVector3(dimensions.X/objWidth, dimensions.Y/objHeight, dimensions.Z/objDepth)
	half := rl.Vector3{X: dimensions.X / 2, Y: dimensions.Y / 2, Z: dimensions.Z / 2}
	realBbox := rl.BoundingBox{
		Min: rl.Vector3Subtract(position, half),
		Max: rl.Vector3Add(position, half),
	}
	return scale, realBbox
}

func DrawNormals(mesh rl.Mesh, position rl.Vector3) {

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
		worldPos := rl.Vector3Add(position, rl.Vector3Scale(pos, 1.0)) // Apply model position

		// End of normal (just visualize scaled normal direction)
		normalEnd := rl.Vector3Add(worldPos, rl.Vector3Scale(normal, 2)) // scale for visibility

		// Draw line
		rl.DrawLine3D(worldPos, normalEnd, rl.Blue)
	}
}

func SetObjectColor(shader rl.Shader, color rl.Color) {
	objectColorLoc := rl.GetShaderLocation(shader, "objectColor")
	rl.SetShaderValue(shader, objectColorLoc, []float32{float32(color.R) / 255, float32(color.G) / 255, float32(color.B) / 255, float32(color.A) / 255}, rl.ShaderUniformVec4)
}

func FixCameraUp(camera *rl.Camera3D, querVector rl.Vector3) {
	camDirection := rl.Vector3Normalize(rl.Vector3Subtract(camera.Position, camera.Target))
	camera.Up = rl.Vector3Negate(rl.Vector3Normalize(rl.Vector3CrossProduct(camDirection, querVector)))
}
