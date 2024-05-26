package urlProcessor

import "recipebot/urlextract"

type Request struct {
	Details urlextract.WordResult
}

type Result struct {
	StoredUrl string
}

type UrlProcessor interface {
	CanHandle(urlextract.UrlType) bool
	Process(*Request) (*Result, error)
}
