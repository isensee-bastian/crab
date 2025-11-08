package crab

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/isensee-bastian/crab/resources/images/sprites"
	_ "image/png"
)

// Constant values like static screen positions are defined here.
const (
	ScreenWidth  = 1000
	ScreenHeight = 800

	crabStartX = 450
	crabStartY = 450
)

// Game holds our data required for managing state. All data, like images, object positions, and scores, go here.
type Game struct {
	beachImage *ebiten.Image
	crabImages []*ebiten.Image // A slice (list) of images used for animations, first image at index 0, second at index 1 and so on.
	crabX      int             // The crabs current horizontal position.
	crabY      int             // The crabs current vertical position.
}

// NewGame prepares a fresh game state required for startup.
func NewGame() *Game {
	return &Game{
		beachImage: readImage(sprites.Beach),
		crabImages: readAnimationImages(sprites.Crab), // Load multiple images for representing animations into a slice.
		// Set the crabs start position.
		crabX: crabStartX,
		crabY: crabStartY,
	}
}

// Update processes all games rules, like checking user input and keeping score. All state updates must occur here, NOT in Draw.
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Signal that the game shall terminate normally when the Escape key is pressed.
		return ebiten.Termination
	}

	return nil
}

// Draw renders all game images to the screen according to the current game state.
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw beach first, all other images will be placed afterward to make it look like a background.
	{
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(0, 0)
		opts.GeoM.Scale(2, 2)
		screen.DrawImage(g.beachImage, opts)
	}
	// Draw crab image.
	{
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(g.crabX), float64(g.crabY)) // Use current position from game state.
		screen.DrawImage(g.crabImages[0], opts)                 // For now, just draw the first crab image (no animation yet).
	}
}

// Layout returns the logical screen size of the game. It can differ from the native outside size and will be scaled if needed.
func (g *Game) Layout(width, height int) (screenWidth, screenHeight int) {
	// No need to use a different logical screen size here. Our game size shall match the native outside window.
	return width, height
}
