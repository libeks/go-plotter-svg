package scenes

import (
	"github.com/libeks/go-plotter-svg/foldable"
	"github.com/libeks/go-plotter-svg/primitives"
)

func foldableCubeIDScene(b primitives.BBox) Document {
	foldableBase := 1500.0
	patterns := foldable.Cube(b, foldableBase)
	scene := FromFoldableLayers(patterns, b)
	return scene
}

func foldableRhombicuboctahedronIDScene(b primitives.BBox) Document {
	foldableBase := 1500.0
	patterns := foldable.Rhombicuboctahedron(b, foldableBase)
	scene := FromFoldableLayers(patterns, b)
	return scene
}

func foldableRhombicuboctahedronSansCornersScene(b primitives.BBox) Document {
	foldableBase := 1500.0
	patterns := foldable.RhombicuboctahedronWithoutCorners(b, foldableBase)
	scene := FromFoldableLayers(patterns, b)
	return scene
}

func foldableRhombicuboctahedronSansCornersTricolorScene(b primitives.BBox) Document {
	foldableBase := 1500.0
	patterns := foldable.RhombicuboctahedronWithoutCornersTricolor(b, foldableBase)
	scene := FromFoldableLayers(patterns, b)
	return scene
}

func foldableRightTrianglePrismIDScene(b primitives.BBox) Document {
	foldableBase := 1500.0
	patterns := foldable.RightTrianglePrism(b, foldableBase, foldableBase, foldableBase)
	scene := FromFoldableLayers(patterns, b)
	return scene
}

func foldableCutCornerScene(b primitives.BBox) Document {
	foldableBase := 1500.0
	patterns := foldable.CutCube(b, foldableBase, 0.5)
	scene := FromFoldableLayers(patterns, b)
	return scene
}

func foldableVoronoiScene(b primitives.BBox) Document {
	patterns := foldable.VoronoiFoldable(b)
	scene := FromFoldableLayers(patterns, b)
	return scene
}
