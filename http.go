package antidote

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func addHttpProtocolIfNotExists(url string) string {
	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		return url
	}

	return "http://" + url
}

// normalizeSourceUrl converts relative URL's like '/css/foo/bar.css' into HTTP requestable URL's
// like 'http://domain.com/css/foo/bar.css' based on the source origin.
func normalizeSourceUrl(assetPath string, origin *url.URL) (string, error) {
	s, err := url.Parse(assetPath)
	if err != nil {
		return "", err
	}

	// Remove '//' from assets. Ex: //foo.bar/baz.css => foo.bar/baz.css
	if strings.HasPrefix(assetPath, "//") {
		assetPath = strings.Replace(assetPath, "//", "", 1)
	}

	// Remove relative '..' paths. Ex: ../../app.css => app.css
	assetPath = strings.Replace(assetPath, "../", "", -1)

	// If the asset path doesn't contain a host (meaning it's a relative path), prefix it with the origin host.
	if len(s.Host) == 0 {
		assetPath = addHttpProtocolIfNotExists(origin.Host + "/" + assetPath)
	} else {
		assetPath = addHttpProtocolIfNotExists(assetPath)
	}

	return assetPath, nil
}
