// +build ignore

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Card dimensions
const (
	// Icon: 32x48
	IconWidth  = 32
	IconHeight = 48
	
	// Small: 64x90
	SmallWidth  = 64
	SmallHeight = 90
	
	// Large: 200x280
	LargeWidth  = 200
	LargeHeight = 280
)

// Colors
var (
	white = color.RGBA{255, 255, 255, 255}
	black = color.RGBA{0, 0, 0, 255}
	red   = color.RGBA{220, 20, 60, 255}
	darkBlue = color.RGBA{0, 0, 139, 255}
)

type CardInfo struct {
	Rank     string
	RankVal  int
	Suit     string
	Symbol   string
	Color    color.RGBA
	FileName string
}

func getFontFace(ttf []byte, size float64) font.Face {
	f, err := opentype.Parse(ttf)
	if err != nil {
		log.Fatal(err)
	}
	
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size: size,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}
	
	return face
}

func main() {
	suits := []struct {
		name   string
		symbol string
		color  color.RGBA
		id     int
	}{
		{"hearts", "♥", red, 0},
		{"diamonds", "♦", red, 1},
		{"clubs", "♣", black, 2},
		{"spades", "♠", black, 3},
	}
	
	ranks := []struct {
		name  string
		short string
		id    int
	}{
		{"ace", "A", 1},
		{"2", "2", 2},
		{"3", "3", 3},
		{"4", "4", 4},
		{"5", "5", 5},
		{"6", "6", 6},
		{"7", "7", 7},
		{"8", "8", 8},
		{"9", "9", 9},
		{"10", "10", 10},
		{"jack", "J", 11},
		{"queen", "Q", 12},
		{"king", "K", 13},
	}
	
	// Generate card for each combination
	cardCount := 0
	for _, suit := range suits {
		for _, rank := range ranks {
			card := CardInfo{
				Rank:     rank.short,
				RankVal:  rank.id,
				Suit:     suit.name,
				Symbol:   suit.symbol,
				Color:    suit.color,
				FileName: fmt.Sprintf("%d_%d", rank.id, suit.id), // e.g., "1_0" for Ace of Hearts
			}
			
			// Generate card quietly
			
			// Generate all three sizes
			generateCard(card, IconWidth, IconHeight, "icon")
			generateCard(card, SmallWidth, SmallHeight, "small")
			generateCard(card, LargeWidth, LargeHeight, "large")
			cardCount++
		}
	}
	
	// Generate card back for all sizes
	generateCardBack(IconWidth, IconHeight, "icon")
	generateCardBack(SmallWidth, SmallHeight, "small")
	generateCardBack(LargeWidth, LargeHeight, "large")
	
	fmt.Println("Card generation complete!")
	fmt.Printf("Generated %d cards in 3 sizes (should be 52)\n", cardCount)
}

func generateCard(card CardInfo, width, height int, size string) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Fill with white background
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)
	
	// Draw border
	drawBorder(img, black, 1)
	
	// Get appropriate font sizes
	var rankFont, suitFont font.Face
	var rankOffset, suitOffset int
	switch size {
	case "icon":
		rankFont = getFontFace(gobold.TTF, 8)
		suitFont = getFontFace(goregular.TTF, 10)
		rankOffset = 10
		suitOffset = 18
	case "small":
		rankFont = getFontFace(gobold.TTF, 14)
		suitFont = getFontFace(goregular.TTF, 16)
		rankOffset = 16
		suitOffset = 30
	case "large":
		rankFont = getFontFace(gobold.TTF, 28)
		suitFont = getFontFace(goregular.TTF, 32)
		rankOffset = 32
		suitOffset = 58
	}
	
	// Draw rank and suit in top-left corner
	margin := 4
	drawTextWithFont(img, card.Rank, margin, rankOffset, card.Color, rankFont)
	drawTextWithFont(img, card.Symbol, margin, suitOffset, card.Color, suitFont)
	
	// Draw rank and suit in bottom-right corner (upside down appearance)
	// Calculate text width for proper positioning
	rankWidth := getTextWidth(card.Rank, rankFont)
	suitWidth := getTextWidth(card.Symbol, suitFont)
	
	drawTextWithFont(img, card.Rank, width-rankWidth-margin, height-suitOffset+rankOffset, card.Color, rankFont)
	drawTextWithFont(img, card.Symbol, width-suitWidth-margin, height-margin, card.Color, suitFont)
	
	// Draw center symbols
	drawCenterSymbols(img, card, size)
	
	// Save the image
	outputPath := filepath.Join("static", "cards", size, card.FileName+".png")
	file, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Error creating file %s: %v", outputPath, err)
		return
	}
	defer file.Close()
	
	if err := png.Encode(file, img); err != nil {
		log.Printf("Error encoding PNG: %v", err)
	}
}

func generateCardBack(width, height int, size string) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Fill with dark blue background
	draw.Draw(img, img.Bounds(), &image.Uniform{darkBlue}, image.Point{}, draw.Src)
	
	// Draw white border
	drawBorder(img, white, 2)
	
	// Draw decorative pattern
	drawBackPattern(img, size)
	
	// Save the image
	outputPath := filepath.Join("static", "cards", size, "back.png")
	file, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Error creating file %s: %v", outputPath, err)
		return
	}
	defer file.Close()
	
	if err := png.Encode(file, img); err != nil {
		log.Printf("Error encoding PNG: %v", err)
	}
}

func drawBorder(img *image.RGBA, c color.RGBA, thickness int) {
	bounds := img.Bounds()
	for i := 0; i < thickness; i++ {
		// Top and bottom
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, bounds.Min.Y+i, c)
			img.Set(x, bounds.Max.Y-1-i, c)
		}
		// Left and right
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			img.Set(bounds.Min.X+i, y, c)
			img.Set(bounds.Max.X-1-i, y, c)
		}
	}
}

func drawTextWithFont(img *image.RGBA, text string, x, y int, c color.RGBA, face font.Face) {
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: face,
		Dot:  point,
	}
	d.DrawString(text)
}

func getTextWidth(text string, face font.Face) int {
	d := &font.Drawer{Face: face}
	return d.MeasureString(text).Round()
}

func drawCenterSymbols(img *image.RGBA, card CardInfo, size string) {
	bounds := img.Bounds()
	
	var suitFont font.Face
	var symbolOffsetX, symbolOffsetY int
	switch size {
	case "icon":
		suitFont = getFontFace(goregular.TTF, 10)
		symbolOffsetX = 4
		symbolOffsetY = 7
	case "small":
		suitFont = getFontFace(goregular.TTF, 18)
		symbolOffsetX = 7
		symbolOffsetY = 12
	case "large":
		suitFont = getFontFace(goregular.TTF, 48)
		symbolOffsetX = 18
		symbolOffsetY = 32
	}
	
	// Get positions for symbols based on card rank
	positions := getSuitPositions(card.RankVal, bounds, size)
	
	// Draw each symbol
	for _, pos := range positions {
		drawTextWithFont(img, card.Symbol, pos.X-symbolOffsetX, pos.Y+symbolOffsetY, card.Color, suitFont)
	}
	
	// For face cards, draw a simple letter in the center
	if card.RankVal >= 11 {
		var faceFont font.Face
		switch size {
		case "icon":
			faceFont = getFontFace(gobold.TTF, 16)
		case "small":
			faceFont = getFontFace(gobold.TTF, 28)
		case "large":
			faceFont = getFontFace(gobold.TTF, 72)
		}
		
		centerX := bounds.Max.X / 2
		centerY := bounds.Max.Y / 2
		letterWidth := getTextWidth(card.Rank, faceFont)
		drawTextWithFont(img, card.Rank, centerX-letterWidth/2, centerY+symbolOffsetY/2, card.Color, faceFont)
	}
}

func getSuitPositions(rank int, bounds image.Rectangle, size string) []image.Point {
	width := bounds.Max.X
	height := bounds.Max.Y
	centerX := width / 2
	centerY := height / 2
	
	// Define margins based on size to avoid overlap with corner text
	var marginTop, marginBottom, marginSide int
	switch size {
	case "icon":
		marginTop = 20
		marginBottom = 8
		marginSide = 4
	case "small":
		marginTop = 34
		marginBottom = 14
		marginSide = 8
	case "large":
		marginTop = 70
		marginBottom = 30
		marginSide = 20
	}
	
	// Calculate usable area
	usableWidth := width - (marginSide * 2)
	usableHeight := height - marginTop - marginBottom
	topY := marginTop
	bottomY := height - marginBottom
	leftX := marginSide
	rightX := width - marginSide
	middleY := marginTop + (usableHeight / 2)
	
	// Column positions
	col1 := leftX + (usableWidth / 4)
	col3 := rightX - (usableWidth / 4)
	
	// Row positions for cards with many symbols
	row1 := topY + (usableHeight / 4)
	row2 := middleY
	row3 := bottomY - (usableHeight / 4)
	
	switch rank {
	case 1: // Ace - one large center symbol
		return []image.Point{{centerX, centerY}}
		
	case 2:
		return []image.Point{
			{centerX, row1},
			{centerX, row3},
		}
		
	case 3:
		return []image.Point{
			{centerX, row1},
			{centerX, row2},
			{centerX, row3},
		}
		
	case 4:
		return []image.Point{
			{col1, row1},
			{col3, row1},
			{col1, row3},
			{col3, row3},
		}
		
	case 5:
		return []image.Point{
			{col1, row1},
			{col3, row1},
			{centerX, row2},
			{col1, row3},
			{col3, row3},
		}
		
	case 6:
		return []image.Point{
			{col1, row1},
			{col3, row1},
			{col1, row2},
			{col3, row2},
			{col1, row3},
			{col3, row3},
		}
		
	case 7:
		return []image.Point{
			{col1, row1},
			{col3, row1},
			{centerX, topY + (usableHeight / 3)}, // Special position for 7th symbol
			{col1, row2},
			{col3, row2},
			{col1, row3},
			{col3, row3},
		}
		
	case 8:
		return []image.Point{
			{col1, row1},
			{col3, row1},
			{centerX, topY + (usableHeight / 3)},
			{col1, row2},
			{col3, row2},
			{centerX, bottomY - (usableHeight / 3)},
			{col1, row3},
			{col3, row3},
		}
		
	case 9:
		// 9 needs special handling - 4 in corners, 4 in middle positions, 1 center
		midRow1 := topY + (usableHeight * 2 / 5)
		midRow2 := bottomY - (usableHeight * 2 / 5)
		return []image.Point{
			{col1, row1},
			{col3, row1},
			{col1, midRow1},
			{col3, midRow1},
			{centerX, row2},
			{col1, midRow2},
			{col3, midRow2},
			{col1, row3},
			{col3, row3},
		}
		
	case 10:
		// 10 needs special handling - similar to 9 but with 2 center symbols
		midRow1 := topY + (usableHeight * 2 / 5)
		midRow2 := bottomY - (usableHeight * 2 / 5)
		centerRow1 := row2 - (usableHeight / 10)
		centerRow2 := row2 + (usableHeight / 10)
		return []image.Point{
			{col1, row1},
			{col3, row1},
			{col1, midRow1},
			{col3, midRow1},
			{centerX, centerRow1},
			{centerX, centerRow2},
			{col1, midRow2},
			{col3, midRow2},
			{col1, row3},
			{col3, row3},
		}
		
	default: // Face cards (11, 12, 13) - no suit symbols in center
		return []image.Point{}
	}
}

func drawBackPattern(img *image.RGBA, size string) {
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	
	// Create a diamond pattern
	patternColor := color.RGBA{30, 30, 150, 255}
	lightPatternColor := color.RGBA{50, 50, 170, 255}
	
	// Draw diagonal lines to create diamond pattern
	spacing := 10
	if size == "icon" {
		spacing = 5
	} else if size == "large" {
		spacing = 20
	}
	
	// Create a more intricate pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Skip border area
			if x < 2 || x >= width-2 || y < 2 || y >= height-2 {
				continue
			}
			
			// Create diamond pattern
			if ((x/spacing)+(y/spacing))%2 == 0 {
				img.Set(x, y, patternColor)
			} else if ((x/spacing)%2 == 0 || (y/spacing)%2 == 0) {
				img.Set(x, y, lightPatternColor)
			}
		}
	}
}