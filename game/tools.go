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
