package main

import (
	"math"

	m "github.com/tilleyd/rasteroids/math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const PLAYER_ACCEL = 640
const PLAYER_MAX_SPEED = 400
const PLAYER_TURN_SPEED = math.Pi
const PLAYER_RADIUS = 15
const BULLET_RADIUS = 2
const BULLET_SPEED = 500
const BULLET_COOLDOWN_S = 0.2
const BULLET_LIFETIME_S = 1.5

type Player struct {
	position  rl.Vector2
	velocity  rl.Vector2
	direction float32
	cooldown  float32
}

type Bullet struct {
	position rl.Vector2
	velocity rl.Vector2
	timer    float32
}

type Game struct {
	gameOver bool
	player   Player
	bullets  []Bullet
}

func main() {
	rl.InitWindow(1280, 720, "rasteroids")
	rl.SetWindowState(rl.FlagVsyncHint)
	defer rl.CloseWindow()

	game := NewGame()

	for !(rl.WindowShouldClose() || game.gameOver) {
		game.Update(rl.GetFrameTime())
		game.Draw()
	}
}

func NewGame() (g Game) {
	g.gameOver = false
	g.player.position.X = float32(rl.GetScreenWidth()) * 0.5
	g.player.position.Y = float32(rl.GetScreenHeight()) * 0.5
	return
}

func (g *Game) Update(delta float32) {
	xUnit := m.Cosf(g.player.direction)
	yUnit := m.Sinf(g.player.direction)

	if rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
		g.player.velocity.Y += PLAYER_ACCEL * yUnit * delta
		g.player.velocity.X += PLAYER_ACCEL * xUnit * delta
	}

	if rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
		g.player.velocity.Y -= PLAYER_ACCEL * yUnit * delta
		g.player.velocity.X -= PLAYER_ACCEL * xUnit * delta
	}

	if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		g.player.direction += PLAYER_TURN_SPEED * delta
		for g.player.direction > 2*math.Pi {
			g.player.direction -= 2 * math.Pi
		}
	}

	if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		g.player.direction -= PLAYER_TURN_SPEED * delta
		for g.player.direction < 0 {
			g.player.direction += 2 * math.Pi
		}
	}

	if rl.IsKeyDown(rl.KeySpace) && g.player.cooldown <= 0 {
		velocity := rl.Vector2{
			X: xUnit*BULLET_SPEED + g.player.velocity.X,
			Y: yUnit*BULLET_SPEED + g.player.velocity.Y,
		}
		bullet := Bullet{
			position: g.player.position,
			velocity: velocity,
			timer:    BULLET_LIFETIME_S,
		}

		g.bullets = append(g.bullets, bullet)
		g.player.cooldown = BULLET_COOLDOWN_S
	}

	if !rl.IsKeyDown(rl.KeySpace) {
		g.player.cooldown = 0.0
	}
	if g.player.cooldown > 0 {
		g.player.cooldown -= delta
	}

	g.player.velocity = rl.Vector2ClampValue(g.player.velocity, 0, PLAYER_MAX_SPEED)
	g.player.position.X += g.player.velocity.X * delta
	g.player.position.Y += g.player.velocity.Y * delta
	g.player.position = LazyWrap(g.player.position, PLAYER_RADIUS)

	for i := 0; i < len(g.bullets); i++ {
		g.bullets[i].position.X += g.bullets[i].velocity.X * delta
		g.bullets[i].position.Y += g.bullets[i].velocity.Y * delta
		g.bullets[i].position = LazyWrap(g.bullets[i].position, BULLET_RADIUS)
		g.bullets[i].timer -= delta

		if g.bullets[i].timer <= 0 {
			// delete by swapping with last element
			g.bullets[i] = g.bullets[len(g.bullets)-1]
			g.bullets = g.bullets[:len(g.bullets)-1]
		}
	}
}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	for _, bullet := range g.bullets {
		rl.DrawCircle(
			int32(bullet.position.X),
			int32(bullet.position.Y),
			BULLET_RADIUS,
			rl.White,
		)
	}

	unit := rl.Vector2{X: PLAYER_RADIUS, Y: 0}
	v1 := rl.Vector2Rotate(unit, g.player.direction)
	v2 := rl.Vector2Rotate(unit, g.player.direction-140*rl.Deg2rad)
	v3 := rl.Vector2Rotate(unit, g.player.direction+140*rl.Deg2rad)

	rl.DrawTriangle(
		rl.Vector2Add(g.player.position, v1),
		rl.Vector2Add(g.player.position, v2),
		rl.Vector2Add(g.player.position, v3),
		rl.White,
	)

	rl.EndDrawing()
}

func LazyWrap(v rl.Vector2, r float32) rl.Vector2 {
	if v.X > float32(rl.GetScreenWidth())+r {
		v.X -= float32(rl.GetScreenWidth()) + 2*r
	}
	if v.X < -r {
		v.X += float32(rl.GetScreenWidth()) + 2*r
	}
	if v.Y > float32(rl.GetScreenHeight())+r {
		v.Y -= float32(rl.GetScreenHeight()) + 2*r
	}
	if v.Y < -r {
		v.Y += float32(rl.GetScreenHeight()) + 2*r
	}

	return v
}
