package foldable

type CutOut struct {
	//
	Faces       []Face
	Connections map[int]Connection
}

func NewCutOut(faces []Face, connections map[int]Connection) CutOut {
	return CutOut{
		Faces:       faces,
		Connections: connections,
	}
}
