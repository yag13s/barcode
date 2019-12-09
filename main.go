package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"image/color"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code39"
	"github.com/go-pp/pp"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	rs6Letters       = "0123456789ABCD"
	rs6LetterIdxBits = 6
	rs6LetterIdxMask = 1<<rs6LetterIdxBits - 1
	rs6LetterIdxMax  = 63 / rs6LetterIdxBits
	margin           = 40
)

func main() {

	// A4 : 595pt × 842pt
	cambus := image.NewRGBA(image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: 595,
			Y: 842,
		},
	})

	var imgFile *os.File
	var textFile *os.File

	// create the output imgFile

	name := getFilename()

	imgFile, err := os.Create(name + ".png")
	defer imgFile.Close()

	if err != nil {
		pp.Println(err)
	}

	textFile, err = os.Create(name + ".txt")
	defer textFile.Close()

	if err != nil {
		pp.Println(err)
	}

	barcodeNumbers := []string{}

	for i := 0; i < 20; i++ {

		value := randString(15)

		fmt.Println("*" + value)
		barcodeNumbers = append(barcodeNumbers, fmt.Sprintln("*"+value))

		// Create the barcode
		c, err := code39.Encode(value, true, false)

		if err != nil {
			pp.Println(err)
		}

		bcd, err := barcode.Scale(c, c.Bounds().Dx(), 30)

		if err != nil {
			pp.Println(err)
		}

		selection := i % 2

		if selection == 0 {
			draw.Copy(cambus, image.Pt(selection*margin+30, margin*i+margin), bcd, bcd.Bounds(), draw.Over, nil)
			addLabel(cambus, selection*margin+30, margin*i+margin+40, value)
		} else {
			draw.Copy(cambus, image.Pt(selection*bcd.Bounds().Bounds().Dx()+margin+30, margin*(i-1)+margin), bcd, bcd.Bounds(), draw.Over, nil)
			addLabel(cambus, selection*bcd.Bounds().Bounds().Dx()+margin+30, margin*(i-1)+margin+40, value)
		}
	}
	png.Encode(imgFile, cambus)
	textFile.Write(([]byte)(strings.Join(barcodeNumbers, "")))
}

// `math/rand` パッケージの `lockedSource` を引用
// https://golang.org/src/math/rand/rand.go?s=12478:12833
type lockedSource struct {
	lk  sync.Mutex
	src rand.Source64
}

func (r *lockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Uint64() (n uint64) {
	r.lk.Lock()
	n = r.src.Uint64()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}

var rng = rand.New(&lockedSource{
	src: rand.NewSource(time.Now().UnixNano()).(rand.Source64),
})

func randString(n int) string {
	b := make([]byte, n)
	cache, remain := rng.Int63(), rs6LetterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = rng.Int63(), rs6LetterIdxMax
		}
		idx := int(cache & rs6LetterIdxMask)
		if idx < len(rs6Letters) {
			b[i] = rs6Letters[idx]
			i--
		}
		cache >>= rs6LetterIdxBits
		remain--
	}
	return string(b)
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{150, 150, 150, 255}
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func getFilename() string {
	u := uuid.NewV4()
	return fmt.Sprintf("created/barcode-%s", u.String())
}
