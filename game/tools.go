package game

import (
	"math"

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

func normalize(v rl.Vector2) rl.Vector2 {
	l := float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
	return scale(v, 1/l)
}

func scale(v rl.Vector2, factor float32) rl.Vector2 {
	return rl.Vector2{X: v.X * factor, Y: v.Y * factor}
}

func add(a, b rl.Vector2) rl.Vector2 {
	return rl.Vector2{X: a.X + b.X, Y: a.Y + b.Y}
}

func crop(v, min, max rl.Vector2) rl.Vector2 {
	return rl.Vector2{X: MaxFloat(MinFloat(v.X, max.X), min.X), Y: MaxFloat(MinFloat(v.Y, max.Y), min.Y)}
}
