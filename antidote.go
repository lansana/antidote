package antidote

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// Ingredients object represents options for Antidote.
type Ingredients struct {
	URL string
}

// Antidote object provides the APi operation methods for curing a site.
type Antidote struct {
	ingredients *Ingredients
	parsedUrl   *url.URL
	website     *goquery.Document
	curedHtml   string
}

// New creates a new instance of an Antidote pointer.
func New() *Antidote {
	return new(Antidote)
}

// Mix sets the options of Antidote.
func (a *Antidote) Mix(ingredients *Ingredients) {
	a.ingredients = ingredients
}

// Html retrieves the cured HTML (it will be empty if called before Antidote.Cure() has been called).
func (a *Antidote) Html() string {
	return a.curedHtml
}

// Cure will begin running the algorithms to cure a websites source of any CORS
// restrictions enforced by browsers.
func (a *Antidote) Cure() (string, error) {
	var err error

	if a.ingredients == nil {
		return "", errors.New("Antidote.Mix() must be called before Antidote.Cure().")
	}

	a.parsedUrl, err = url.Parse(a.ingredients.URL)
	if err != nil {
		return "", err
	}

	a.website, err = goquery.NewDocument(a.ingredients.URL)
	if err != nil {
		return "", err
	}

	a.cureAssets()

	a.curedHtml, err = a.website.Html()
	if err != nil {
		return "", err
	}

	return a.curedHtml, nil
}

// cureAssets will run all cure methods concurrently and wait for them to be complete.
func (a *Antidote) cureAssets() {
	var wg sync.WaitGroup
	wg.Add(3)

	go (func() {
		defer wg.Done()
		a.cureCSS()
	})()

	go (func() {
		defer wg.Done()
		a.cureJS()
	})()

	go (func() {
		defer wg.Done()
		a.cureImages()
	})()

	wg.Wait()
}

// cureCSS will fetch the CSS source of all <link> elements concurrently and wait for them to be complete.
// Then it will append a <style> node in the <head> with the raw CSS as the content, and remove the
// pre-existing <link> referencing the external so the browser doesn't throw any errors.
func (a *Antidote) cureCSS() {
	links := a.website.Find("link")

	var wg sync.WaitGroup
	wg.Add(links.Length())

	links.Each(func(index int, link *goquery.Selection) {
		go (func() {
			defer wg.Done()

			if href, ok := link.Attr("href"); ok {
				matchedExtension, err := hasExtension(href, ".css")
				if err != nil {
					log.Println(err)
					return
				}

				if matchedExtension != "" {
					normalizedHref, err := normalizeSourceUrl(href, a.parsedUrl)
					if err != nil {
						log.Println(err)
						return
					}

					source, err := fetch(normalizedHref)
					if err != nil {
						log.Println(err)
						return
					}

					link.AfterHtml(fmt.Sprintf(`<style>%s</style>`, source))
					link.Remove()
				}
			}
		})()
	})

	wg.Wait()
}

// cureJS will fetch the JS source of all <script> elements concurrently and wait for them to be complete.
// Then it will append a <script> node in the <head> with the raw JS as the content, and remove the
// pre-existing <script> referencing the external JS so the browser doesn't throw any errors.
func (a *Antidote) cureJS() {
	scripts := a.website.Find("script")

	var wg sync.WaitGroup
	wg.Add(scripts.Length())

	scripts.Each(func(index int, script *goquery.Selection) {
		go (func() {
			defer wg.Done()

			if src, ok := script.Attr("src"); ok {
				matchedExtension, err := hasExtension(src, ".js")
				if err != nil {
					log.Println(err)
					return
				}

				if matchedExtension != "" {
					normalizedSrc, err := normalizeSourceUrl(src, a.parsedUrl)
					if err != nil {
						log.Println(err)
						return
					}

					source, err := fetch(normalizedSrc)
					if err != nil {
						log.Println(err)
						return
					}

					script.AfterHtml(fmt.Sprintf(`<script>%s</script>`, source))
					script.Remove()
				}
			}
		})()
	})

	wg.Wait()
}

var isImageExtension map[string]bool = map[string]bool{
	"JPEG": true,
	"jpeg": true,
	"JPG":  true,
	"jpg":  true,
	"GIF":  true,
	"gif":  true,
	"PNG":  true,
	"png":  true,
	"BMP":  true,
	"bmp":  true,
	"TIFF": true,
	"tiff": true,
}

// cureImages will fetch the image of all <img> elements concurrently and wait for them to be complete.
// Then it will convert the image into a base64 data URL and replace the src value with the data URL.
func (a *Antidote) cureImages() {
	images := a.website.Find("img")

	var wg sync.WaitGroup
	wg.Add(images.Length())

	images.Each(func(index int, img *goquery.Selection) {
		go (func() {
			defer wg.Done()

			if src, ok := img.Attr("src"); ok {
				imgExtensions := make([]string, len(isImageExtension), 0)
				for k, _ := range isImageExtension {
					imgExtensions = append(imgExtensions, "."+k)
				}

				matchedExtension, err := hasExtension(src, imgExtensions...)
				if err != nil {
					log.Println(err)
					return
				}

				if matchedExtension != "" {
					normalizedSrc, err := normalizeSourceUrl(src, a.parsedUrl)
					if err != nil {
						log.Println(err)
						return
					}

					source, err := fetch(normalizedSrc)
					if err != nil {
						log.Println(err)
						return
					}

					img.SetAttr(
						"src",
						fmt.Sprintf(
							"data:image/%s;base64,%s",
							strings.ToLower(matchedExtension),
							base64.StdEncoding.EncodeToString([]byte(source)),
						),
					)
				}
			}
		})()
	})

	wg.Wait()
}

// hasExtension matches an extension to a URL. If there is a match, the extension is returned.
func hasExtension(src string, extensions ...string) (string, error) {
	for _, extension := range extensions {
		found, err := regexp.MatchString(extension, src)
		if err != nil {
			return "", err
		}
		if found {
			return extension, nil
		}
	}

	return "", nil
}
