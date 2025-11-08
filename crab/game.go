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
)

// Game holds our data required for managing state. All data, like images, object positions, and scores, go here.
type Game struct {
	beachImage *ebiten.Image // Hold image bytes in game state for drawing.
}

// NewGame prepares a fresh game state required for startup.
func NewGame() *Game {
	return &Game{
		beachImage: readImage(sprites.Beach), // Load image from sprites folder into memory.
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
	opts := &ebiten.DrawImageOptions{}   // Image properties are configured via DrawImageOptions.
	opts.GeoM.Translate(0, 0)            // Place image at position x=0, y=0 (upper left corner).
	opts.GeoM.Scale(2, 2)                // Scale image by factor 2 (double its size).
	screen.DrawImage(g.beachImage, opts) // Draw the actual image using our specified options.
}

// Layout returns the logical screen size of the game. It can differ from the native outside size and will be scaled if needed.
func (g *Game) Layout(width, height int) (screenWidth, screenHeight int) {
	// No need to use a different logical screen size here. Our game size shall match the native outside window.
	return width, height
}
