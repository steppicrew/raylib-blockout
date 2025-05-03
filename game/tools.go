package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func MinInt(a, b int32) int32 {
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

func MaxInt(a, b int32) int32 {
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
