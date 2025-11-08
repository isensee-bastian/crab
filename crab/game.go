package crab

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/isensee-bastian/crab/resources/images/sprites"
	"image"
	"image/color"
	_ "image/png"
	"math/rand/v2"
)

// Constant values like static screen positions are defined here.
const (
	ScreenWidth  = 1000
	ScreenHeight = 800

	beachScaleFactor = 2 // The beach image has half the size of our game screen, hence we double its size to fill it.

	crabStartX   = 450
	crabStartY   = 450
	crabStepSize = 2

	walkableMinY = 180 * beachScaleFactor // Vertical position where the sand area starts.
	walkableMaxY = 320 * beachScaleFactor // Vertical position where the sand area ends.

	scoreX = 10
	scoreY = 740
)

// Game holds our data required for managing state. All data, like images, object positions, and scores, go here.
type Game struct {
	beachImage *ebiten.Image

	crabImages     []*ebiten.Image // A slice (list) of images used for animations, first image at index 0, second at index 1 and so on.
	crabImageIndex int             // Index of the current crab image to show for the slice above. Required to show animations via alternating crab images.
	crabX          int             // The crabs current horizontal position.
	crabY          int             // The crabs current vertical position.

	fishImage *ebiten.Image
	fishX     int // The collectible fishes horizontal position.
	fishY     int // The collectible fishes vertical position.

	score int // The number of collected fishes.

	tickInSecond int // The tick count of the current second, reset to zero after one second has passed. Required for updating animation indexes, i.e. showing a single image for N ticks.
}

// NewGame prepares a fresh game state required for startup.
func NewGame() *Game {
	fishX, fishY := randomFishPosition()

	return &Game{
		beachImage: readImage(sprites.Beach),

		crabImages: readAnimationImages(sprites.Crab), // Load multiple images for representing animations into a slice.
		// Set the crabs start position.
		crabX: crabStartX,
		crabY: crabStartY,

		fishImage: readImage(sprites.Fish),
		// Set the first collectible fishes position.
		fishX: fishX,
		fishY: fishY,
	}
}

// Update processes all games rules, like checking user input and keeping score. All state updates must occur here, NOT in Draw.
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Signal that the game shall terminate normally when the Escape key is pressed.
		return ebiten.Termination
	}

	// Track ticks in second to determine which crab image to show to achieve an animation effect. We want to show each
	// crab image of the animation for an equal portion of time in a single second (per TPS).
	maxTicksPerSecond := ebiten.TPS()                         // Typically 60 ticks per second
	maxTicksPerFrame := maxTicksPerSecond / len(g.crabImages) // Each animation image shall take up an equal portion of a second
	g.tickInSecond = (g.tickInSecond + 1) % maxTicksPerSecond // The current tick in the current second, typically, between 0 and 59
	g.crabImageIndex = g.tickInSecond / maxTicksPerFrame      // Pick the next crab image if needed, ensuring that every animation image is shown in a second, Typically 0 to 59 / 15 = between 0 and 3 (inclusive)

	// Move crab according to pressed arrow keys. KeyPressDuration returns the number of ticks that passed since the
	// user started pressing the key (without releasing it). IsKeyJustPressed from above would not work for us here
	// since we want to keep the crab moving until a key is no longer pressed.
	if inpututil.KeyPressDuration(ebiten.KeyArrowLeft) > 0 {
		g.crabX = max(g.crabX-crabStepSize, 0) // Move no further left than where the screen starts (prevent negative X position).
	}
	if inpututil.KeyPressDuration(ebiten.KeyArrowRight) > 0 {
		g.crabX = min(g.crabX+crabStepSize, ScreenWidth-spriteWidth) // Move no further right than where the screen ends (include width of crab image, so it's still fully shown at the edge).
	}
	if inpututil.KeyPressDuration(ebiten.KeyArrowUp) > 0 {
		g.crabY = max(g.crabY-crabStepSize, walkableMinY) // Move no further up than where the sand area starts.
	}
	if inpututil.KeyPressDuration(ebiten.KeyArrowDown) > 0 {
		g.crabY = min(g.crabY+crabStepSize, walkableMaxY-spriteHeight) // Move no further down than where the sand area ends (include height of crab image, so it's still fully on the sand).
	}

	// Check for collision of crab and fish by comparing their rectangular areas.
	crabArea := image.Rect(g.crabX, g.crabY, g.crabX+spriteWidth-1, g.crabY+spriteHeight-1)
	fishArea := image.Rect(g.fishX, g.fishY, g.fishX+spriteWidth-1, g.fishY+spriteHeight-1)

	if crabArea.Overlaps(fishArea) {
		// There is some overlap between the rectangular crab area and the rectangular fish area.
		// That means, we collected the fish, remove it from its old position and spawn a new one at a different location.
		g.fishX, g.fishY = randomFishPosition()
		g.score += 1 // Increase score.
	}

	return nil
}

func randomFishPosition() (int, int) {
	// Determine a random position on the walkable sand area. Note that there is a small chance of immediately colliding
	// with our crabs position. This is fine for now and will lead to an immediate collect and respawn.
	randomX := rand.IntN(ScreenWidth - spriteWidth)
	randomY := walkableMinY + rand.IntN(walkableMaxY-walkableMinY-spriteHeight) // Must be in walkable sand area.

	return randomX, randomY
}

// Draw renders all game images to the screen according to the current game state.
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw beach first, all other images will be placed afterward to make it look like a background.
	{
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(0, 0)
		opts.GeoM.Scale(beachScaleFactor, beachScaleFactor)
		screen.DrawImage(g.beachImage, opts)
	}
	// Draw crab image.
	{
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(g.crabX), float64(g.crabY)) // Use current position from game state.
		screen.DrawImage(g.crabImages[g.crabImageIndex], opts)  // Draw animated crab by showing different crab images alternating.
	}
	// Draw collectible fish image.
	{
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(g.fishX), float64(g.fishY)) // Use current position from game state.
		screen.DrawImage(g.fishImage, opts)
	}
	// Draw current score label.
	drawBigText(screen, scoreX, scoreY, color.Black, fmt.Sprintf("Score: %d", g.score))
}

// Layout returns the logical screen size of the game. It can differ from the native outside size and will be scaled if needed.
func (g *Game) Layout(width, height int) (screenWidth, screenHeight int) {
	// No need to use a different logical screen size here. Our game size shall match the native outside window.
	return width, height
}
