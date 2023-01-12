package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"strconv"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/sunshineplan/imgconv"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {

		req, err := http.NewRequest("HEAD", r.URL.RequestURI()[1:], nil)
		if err != nil {
			errorResponse(w, errors.New("Invalid URL"))
			return
		}

		res, err := spinhttp.Send(req)
		if err != nil {
			errorResponse(w, errors.New("Invalid URL"))
			return
		}

		//Print Content-Type to Response
		contentType := res.Header.Get("Content-Type")
		fmt.Fprintln(w, contentType)
		//TODO return wrong content type
		//Print Content-Length
		contentLength := res.Header.Get("Content-Length")
		i, _ := strconv.Atoi(contentLength)
		if i == 0 {
			errorResponse(w, errors.New("No Content"))
			return
		}
		fmt.Fprintln(w, i)
		//Make a new request to get the content
		req, err = http.NewRequest("GET", r.URL.RequestURI()[1:], nil)
		if err != nil {
			errorResponse(w, errors.New("Something went wrong"))
			return
		}
		res, err = spinhttp.Send(req)
		if err != nil {
			errorResponse(w, errors.New("Something went wrong"))
			return
		}
		//Convert the content to image.Image
		img, _, err := image.Decode(res.Body)
		if err != nil {
			errorResponse(w, errors.New("Something went wrong"))
			return
		}

		// var imgRes io.Writer
		var imgBuffer bytes.Buffer
		imgRes := io.MultiWriter(&imgBuffer)
		// imgRes := new(bytes.Buffer)

		opts := imgconv.Options{Format: &imgconv.FormatOption{Format: imgconv.PNG}}
		opts.Convert(imgRes, img)
		w.Header().Set("Content-Type", "image/png")
		imgReader := bytes.NewReader(imgBuffer.Bytes())
		io.Copy(w, imgReader)


	})
}


