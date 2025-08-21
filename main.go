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
const ASTEROID_DEVIATION = 0.1
const ASTEROID_LARGE_RADIUS = 88
const ASTEROID_MEDIUM_RADIUS = 64
const ASTEROID_SMALL_RADIUS = 44
const ASTEROID_MIN_VCOUNT = 16
const ASTEROID_MAX_VCOUNT = 32
const ASTEROID_SPEED = 100

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

type AsteroidSize int

const (
	ASTEROID_SIZE_SMALL AsteroidSize = iota
	ASTEROID_SIZE_MEDIUM
	ASTEROID_SIZE_LARGE
)

type Asteroid struct {
	size            AsteroidSize
	position        rl.Vector2
	velocity        rl.Vector2
	angle           float32
	angularVelocity float32
	maxRadius       float32
	vertices        []rl.Vector2
}

type Game struct {
	gameOver  bool
	player    Player
	bullets   []Bullet
	asteroids []Asteroid
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

	for range 3 {
		g.SpawnAsteroid()
	}
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
	}

	if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		g.player.direction -= PLAYER_TURN_SPEED * delta
	}
	g.player.direction = AngleWrap(g.player.direction)

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

	for i := range len(g.asteroids) {
		g.asteroids[i].position.X += g.asteroids[i].velocity.X * delta
		g.asteroids[i].position.Y += g.asteroids[i].velocity.Y * delta
		g.asteroids[i].position = LazyWrap(g.asteroids[i].position, g.asteroids[i].maxRadius)

		g.asteroids[i].angle += g.asteroids[i].angularVelocity * delta
		g.asteroids[i].angle = AngleWrap(g.asteroids[i].angle)
	}

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

	for _, asteroid := range g.asteroids {
		points := make([]rl.Vector2, len(asteroid.vertices)+2)
		points[0] = asteroid.position
		for i, v := range asteroid.vertices {
			points[i+1] = rl.Vector2Add(asteroid.position, rl.Vector2Rotate(v, asteroid.angle))
		}
		points[len(asteroid.vertices)+1] = points[1] // first vertex must also be the last point
		rl.DrawTriangleFan(points, rl.White)
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

func (g *Game) SpawnAsteroid() {
	position := rl.Vector2{
		X: m.RandRangef(0, float32(rl.GetScreenWidth())),
		Y: m.RandRangef(0, float32(rl.GetScreenHeight())),
	}
	direction := m.RandRangef(0, 2*math.Pi)
	velocity := rl.Vector2Rotate(rl.Vector2{X: 0, Y: ASTEROID_SPEED}, direction)
	angularVelocity := 30 * float32(rl.Deg2rad)
	if m.RandBool() {
		angularVelocity = -angularVelocity
	}

	asteroid := NewAsteroid(ASTEROID_SIZE_LARGE, position, velocity, angularVelocity)
	g.asteroids = append(g.asteroids, asteroid)
}

func NewAsteroid(
	size AsteroidSize,
	position rl.Vector2,
	velocity rl.Vector2,
	angularVelocity float32,
) (a Asteroid) {
	var radius float32
	switch size {
	case ASTEROID_SIZE_LARGE:
		radius = ASTEROID_LARGE_RADIUS
	case ASTEROID_SIZE_MEDIUM:
		radius = ASTEROID_MEDIUM_RADIUS
	case ASTEROID_SIZE_SMALL:
		radius = ASTEROID_SMALL_RADIUS
	default:
		panic("Unexpected asteroid size")
	}
	vertexCount := m.RandRangei(ASTEROID_MIN_VCOUNT, ASTEROID_MAX_VCOUNT)

	a.position = position
	a.velocity = velocity
	a.vertices = make([]rl.Vector2, vertexCount)
	a.angularVelocity = angularVelocity
	angleIncrement := 2 * math.Pi / float32(vertexCount)
	lowerRad := radius * (1 - ASTEROID_DEVIATION)
	upperRad := radius * (1 + ASTEROID_DEVIATION)
	for i := range vertexCount {
		angle := -float32(i) * angleIncrement // negate to force clockwise direction
		x := m.Cosf(angle)
		y := m.Sinf(angle)

		vRad := m.RandRangef(lowerRad, upperRad)
		if vRad > a.maxRadius {
			a.maxRadius = vRad
		}

		v := rl.Vector2Scale(rl.Vector2{X: x, Y: y}, vRad)
		a.vertices[i] = v
	}
	return
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

func AngleWrap(a float32) float32 {
	for a < 0 {
		a += 2 * math.Pi
	}
	for a > 2*math.Pi {
		a -= 2 * math.Pi
	}
	return a
}
