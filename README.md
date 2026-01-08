# svg

A clean, efficient SVG rendering library that integrates with the [layout](https://github.com/SCKelemen/layout) engine to produce beautiful SVG graphics from layout trees.

## Features

- **Layout Integration**: Seamlessly renders `layout.Node` trees to SVG
- **Transform Support**: Full 2D transform support (translate, rotate, scale, skew)
- **Styling System**: Colors, borders, backgrounds, shadows
- **Text Rendering**: SVG text elements with proper positioning
- **ClipPath Management**: Thread-safe unique ID generation for clipping
- **Gradient Support**: Linear and radial gradients with multiple color spaces (OKLCH, OKLAB, sRGB, Display P3)
- **Design Tokens**: Themeable styling system

## Installation

```bash
go get github.com/SCKelemen/svg
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/SCKelemen/layout"
    "github.com/SCKelemen/svg"
)

func main() {
    // Create a layout tree
    root := &layout.Node{
        Style: layout.Style{
            Display: layout.DisplayFlex,
            FlexDirection: layout.FlexDirectionRow,
            Width: 400,
            Height: 200,
        },
        Children: []*layout.Node{
            {Style: layout.Style{Width: 100, Height: 100}},
            {Style: layout.Style{Width: 100, Height: 100}},
        },
    }

    // Perform layout
    constraints := layout.Loose(800, 600)
    layout.Layout(root, constraints, nil)

    // Render to SVG
    output := svg.RenderToSVG(root, svg.Options{
        Width: 400,
        Height: 200,
    })

    fmt.Println(output)
}
```

### With Styling

```go
// Create a styled renderer
renderer := svg.NewRenderer(svg.Options{
    Width: 800,
    Height: 600,
    StyleSheet: svg.DefaultStyles(),
})

// Render with custom styles
output := renderer.Render(root)
```

### Gradients

```go
// Create a gradient with multiple color spaces
gradient := svg.LinearGradient{
    ID: "myGradient",
    X1: 0, Y1: 0,
    X2: 100, Y2: 0,
    Stops: []svg.GradientStop{
        {Offset: "0%", Color: "#3B82F6"},
        {Offset: "100%", Color: "#8B5CF6"},
    },
    ColorSpace: color.GradientOKLCH, // Perceptually uniform gradients
}

// Apply gradient to elements
svgElement := fmt.Sprintf(`<rect fill="url(#myGradient)" x="0" y="0" width="100" height="50"/>`)
```

## Design Philosophy

This library focuses on:

1. **Simplicity**: Clean API with sensible defaults
2. **Integration**: First-class support for the layout engine
3. **Extensibility**: Easy to add custom rendering logic
4. **Performance**: Efficient string building and minimal allocations

## Related Projects

- [layout](https://github.com/SCKelemen/layout) - CSS Grid/Flexbox layout engine
- [text](https://github.com/SCKelemen/text) - Unicode text handling
- [color](https://github.com/SCKelemen/color) - Color space handling
- [cli](https://github.com/SCKelemen/cli) - Terminal rendering
