package main

import (
	"errors"
	"image/color"
	"strconv"

	"bitbucket.org/oakmoundstudio/oak"
	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
	pt "bitbucket.org/oakmoundstudio/oak/render/particle"
	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
)

var (
	startColor     color.Color
	startColorRand color.Color
	endColor       color.Color
	endColorRand   color.Color
	src            *pt.Source
)

func main() {
	oak.AddCommand("shape", func(args []string) {
		if len(args) > 1 {
			switch args[1] {
			case "square":
				src.Generator.(pt.Shapeable).SetShape(pt.Square)
			case "circle":
				src.Generator.(pt.Shapeable).SetShape(pt.Circle)
			case "diamond":
				src.Generator.(pt.Shapeable).SetShape(pt.Diamond)
			}
		}
	})

	oak.AddCommand("size", func(args []string) {
		f1, f2, two, err := parseInts(args)
		if err != nil {
			return
		}
		if !two {
			src.Generator.(pt.Sizeable).SetSize(intrange.NewConstant(f1))
		} else {
			src.Generator.(pt.Sizeable).SetSize(intrange.NewLinear(f1, f2))
		}
	})

	oak.AddCommand("count", func(args []string) {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return
		}
		if !two {
			src.Generator.GetBaseGenerator().NewPerFrame = floatrange.NewConstant(npf)
		} else {
			src.Generator.GetBaseGenerator().NewPerFrame = floatrange.NewLinear(npf, npf2)
		}
	})

	oak.AddCommand("life", func(args []string) {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return
		}
		if !two {
			src.Generator.GetBaseGenerator().LifeSpan = floatrange.NewConstant(npf)
		} else {
			src.Generator.GetBaseGenerator().LifeSpan = floatrange.NewLinear(npf, npf2)
		}
	})

	oak.AddCommand("rotation", func(args []string) {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return
		}
		if !two {
			src.Generator.GetBaseGenerator().Rotation = floatrange.NewConstant(npf)
		} else {
			src.Generator.GetBaseGenerator().Rotation = floatrange.NewLinear(npf, npf2)
		}
	})

	oak.AddCommand("angle", func(args []string) {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return
		}
		if !two {
			src.Generator.GetBaseGenerator().Angle = floatrange.NewConstant(npf)
		} else {
			src.Generator.GetBaseGenerator().Angle = floatrange.NewLinear(npf, npf2)
		}
	})

	oak.AddCommand("speed", func(args []string) {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return
		}
		if !two {
			src.Generator.GetBaseGenerator().Speed = floatrange.NewConstant(npf)
		} else {
			src.Generator.GetBaseGenerator().Speed = floatrange.NewLinear(npf, npf2)
		}
	})

	oak.AddCommand("spread", func(args []string) {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return
		}
		if !two {
			return
		}
		src.Generator.GetBaseGenerator().Spread.SetPos(npf, npf2)
	})

	oak.AddCommand("gravity", func(args []string) {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return
		}
		if !two {
			return
		}
		src.Generator.GetBaseGenerator().Gravity.SetPos(npf, npf2)
	})

	oak.AddCommand("speeddecay", func(args []string) {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return
		}
		if !two {
			return
		}
		src.Generator.GetBaseGenerator().SpeedDecay.SetPos(npf, npf2)
	})

	oak.AddCommand("pos", func(args []string) {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return
		}
		if !two {
			return
		}
		src.Generator.SetPos(npf, npf2)
	})

	oak.AddCommand("startcolor", func(args []string) {
		if len(args) > 4 {
			r, g, b, a, err := parseRGBA(args)
			if err != nil {
				return
			}
			startColor = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			src.Generator.(pt.Colorable).SetStartColor(startColor, startColorRand)
		}
	})

	oak.AddCommand("startrand", func(args []string) {
		if len(args) > 4 {
			r, g, b, a, err := parseRGBA(args)
			if err != nil {
				return
			}
			startColorRand = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			src.Generator.(pt.Colorable).SetStartColor(startColor, startColorRand)
		}
	})

	oak.AddCommand("endcolor", func(args []string) {
		if len(args) > 4 {
			r, g, b, a, err := parseRGBA(args)
			if err != nil {
				return
			}
			endColor = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			src.Generator.(pt.Colorable).SetEndColor(endColor, endColorRand)
		}
	})

	oak.AddCommand("endrand", func(args []string) {
		if len(args) > 4 {
			r, g, b, a, err := parseRGBA(args)
			if err != nil {
				return
			}
			endColorRand = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			src.Generator.(pt.Colorable).SetEndColor(endColor, endColorRand)
		}
	})

	oak.AddScene("demo", func(string, interface{}) {
		x := 320.0
		y := 240.0
		newPf := floatrange.NewLinear(1, 2)
		life := floatrange.NewLinear(100, 120)
		angle := floatrange.NewLinear(0, 360)
		speed := floatrange.NewLinear(1, 5)
		size := intrange.Constant(1)
		layerFn := func(v physics.Vector) int {
			return 1
		}
		startColor = color.RGBA{255, 255, 255, 255}
		startColorRand = color.RGBA{0, 0, 0, 0}
		endColor = color.RGBA{255, 255, 255, 255}
		endColorRand = color.RGBA{0, 0, 0, 0}
		shape := pt.Square
		src = pt.NewColorGenerator(
			pt.Pos(x, y),
			pt.Duration(pt.Inf),
			pt.LifeSpan(life),
			pt.Angle(angle),
			pt.Speed(speed),
			pt.Layer(layerFn),
			pt.Shape(shape),
			pt.Size(size),
			pt.Color(startColor, startColorRand, endColor, endColorRand),
			pt.NewPerFrame(newPf)).Generate(0)
	}, func() bool {
		return true
	}, func() (string, *oak.SceneResult) {
		return "demo", nil
	})

	render.SetDrawStack(
		render.NewCompositeR([]render.Renderable{}),
		render.NewDrawFPS(),
	)

	oak.Init("demo")
}

func parseRGBA(args []string) (r, g, b, a int, err error) {
	if len(args) < 5 {
		return
	}
	r, err = strconv.Atoi(args[1])
	if err != nil {
		return
	}
	g, err = strconv.Atoi(args[2])
	if err != nil {
		return
	}
	b, err = strconv.Atoi(args[3])
	if err != nil {
		return
	}
	a, err = strconv.Atoi(args[4])
	return
}

func parseFloats(args []string) (f1, f2 float64, two bool, err error) {
	if len(args) < 2 {
		err = errors.New("No args")
		return
	}
	f1, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return
	}
	if len(args) < 3 {
		return
	}
	f2, err = strconv.ParseFloat(args[2], 64)
	if err != nil {
		return
	}
	two = true
	return
}

func parseInts(args []string) (i1, i2 int, two bool, err error) {
	if len(args) < 2 {
		err = errors.New("No args")
		return
	}
	i1, err = strconv.Atoi(args[1])
	if err != nil {
		return
	}
	if len(args) < 3 {
		return
	}
	i2, err = strconv.Atoi(args[2])
	if err != nil {
		return
	}
	two = true
	return
}
