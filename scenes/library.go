package scenes

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/libeks/go-plotter-svg/primitives"
)

func SceneLibrary() sceneLibrary {
	return sceneLibrary{
		scenes: make(map[string]func(b primitives.BBox) Document),
	}
}

type sceneLibrary struct {
	// map keys are considered to be case insensitive
	scenes map[string]func(b primitives.BBox) Document
}

func (l *sceneLibrary) Add(name string, scene func(b primitives.BBox) Document) error {
	lowerName := strings.ToLower(name)
	if _, ok := l.scenes[lowerName]; ok {
		return errors.New(fmt.Sprintf("Scene with name '%s' already added", name))
	}
	l.scenes[lowerName] = scene
	return nil
}

func (l *sceneLibrary) Get(name string) (func(b primitives.BBox) Document, error) {
	lowerName := strings.ToLower(name)
	if scene, ok := l.scenes[lowerName]; ok {
		return scene, nil
	}
	return nil, errors.New(fmt.Sprintf("Couldn't find scene with name '%s'", name))
}

func (l *sceneLibrary) GetNames() []string {
	names := slices.Collect(maps.Keys(l.scenes))
	slices.Sort(names)
	return names
}
