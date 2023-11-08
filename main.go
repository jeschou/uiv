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
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"strings"
)

var selectAll = &fyne.ShortcutSelectAll{}
var w fyne.Window

func main() {
	myApp := app.New()
	w = myApp.NewWindow("URL Image Viewer")
	urlInput := widget.NewEntry()
	urlInput.SetPlaceHolder("Paste url anywhere")

	imageContainer := container.NewStack()
	borderLayout := container.NewBorder(urlInput, nil, nil, nil, imageContainer)

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
		imageContainer.RemoveAll()
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
