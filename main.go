package main

import (
	"bytes"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	_ "github.com/chai2010/webp"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"strings"
)

var selectAll = &fyne.ShortcutSelectAll{}
var w fyne.Window

func init() {

}

func main() {
	//os.Setenv("FYNE_FONT", "/System/Library/Fonts/PingFang.ttc")
	//defer os.Unsetenv("FYNE_FONT")
	myApp := app.New()
	w = myApp.NewWindow("URL Image Viewer")
	urlInput := widget.NewEntry()
	urlInput.SetPlaceHolder("Paste url anywhere")

	imageContainer := container.NewStack()
	borderLayout := container.NewBorder(urlInput, nil, nil, nil, imageContainer)

	bg := canvas.NewImageFromImage(createBg())
	bg.FillMode = canvas.ImageFillStretch

	urlInput.OnSubmitted = func(s string) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		urlInput.TypedShortcut(selectAll)
		imageContainer.RemoveAll()
		img := loadImage(s)
		if img == nil {
			return
		}
		image0 := canvas.NewImageFromImage(img)
		image0.FillMode = canvas.ImageFillContain
		image1 := canvas.NewImageFromImage(createMosaic(image0.Image.Bounds(), 8))
		image1.FillMode = canvas.ImageFillContain
		imageContainer.RemoveAll()
		imageContainer.Add(bg)
		imageContainer.Add(image1)
		imageContainer.Add(image0)
	}
	w.SetContent(borderLayout)
	w.Canvas().AddShortcut(&fyne.ShortcutPaste{}, func(shortcut fyne.Shortcut) {
		urlInput.SetText(w.Clipboard().Content())
		urlInput.TypedShortcut(selectAll)
		w.Canvas().Focus(urlInput)
	})
	w.Resize(fyne.Size{
		Width: 500, Height: 600,
	})
	w.ShowAndRun()
}

func createBg() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	greyCell := &image.Uniform{C: color.NRGBA{R: 128, G: 128, B: 128, A: 128}}
	draw.Draw(img, img.Bounds(), greyCell, image.Point{}, draw.Src)
	return img
}

func createMosaic(bounds image.Rectangle, size int) image.Image {
	img := image.NewRGBA(bounds)

	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)
	greyCell := &image.Uniform{C: color.NRGBA{R: 128, G: 128, B: 128, A: 128}}
	X := bounds.Dx()/size + 1
	Y := bounds.Dy()/size + 1
	for x := 0; x < X; x++ {
		for y := 0; y < Y; y++ {
			if x%2 == y%2 {
				draw.Draw(img, image.Rect(x*size, y*size, x*size+size, y*size+size), greyCell, image.Point{}, draw.Src)
			}
		}
	}
	return img
}

func loadImage(str string) (img image.Image) {
	w.SetTitle("loading...")
	resp, err := http.Get(str)
	if err != nil {
		showErrorMessage(w, err)
		return
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		showErrorMessage(w, err)
		return
	}
	if resp.StatusCode > 400 {
		showErrorMessage(w, errors.New(string(bs)))
		return
	}
	size := float64(len(bs)) / 1024 / 1024
	img, format, err := image.Decode(bytes.NewReader(bs))
	if err != nil {
		showErrorMessage(w, err)
		return
	}
	w.SetTitle(fmt.Sprintf("%d x %d  %s  %.2fMB", img.Bounds().Size().X, img.Bounds().Size().Y, strings.ToUpper(format), size))
	return
}

func showErrorMessage(window fyne.Window, err error) {
	// Create a dialog to display the error message
	errDialog := dialog.NewError(err, window)
	// Show the error dialog
	errDialog.Show()
	errDialog.Refresh()
}
