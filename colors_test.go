package gift

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"testing"
)

func TestLut(t *testing.T) {
	fn := func(v float32) float32 {
		return v
	}
	for _, size := range []int{10, 100, 1000} {
		lut := prepareLut(size, fn)
		l := len(lut)
		if l != size {
			t.Errorf("LUT bad size: expected %v got %v", size, l)
		}
		if lut[0] != 0 {
			t.Errorf("LUT bad start value: expected 0 got %v", lut[0])
		}
		if lut[l-1] != 1 {
			t.Errorf("LUT bad end value: expected 1 got %v", lut[l-1])
		}
	}
	lut := prepareLut(10000, fn)
	for _, u := range []float32{0.0, 0.0001, 0.5555, 0.9999, 1.0} {
		v := getFromLut(lut, u)
		if math.Abs(float64(v-u)) > 0.0001 {
			t.Errorf("LUT bad value: expected %v got %v", u, v)
		}
	}
}

func TestInvert(t *testing.T) {
	src := image.NewGray(image.Rect(0, 0, 256, 1))
	for i := 0; i <= 255; i++ {
		src.Pix[i] = uint8(i)
	}
	g := New(Invert())
	dst := image.NewGray(g.Bounds(src.Bounds()))
	g.Draw(dst, src)

	for i := 0; i <= 255; i++ {
		if dst.Pix[i] != 255-src.Pix[i] {
			t.Errorf("InvertColors: index %d: expected %d got %d", i, 255-src.Pix[i], dst.Pix[i])
		}
	}
}

func TestColorspaceSRGBToLinear(t *testing.T) {
	vals := []float32{
		0.00000,
		0.01002,
		0.03310,
		0.07324,
		0.13287,
		0.21404,
		0.31855,
		0.44799,
		0.60383,
		0.78741,
		1.00000,
	}

	imgs := []draw.Image{
		image.NewGray(image.Rect(0, 0, 11, 11)),
		image.NewGray(image.Rect(0, 0, 111, 111)),
		image.NewGray16(image.Rect(0, 0, 11, 11)),
		image.NewGray16(image.Rect(0, 0, 1111, 1111)),
	}
	for _, img := range imgs {
		for i := 0; i <= 10; i++ {
			img.Set(i, 0, color.Gray{uint8(255 * float32(i) / 10.0)})
		}
		img2 := image.NewGray(img.Bounds())
		New(ColorspaceSRGBToLinear()).Draw(img2, img)
		if !img2.Bounds().Size().Eq(img.Bounds().Size()) {
			t.Errorf("ColorspaceSRGBToLinear bad result size: expected %v got %v", img.Bounds().Size(), img2.Bounds().Size())
		}
		for i := 0; i <= 10; i++ {
			expected := uint8(vals[i]*255.0 + 0.5)
			c := img2.At(i, 0).(color.Gray)
			if math.Abs(float64(c.Y)-float64(expected)) > 1 {
				t.Errorf("ColorspaceSRGBToLinear bad color value at index %v expected %v got %v", i, expected, c.Y)
			}
		}
	}
}

func TestColorspaceLinearToSRGB(t *testing.T) {
	vals := []float32{
		0.00000,
		0.34919,
		0.48453,
		0.58383,
		0.66519,
		0.73536,
		0.79774,
		0.85431,
		0.90633,
		0.95469,
		1.00000,
	}

	imgs := []draw.Image{
		image.NewGray(image.Rect(0, 0, 11, 11)),
		image.NewGray(image.Rect(0, 0, 111, 111)),
		image.NewGray16(image.Rect(0, 0, 11, 11)),
		image.NewGray16(image.Rect(0, 0, 1111, 1111)),
	}
	for _, img := range imgs {
		for i := 0; i <= 10; i++ {
			img.Set(i, 0, color.Gray{uint8(255 * float32(i) / 10.0)})
		}
		img2 := image.NewGray(img.Bounds())
		New(ColorspaceLinearToSRGB()).Draw(img2, img)
		if !img2.Bounds().Size().Eq(img.Bounds().Size()) {
			t.Errorf("ColorspaceLinearRGBToSRGB bad result size: expected %v got %v", img.Bounds().Size(), img2.Bounds().Size())
		}
		for i := 0; i <= 10; i++ {
			expected := uint8(vals[i]*255.0 + 0.5)
			c := img2.At(i, 0).(color.Gray)
			if math.Abs(float64(c.Y)-float64(expected)) > 1 {
				t.Errorf("ColorspaceLinearRGBToSRGB bad color value at index %v expected %v got %v", i, expected, c.Y)
			}
		}
	}

}

func TestAdjustGamma(t *testing.T) {
	src := image.NewGray(image.Rect(0, 0, 256, 1))
	dst := image.NewGray(image.Rect(0, 0, 256, 1))
	for i := 0; i <= 255; i++ {
		src.Pix[i] = uint8(i)
	}
	ag := Gamma(2.0)
	ag.Draw(dst, src, nil)

	for i := 100; i <= 150; i++ {
		if dst.Pix[i] <= src.Pix[i] {
			t.Errorf("Gamma unexpected color")
		}
	}

	ag = Gamma(0.5)
	ag.Draw(dst, src, nil)

	for i := 100; i <= 150; i++ {
		if dst.Pix[i] >= src.Pix[i] {
			t.Errorf("Gamma unexpected color")
		}
	}

	ag = Gamma(1.0)
	ag.Draw(dst, src, nil)

	for i := 100; i <= 150; i++ {
		if dst.Pix[i] != src.Pix[i] {
			t.Errorf("Gamma unexpected color")
		}
	}
}

func TestContrast(t *testing.T) {
	testData := []struct {
		desc           string
		p              float32
		srcb, dstb     image.Rectangle
		srcPix, dstPix []uint8
	}{
		{
			"contrast (0)",
			0,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
		},
		{
			"contrast (30)",
			30,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x00, 0x25, 0x00, 0x25, 0x00,
				0x53, 0xC5, 0xAE, 0xC5, 0x53,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
		},
		{
			"contrast (-30)",
			-30,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x26, 0x53, 0x26, 0x53, 0x26,
				0x69, 0xA1, 0x96, 0xA1, 0x69,
				0x26, 0x80, 0x26, 0x80, 0x26,
			},
		},
		{
			"contrast (100)",
			100,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0xFF, 0xFF, 0xFF, 0x00,
				0x00, 0xFF, 0x00, 0xFF, 0x00,
			},
		},
		{
			"contrast (200)",
			200,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0xFF, 0xFF, 0xFF, 0x00,
				0x00, 0xFF, 0x00, 0xFF, 0x00,
			},
		},
		{
			"contrast (-100)",
			-100,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x80, 0x80, 0x80, 0x80, 0x80,
				0x80, 0x80, 0x80, 0x80, 0x80,
				0x80, 0x80, 0x80, 0x80, 0x80,
			},
		},
		{
			"contrast (-200)",
			-200,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x80, 0x80, 0x80, 0x80, 0x80,
				0x80, 0x80, 0x80, 0x80, 0x80,
				0x80, 0x80, 0x80, 0x80, 0x80,
			},
		},
	}

	for _, d := range testData {
		src := image.NewGray(d.srcb)
		src.Pix = d.srcPix

		f := Contrast(d.p)
		dst := image.NewGray(f.Bounds(src.Bounds()))
		f.Draw(dst, src, nil)

		if !checkBoundsAndPix(dst.Bounds(), d.dstb, dst.Pix, d.dstPix) {
			t.Errorf("test [%s] failed: %#v, %#v", d.desc, dst.Bounds(), dst.Pix)
		}
	}
}

func TestBrightness(t *testing.T) {
	testData := []struct {
		desc           string
		p              float32
		srcb, dstb     image.Rectangle
		srcPix, dstPix []uint8
	}{
		{
			"brightness (0)",
			0,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
		},
		{
			"brightness (30)",
			30,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x4D, 0x8D, 0x4D, 0x8D, 0x4D,
				0xAD, 0xFD, 0xED, 0xFD, 0xAD,
				0x4D, 0xCD, 0x4D, 0xCD, 0x4D,
			},
		},
		{
			"brightness (-30)",
			-30,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x00, 0x00, 0x00, 0x00, 0x00,
				0x14, 0x64, 0x53, 0x64, 0x14,
				0x00, 0x34, 0x00, 0x34, 0x00,
			},
		},

		{
			"brightness (100)",
			100,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			},
		},
		{
			"brightness (200)",
			200,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			},
		},
		{
			"brightness (-100)",
			-100,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"brightness (-200)",
			-200,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
	}

	for _, d := range testData {
		src := image.NewGray(d.srcb)
		src.Pix = d.srcPix

		f := Brightness(d.p)
		dst := image.NewGray(f.Bounds(src.Bounds()))
		f.Draw(dst, src, nil)

		if !checkBoundsAndPix(dst.Bounds(), d.dstb, dst.Pix, d.dstPix) {
			t.Errorf("test [%s] failed: %#v, %#v", d.desc, dst.Bounds(), dst.Pix)
		}
	}
}

func TestSigmoid(t *testing.T) {
	testData := []struct {
		desc             string
		midpoint, factor float32
		srcb, dstb       image.Rectangle
		srcPix, dstPix   []uint8
	}{
		{
			"sigmoid (0.5, 0)",
			0.5, 0,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
		},
		{
			"sigmoid (0.5, 3)",
			0.5, 3,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x00, 0x38, 0x00, 0x38, 0x00,
				0x5B, 0xB7, 0xA5, 0xB7, 0x5B,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
		},
		{
			"sigmoid (0.5, -3)",
			0.5, -3,
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x00, 0x48, 0x00, 0x48, 0x00,
				0x65, 0xA9, 0x9B, 0xA9, 0x65,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
		},
	}

	for _, d := range testData {
		src := image.NewGray(d.srcb)
		src.Pix = d.srcPix

		f := Sigmoid(d.midpoint, d.factor)
		dst := image.NewGray(f.Bounds(src.Bounds()))
		f.Draw(dst, src, nil)

		if !checkBoundsAndPix(dst.Bounds(), d.dstb, dst.Pix, d.dstPix) {
			t.Errorf("test [%s] failed: %#v, %#v", d.desc, dst.Bounds(), dst.Pix)
		}
	}
}

func TestGrayscale(t *testing.T) {
	testData := []struct {
		desc           string
		srcb, dstb     image.Rectangle
		srcPix, dstPix []uint8
	}{

		{
			"grayscale 0x0",
			image.Rect(0, 0, 0, 0),
			image.Rect(0, 0, 0, 0),
			[]uint8{},
			[]uint8{},
		},
		{
			"grayscale 2x2",
			image.Rect(-1, -1, 1, 2),
			image.Rect(0, 0, 2, 3),
			[]uint8{
				0x00, 0x10, 0x20, 0x30, 0xFF, 0x00, 0x88, 0xFF,
				0xF0, 0xE0, 0xD0, 0xC0, 0x11, 0x66, 0xBB, 0x00,
				0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF,
			},
			[]uint8{
				0x0D, 0x0D, 0x0D, 0x30, 0x5C, 0x5C, 0x5C, 0xFF,
				0xe3, 0xe3, 0xe3, 0xC0, 0x56, 0x56, 0x56, 0x00,
				0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF,
			},
		},
	}

	for _, d := range testData {
		src := image.NewNRGBA(d.srcb)
		src.Pix = d.srcPix

		f := Grayscale()
		dst := image.NewNRGBA(f.Bounds(src.Bounds()))
		f.Draw(dst, src, nil)

		if !checkBoundsAndPix(dst.Bounds(), d.dstb, dst.Pix, d.dstPix) {
			t.Errorf("test [%s] failed: %#v, %#v", d.desc, dst.Bounds(), dst.Pix)
		}
	}
}

func TestSepia(t *testing.T) {
	testData := []struct {
		desc           string
		srcb, dstb     image.Rectangle
		srcPix, dstPix []uint8
	}{

		{
			"sepia 0x0",
			image.Rect(0, 0, 0, 0),
			image.Rect(0, 0, 0, 0),
			[]uint8{},
			[]uint8{},
		},
		{
			"sepia 2x2",
			image.Rect(-1, -1, 1, 2),
			image.Rect(0, 0, 2, 3),
			[]uint8{
				0x00, 0x10, 0x20, 0x30, 0xFF, 0x00, 0x88, 0xFF,
				0xF0, 0xE0, 0xD0, 0xC0, 0x11, 0x66, 0xBB, 0x00,
				0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF,
			},
			[]uint8{
				0x12, 0x10, 0x0D, 0x30, 0x7E, 0x70, 0x57, 0xFF,
				0xFF, 0xFF, 0xD4, 0xC0, 0x78, 0x6B, 0x54, 0x00,
				0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xEF, 0xFF,
			},
		},
	}

	for _, d := range testData {
		src := image.NewNRGBA(d.srcb)
		src.Pix = d.srcPix

		f := Sepia()
		dst := image.NewNRGBA(f.Bounds(src.Bounds()))
		f.Draw(dst, src, nil)

		if !checkBoundsAndPix(dst.Bounds(), d.dstb, dst.Pix, d.dstPix) {
			t.Errorf("test [%s] failed: %#v, %#v", d.desc, dst.Bounds(), dst.Pix)
		}
	}
}
