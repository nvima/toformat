package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/pkg/errors"
)

//go:embed docs/dist/index.html
var WebDoc string

//go:embed img.png
var samplePng []byte

//go:embed img.jpg
var sampleJpg []byte

var contentTypes = []string{"image/png", "image/jpeg", "image/jpg", "image/gif"}
var prefixes = []string{"/png=", "/jpeg=", "/jpg=", "/gif="}
var sampleFiles = map[string]map[string]*[]byte{
	"/img.png": {"image/png": &samplePng},
	"/img.jpg": {"image/jpeg": &sampleJpg},
}

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		for _, prefix := range prefixes {
			if len(r.URL.RequestURI()) > len(prefix) && r.URL.RequestURI()[:len(prefix)] == prefix {
				err := convertImage(prefix, w, r)
				if err != nil {
					errorResponse(w, err)
				}
				return
			}
		}
		setDefaultResponse(w, r)
	})
}

func convertImage(prefix string, w http.ResponseWriter, r *http.Request) error {
	url := r.URL.RequestURI()[len(prefix):]

	contentType, _, err := checkUrl(url)
	if err != nil {
		return errors.New("Invalid image url")
	}

	res, err := loadImage(url)
	if err != nil {
		return errors.New("Failed to load image")
	}

	img, err := decodeImage(contentType, res)
	if err != nil {
		return errors.New("Failed to decode image")
	}

	err = encodeImage(prefix, img, w)
	if err != nil {
		return errors.New("Failed to encode image")
	}

	defer res.Body.Close()
	return nil
}

func encodeImage(prefix string, image image.Image, w http.ResponseWriter) error {
	switch prefix {
	case "/png=":
		w.Header().Set("Content-Type", "image/png")
		return png.Encode(w, image)
	case "/jpeg=":
		w.Header().Set("Content-Type", "image/jpeg")
		return jpeg.Encode(w, image, nil)
	case "/jpg=":
		w.Header().Set("Content-Type", "image/jpg")
		return jpeg.Encode(w, image, nil)
	case "/gif=":
		w.Header().Set("Content-Type", "image/gif")
		return gif.Encode(w, image, nil)
	default:
		return errors.New("Unsupported image type")
	}
}

func decodeImage(contentType string, res *http.Response) (image.Image, error) {
	switch contentType {
	case "image/png":
		imgSrc, err := png.Decode(res.Body)
		if err != nil {
			return nil, err
		}
		//Set White Background for transparent PNGs
		newImg := image.NewRGBA(imgSrc.Bounds())
		draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
		draw.Draw(newImg, newImg.Bounds(), imgSrc, imgSrc.Bounds().Min, draw.Over)
		return newImg, nil
	case "image/jpeg":
		return jpeg.Decode(res.Body)
	case "image/jpg":
		return jpeg.Decode(res.Body)
	case "image/gif":
		return gif.Decode(res.Body)
	}
	return nil, errors.New("Unknown image type")
}

func loadImage(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := spinhttp.Send(req)
	if err != nil {
		return nil, err
	}

	// defer res.Body.Close()
	return res, nil
}

func checkUrl(url string) (string, string, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", "", err
	}

	res, err := spinhttp.Send(req)
	if err != nil {
		return "", "", err
	}

	if !contains(contentTypes, res.Header.Get("Content-Type")) {
		return "", "", errors.New("Unsupported image type")
	}

	defer res.Body.Close()
	return res.Header.Get("Content-Type"), res.Header.Get("Content-Length"), nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func errorResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, err.Error())
}

func setDefaultResponse(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		return
	}
	//Check if Path countains Sample File
	for path, content := range sampleFiles {
		if r.URL.Path == path {
			for contentType, data := range content {
				w.Header().Set("Content-Type", contentType)
				w.Write(*data)
				return
			}
		}
	}
	buildDocResponse(w, r)
}

func buildDocResponse(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("User-Agent"), "curl") {
		buildTerminalDocResponse(w)
	} else {
		buildWebDocResponse(w)
	}
}

func buildWebDocResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%s", WebDoc)
}

func buildTerminalDocResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "# toformat.link Terminal Documentation\n")
	fmt.Fprint(w, "-------------------------------\n")
	fmt.Fprint(w, "## Format to PNG\n")
	fmt.Fprint(w, "curl -L toqr.link/png=https://toformat.link/sample.jpg > img.png\n\n")
	fmt.Fprint(w, "## Format to GIF\n")
	fmt.Fprint(w, "curl -L toqr.link/png=https://toformat.link/sample.jpg > img.gif\n\n")
	fmt.Fprint(w, "## Format to JPG\n")
	fmt.Fprint(w, "curl -L toqr.link/png=https://toformat.link/sample.png > img.jpg\n\n")
	fmt.Fprint(w, "## Format to JPEG\n")
	fmt.Fprint(w, "curl -L toqr.link/png=https://toformat.link/sample.png > img.jpeg\n\n")
}

func main() {}
