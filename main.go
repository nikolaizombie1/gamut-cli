package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/muesli/gamut"
)

type Color struct {
	Color string
}

func main() {
	var darkerFlag = flag.Float64("Darker", 0.0, "Make input color darker by a specific percentage.")
	var lighterFlag = flag.Float64("Lighter", 0.0, "Make input color lighter by a specific percentage.")
	var complementaryFlag = flag.Bool("Complementary", false, "Get the complementary color of the input color.")
	var contrastFlag = flag.Bool("Contrast", false, "Get the color with the highest contrast of the input color, either black or white.")
	var hueOffsetFlag = flag.Int("HueOffset", 0, "Change the angle of the input color without changing the lightness or saturation.")
	var triadicSchemeFlag = flag.Bool("Triadic", false, "A color scheme comprised of three equally spaced colors around the color wheel based on the input color.")
	var quadraticSchemeFlag = flag.Bool("Quadratic", false, "A color scheme comprised of four equally spaced colors around the color wheel based on the input color.")
	var tetraticSchemeFlag = flag.Bool("Tetratic", false, "A color scheme made up by two input colors and their complementary values.")
	var analogousSchemeFlag = flag.Bool("Analogous", false, "A color scheme created by getting the two colors that sit next to the input color on the color wheel.")
	var splitComplementaryFlag = flag.Bool("SplitComplementary", false, "A color scheme created by getting the two colors that sit next to the complement of the input color on the color wheel.")
	var warmFlag = flag.Bool("Warm", false, "Determine if the input color is warm.")
	var coolFlag = flag.Bool("Cool", false, "Determine if the input color is cool.")
	var monochromaticFlag = flag.Int("Monochromatic", 0, "Number of colors of the same hue, but with a different saturation/lightness based on the input color.")
	var shadesFlag = flag.Int("Shades", 0, "Number of colors, based on the input color, blended from the given color to black.")
	var tintsFlag = flag.Int("Tints", 0, "Number of colors, based on the input color, blended from the given color to white.")
	var tonesFlag = flag.Int("Tones", 0, "Number of colors, based on the input color, blended from the given color to gray.")
	var blendsFlag = flag.Int("Blends", 0, "Number of colors, based on the two input colors, interpolated together.")
	var color1Flag = flag.String("Color1", "", "Hex RGB value of the color1. Can start with and without #. Can be short or standard formats.")
	var color2Flag = flag.String("Color2", "", "Hex RGB value of the color2. Can start with and without #. Can be short or standard formats.")

	flag.Parse()

	if *color1Flag == "" {
		log.Fatal("Color1 flag was not specified. Aborting.")
	}

	var color1 = DecodeColor(*color1Flag)

	var retStr string

	if *darkerFlag > 0 {
		retStr = MarshallInput(Color{Color: gamut.ToHex(gamut.Darker(color1, *darkerFlag))})
	} else if *lighterFlag > 0 {
		retStr = MarshallInput(Color{Color: gamut.ToHex(gamut.Lighter(color1, *lighterFlag))})
	} else if *complementaryFlag {
		retStr = MarshallInput(Color{Color: gamut.ToHex(gamut.Complementary(color1))})
	} else if *contrastFlag {
		retStr = MarshallInput(Color{Color: gamut.ToHex(gamut.Contrast(color1))})
	} else if *hueOffsetFlag > 0 {
		retStr = MarshallInput(Color{Color: gamut.ToHex(gamut.HueOffset(color1, *hueOffsetFlag))})
	} else if *triadicSchemeFlag {
		retStr = MarshallInput(GetColors1Color(gamut.Triadic, color1))
	} else if *quadraticSchemeFlag {
		retStr = MarshallInput(GetColors1Color(gamut.Quadratic, color1))
	} else if *tetraticSchemeFlag {
		if *color2Flag == "" {
			log.Fatal("Color2 flag is not specified. Aborting.")
		}
		retStr = MarshallInput(GetColors2Color(gamut.Tetradic, color1, DecodeColor(*color2Flag)))
	} else if *analogousSchemeFlag {
		retStr = MarshallInput(GetColors1Color(gamut.Analogous, color1))
	} else if *splitComplementaryFlag {
		retStr = MarshallInput(GetColors1Color(gamut.SplitComplementary, color1))
	} else if *warmFlag {
		fmt.Println(gamut.Warm(color1))
		os.Exit(0)
	} else if *coolFlag {
		fmt.Println(gamut.Cool(color1))
		os.Exit(0)
	} else if *monochromaticFlag > 0 {
		retStr = MarshallInput(GetColors1ColorSTT(gamut.Monochromatic, color1, *monochromaticFlag))
	} else if *shadesFlag > 0 {
		retStr = MarshallInput(GetColors1ColorSTT(gamut.Shades, color1, *shadesFlag))
	} else if *tintsFlag > 0 {
		retStr = MarshallInput(GetColors1ColorSTT(gamut.Tints, color1, *tintsFlag))
	} else if *tonesFlag > 0 {
		retStr = MarshallInput(GetColors1ColorSTT(gamut.Tones, color1, *tonesFlag))
	} else if *blendsFlag > 0 {
		if *color2Flag == "" {
			log.Fatal("Color2 flag is not specified. Aborting.")
		}
		retStr = MarshallInput(GetColors2ColorSTT(gamut.Blends, color1, DecodeColor(*color2Flag), *blendsFlag))
	} else {
		log.Fatal("No operation flag specified. Aborting.")
	}

	fmt.Println(retStr)
}

func MarshallInput(a any) string {
	retStr, err := json.Marshal(a)
	if err != nil {
		log.Fatal("Failure converting to JSON.")
	}
	return string(retStr)

}

func DecodeColor(str string) color.Color {
	var color1Str string
	if str[0] != '#' {
		color1Str = fmt.Sprint("#", str)
	} else {
		color1Str = str
	}

	return gamut.Hex(color1Str)
}

func GetColors1Color(f func(color.Color) []color.Color, color color.Color) []Color {
	colors := f(color)
	colorsStrs := []Color{{Color: gamut.ToHex(color)}}
	for _, v := range colors {
		colorsStrs = append(colorsStrs, Color{Color: gamut.ToHex(v)})
	}
	return colorsStrs
}

func GetColors2Color(f func(color.Color, color.Color) []color.Color, color1 color.Color, color2 color.Color) []Color {
	colors := f(color1, color2)
	colorsStrs := []Color{{Color: gamut.ToHex(color1)}, {Color: gamut.ToHex(color2)}}
	for _, v := range colors {
		colorsStrs = append(colorsStrs, Color{Color: gamut.ToHex(v)})
	}
	return colorsStrs
}

func GetColors1ColorSTT(f func(color.Color, int) []color.Color, color color.Color, num int) []Color {
	colors := f(color, num)
	colorsStrs := []Color{{Color: gamut.ToHex(color)}}
	for _, v := range colors {
		colorsStrs = append(colorsStrs, Color{Color: gamut.ToHex(v)})
	}
	return colorsStrs
}

func GetColors2ColorSTT(f func(color.Color, color.Color, int) []color.Color, color1 color.Color, color2 color.Color, num int) []Color {
	colors := f(color1, color2, num)
	colorsStrs := []Color{{Color: gamut.ToHex(color1)}, {Color: gamut.ToHex(color2)}}
	for _, v := range colors {
		colorsStrs = append(colorsStrs, Color{Color: gamut.ToHex(v)})
	}
	return colorsStrs
}
