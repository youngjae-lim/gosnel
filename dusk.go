package gosnel

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
)

func (g *Gosnel) TakeScreenShot(pageURL, testName string, w, h float64) {
	page := rod.New().MustConnect().MustIgnoreCertErrors(true).MustPage(pageURL).MustWaitLoad()

	img, _ := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
		Clip: &proto.PageViewport{
			X:      0,
			Y:      0,
			Width:  w,
			Height: h,
			Scale:  1,
		},
	})
	fileName := time.Now().Format("2006-01-02-15-04-05.000000")
	_ = utils.OutputFile(fmt.Sprintf("%s/screenshots/%s-%s.png", g.RootPath, testName, fileName), img)
}

func (g *Gosnel) FetchPage(pageURL string) *rod.Page {
	return rod.New().MustConnect().MustIgnoreCertErrors(true).MustPage(pageURL).MustWaitLoad()
}

func (g *Gosnel) SelectElementById(page *rod.Page, id string) *rod.Element {
	return page.MustElementByJS(fmt.Sprintf("document.getElementById('%s')", id))
}
