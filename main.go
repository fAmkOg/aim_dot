package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

func run(yOffset int) error {
	// Initialize SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return err
	}
	defer sdl.Quit()

	numDisplays, err := sdl.GetNumVideoDisplays()
	if err != nil {
		return err
	}
	fmt.Printf("Number of displays: %d\n", numDisplays)

	// Print all displays' info
	for i := range numDisplays {
		displayMode, err := sdl.GetCurrentDisplayMode(i)
		if err != nil {
			return err
		}
		fmt.Printf("Display %d: Width=%d, Height=%d\n", i, displayMode.W, displayMode.H)
	}

	// Get screen dimensions
	displayIndex := 0

	/*
		displayMode, err := sdl.GetCurrentDisplayMode(displayIndex)
		if err != nil {
			return err
		}
		screenWidth := displayMode.W
		screenHeight := displayMode.H


		fmt.Printf("Using Display %d with Width=%d, Height=%d\n", displayIndex, screenWidth, screenHeight)
	*/

	displayBounds, err := sdl.GetDisplayBounds(displayIndex)
	if err != nil {
		return err
	}
	fmt.Printf("Display %d Bounds: X=%d, Y=%d, Width=%d, Height=%d\n", displayIndex, displayBounds.X, displayBounds.Y, displayBounds.W, displayBounds.H)
	centerX := displayBounds.X + (displayBounds.W / 2)
	centerY := displayBounds.Y + (displayBounds.H / 2) + int32(yOffset)

	// Create a borderless, always-on-top window
	windowSize := int32(5)
	window, err := sdl.CreateWindow(
		"Aim Dot",
		// int32(screenWidth/2-windowSize/2), int32(screenHeight/2-windowSize/2), // Centered
		int32(centerX-windowSize/2), int32(centerY-windowSize/2), // Centered
		// 50, 50, // Small size
		windowSize, windowSize,
		sdl.WINDOW_ALWAYS_ON_TOP|sdl.WINDOW_BORDERLESS|sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return err
	}
	defer window.Destroy()

	// Create a renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return err
	}
	defer renderer.Destroy()

	// Set window transparency (Linux requires a compositor to support transparency)
	if runtime.GOOS == "linux" {
		// sdl.SetWindowOpacity(window, 0.5) // Adjust if needed
		// window.SetWindowOpacity(0.5)
	}

	// Main loop
	running := true
	for running {
		// Handle events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		// Clear screen (transparent)
		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.Clear()

		// Draw red dot
		renderer.SetDrawColor(255, 0, 0, 255) // Red
		renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: windowSize, H: windowSize})

		// Present renderer
		renderer.Present()

		// Small delay to reduce CPU usage
		sdl.Delay(16)
	}

	return nil
}

func main() {
	// flag.Parse()
	yOffset := 0
	if len(os.Args) > 1 {
		val, err := strconv.Atoi(os.Args[1])
		if err == nil {
			yOffset = val
		} else {
			fmt.Printf("Invalid argument: %s\n", os.Args[1])
		}
	}
	fmt.Printf("Y Offset: %d\n", yOffset)

	if err := run(yOffset); err != nil {
		log.Fatal(err)
	}
}
