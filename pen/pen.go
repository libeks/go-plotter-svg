package pen

type Pen struct {
	Name    string
	Spacing float64
	XOffset float64
	YOffset float64
}

var (
	Micron005             = Pen{"Micron 005", 6, 0, 0}
	Micron01              = Pen{"Micron 01", 6, 3, -4}
	PilotG207             = Pen{"Pilot G-2 07", 10, 10, -50}
	Micron10              = Pen{"Micron 10", 20, -5, 5}
	TonborABTProThin      = Pen{"Tonbor ABT Pro Thin", 45, -20, -150}
	TonborABTProThick     = Pen{"Tonbor ABT Pro Thick", 45, 15, 10}
	SharpieHighliter      = Pen{"Sharpie Highliter", 20, 0, -50}
	SharpieCreativeMarker = Pen{"Sharpie Creative Marker", 30, -10, -35} // displacement can change depending on positioning
	UniballEcoJapanPen    = Pen{"Uni-ball eco 'Japan' pen", 10, -5, 20}
	WexfordGelInkPen      = Pen{"Wexford Gel Ink Pen", 10, 10, -30}
	BicBU3Grip            = Pen{"Bic BU3 Grip", 7, 25, -110}
	SharpiePen            = Pen{"Sharpie Pen", 10, 5, -25}
	BicIntensityFineTip   = Pen{"Bic Intensity Fine Tip", 15, -10, 25}  // displacement can vary
	BicIntensityBrushTip  = Pen{"Bic Intensity Brush Tip", 30, -5, -50} // displacement can vary
)
