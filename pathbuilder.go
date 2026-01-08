package svg

import (
	"fmt"
	"strings"
)

// PathBuilder provides a fluent API for constructing SVG path data
type PathBuilder struct {
	commands strings.Builder
}

// Point represents a 2D point
type Point struct {
	X, Y float64
}

// NewPathBuilder creates a new path builder
func NewPathBuilder() *PathBuilder {
	return &PathBuilder{}
}

// MoveTo moves the pen to the specified point without drawing
func (pb *PathBuilder) MoveTo(x, y float64) *PathBuilder {
	fmt.Fprintf(&pb.commands, "M %.2f %.2f ", x, y)
	return pb
}

// LineTo draws a line from the current point to the specified point
func (pb *PathBuilder) LineTo(x, y float64) *PathBuilder {
	fmt.Fprintf(&pb.commands, "L %.2f %.2f ", x, y)
	return pb
}

// HorizontalLineTo draws a horizontal line to the specified x coordinate
func (pb *PathBuilder) HorizontalLineTo(x float64) *PathBuilder {
	fmt.Fprintf(&pb.commands, "H %.2f ", x)
	return pb
}

// VerticalLineTo draws a vertical line to the specified y coordinate
func (pb *PathBuilder) VerticalLineTo(y float64) *PathBuilder {
	fmt.Fprintf(&pb.commands, "V %.2f ", y)
	return pb
}

// CurveTo draws a cubic Bézier curve
func (pb *PathBuilder) CurveTo(x1, y1, x2, y2, x, y float64) *PathBuilder {
	fmt.Fprintf(&pb.commands, "C %.2f %.2f, %.2f %.2f, %.2f %.2f ", x1, y1, x2, y2, x, y)
	return pb
}

// SmoothCurveTo draws a smooth cubic Bézier curve (first control point is reflection of previous)
func (pb *PathBuilder) SmoothCurveTo(x2, y2, x, y float64) *PathBuilder {
	fmt.Fprintf(&pb.commands, "S %.2f %.2f, %.2f %.2f ", x2, y2, x, y)
	return pb
}

// QuadraticCurveTo draws a quadratic Bézier curve
func (pb *PathBuilder) QuadraticCurveTo(x1, y1, x, y float64) *PathBuilder {
	fmt.Fprintf(&pb.commands, "Q %.2f %.2f, %.2f %.2f ", x1, y1, x, y)
	return pb
}

// SmoothQuadraticCurveTo draws a smooth quadratic Bézier curve
func (pb *PathBuilder) SmoothQuadraticCurveTo(x, y float64) *PathBuilder {
	fmt.Fprintf(&pb.commands, "T %.2f %.2f ", x, y)
	return pb
}

// ArcTo draws an elliptical arc
// rx, ry: x and y radius
// xAxisRotation: rotation of the ellipse in degrees
// largeArcFlag: 0 for small arc, 1 for large arc
// sweepFlag: 0 for counter-clockwise, 1 for clockwise
// x, y: end point
func (pb *PathBuilder) ArcTo(rx, ry, xAxisRotation float64, largeArcFlag, sweepFlag int, x, y float64) *PathBuilder {
	fmt.Fprintf(&pb.commands, "A %.2f %.2f %.2f %d %d %.2f %.2f ", rx, ry, xAxisRotation, largeArcFlag, sweepFlag, x, y)
	return pb
}

// Close closes the current path by drawing a line back to the first point
func (pb *PathBuilder) Close() *PathBuilder {
	pb.commands.WriteString("Z ")
	return pb
}

// String returns the path data string
func (pb *PathBuilder) String() string {
	return strings.TrimSpace(pb.commands.String())
}

// Build returns the path data string (alias for String)
func (pb *PathBuilder) Build() string {
	return pb.String()
}

// Reset clears the path builder for reuse
func (pb *PathBuilder) Reset() *PathBuilder {
	pb.commands.Reset()
	return pb
}

// Helper functions for common path patterns

// RectPath creates a rectangular path
func RectPath(x, y, width, height float64) string {
	return NewPathBuilder().
		MoveTo(x, y).
		HorizontalLineTo(x + width).
		VerticalLineTo(y + height).
		HorizontalLineTo(x).
		Close().
		String()
}

// RoundedRectPath creates a rounded rectangular path
func RoundedRectPath(x, y, width, height, rx, ry float64) string {
	if ry == 0 {
		ry = rx
	}

	pb := NewPathBuilder()
	pb.MoveTo(x+rx, y)
	pb.HorizontalLineTo(x + width - rx)
	pb.ArcTo(rx, ry, 0, 0, 1, x+width, y+ry)
	pb.VerticalLineTo(y + height - ry)
	pb.ArcTo(rx, ry, 0, 0, 1, x+width-rx, y+height)
	pb.HorizontalLineTo(x + rx)
	pb.ArcTo(rx, ry, 0, 0, 1, x, y+height-ry)
	pb.VerticalLineTo(y + ry)
	pb.ArcTo(rx, ry, 0, 0, 1, x+rx, y)
	pb.Close()

	return pb.String()
}

// CirclePath creates a circular path using arcs
func CirclePath(cx, cy, r float64) string {
	return NewPathBuilder().
		MoveTo(cx-r, cy).
		ArcTo(r, r, 0, 0, 1, cx+r, cy).
		ArcTo(r, r, 0, 0, 1, cx-r, cy).
		Close().
		String()
}

// EllipsePath creates an elliptical path using arcs
func EllipsePath(cx, cy, rx, ry float64) string {
	return NewPathBuilder().
		MoveTo(cx-rx, cy).
		ArcTo(rx, ry, 0, 0, 1, cx+rx, cy).
		ArcTo(rx, ry, 0, 0, 1, cx-rx, cy).
		Close().
		String()
}

// PolylinePath creates a path from multiple points (open path)
func PolylinePath(points []Point) string {
	if len(points) == 0 {
		return ""
	}

	pb := NewPathBuilder()
	pb.MoveTo(points[0].X, points[0].Y)

	for i := 1; i < len(points); i++ {
		pb.LineTo(points[i].X, points[i].Y)
	}

	return pb.String()
}

// PolygonPath creates a closed path from multiple points
func PolygonPath(points []Point) string {
	if len(points) == 0 {
		return ""
	}

	pb := NewPathBuilder()
	pb.MoveTo(points[0].X, points[0].Y)

	for i := 1; i < len(points); i++ {
		pb.LineTo(points[i].X, points[i].Y)
	}

	pb.Close()
	return pb.String()
}

// SmoothLinePath creates a smooth curve through points using cubic Bézier curves
// tension controls how tight the curve is (0 = straight lines, 1 = very curved)
func SmoothLinePath(points []Point, tension float64) string {
	if len(points) < 2 {
		return ""
	}
	if len(points) == 2 {
		return PolylinePath(points)
	}

	pb := NewPathBuilder()
	pb.MoveTo(points[0].X, points[0].Y)

	// Calculate control points for cubic Bézier curves
	for i := 0; i < len(points)-1; i++ {
		var cp1x, cp1y, cp2x, cp2y float64

		if i == 0 {
			// First point
			cp1x = points[i].X + (points[i+1].X-points[i].X)*tension
			cp1y = points[i].Y + (points[i+1].Y-points[i].Y)*tension
		} else {
			// Calculate control point based on previous point
			cp1x = points[i].X + (points[i+1].X-points[i-1].X)*tension
			cp1y = points[i].Y + (points[i+1].Y-points[i-1].Y)*tension
		}

		if i == len(points)-2 {
			// Last segment
			cp2x = points[i+1].X - (points[i+1].X-points[i].X)*tension
			cp2y = points[i+1].Y - (points[i+1].Y-points[i].Y)*tension
		} else {
			// Calculate control point based on next point
			cp2x = points[i+1].X - (points[i+2].X-points[i].X)*tension
			cp2y = points[i+1].Y - (points[i+2].Y-points[i].Y)*tension
		}

		pb.CurveTo(cp1x, cp1y, cp2x, cp2y, points[i+1].X, points[i+1].Y)
	}

	return pb.String()
}

// AreaPath creates a filled area path from points with a baseline
func AreaPath(points []Point, baselineY float64) string {
	if len(points) == 0 {
		return ""
	}

	pb := NewPathBuilder()
	pb.MoveTo(points[0].X, baselineY)
	pb.LineTo(points[0].X, points[0].Y)

	for i := 1; i < len(points); i++ {
		pb.LineTo(points[i].X, points[i].Y)
	}

	pb.LineTo(points[len(points)-1].X, baselineY)
	pb.Close()

	return pb.String()
}

// SmoothAreaPath creates a smooth filled area path from points with a baseline
func SmoothAreaPath(points []Point, baselineY float64, tension float64) string {
	if len(points) == 0 {
		return ""
	}

	pb := NewPathBuilder()
	pb.MoveTo(points[0].X, baselineY)
	pb.LineTo(points[0].X, points[0].Y)

	// Draw smooth curve through top points
	for i := 0; i < len(points)-1; i++ {
		var cp1x, cp1y, cp2x, cp2y float64

		if i == 0 {
			cp1x = points[i].X + (points[i+1].X-points[i].X)*tension
			cp1y = points[i].Y + (points[i+1].Y-points[i].Y)*tension
		} else {
			cp1x = points[i].X + (points[i+1].X-points[i-1].X)*tension
			cp1y = points[i].Y + (points[i+1].Y-points[i-1].Y)*tension
		}

		if i == len(points)-2 {
			cp2x = points[i+1].X - (points[i+1].X-points[i].X)*tension
			cp2y = points[i+1].Y - (points[i+1].Y-points[i].Y)*tension
		} else {
			cp2x = points[i+1].X - (points[i+2].X-points[i].X)*tension
			cp2y = points[i+1].Y - (points[i+2].Y-points[i].Y)*tension
		}

		pb.CurveTo(cp1x, cp1y, cp2x, cp2y, points[i+1].X, points[i+1].Y)
	}

	pb.LineTo(points[len(points)-1].X, baselineY)
	pb.Close()

	return pb.String()
}
