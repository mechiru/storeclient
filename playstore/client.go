package playstore

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const endpoint = "https://play.google.com/store/apps/details"

type Client struct {
	c    *http.Client
	lang string
}

type Option interface{ apply(c *Client) }

// HttpClient

type httpClient struct{ c *http.Client }

func (hc httpClient) apply(c *Client) { c.c = hc.c }

func HTTPClient(c *http.Client) Option { return &httpClient{c: c} }

// Language

type lang struct{ v string }

func (l lang) apply(c *Client) { c.lang = l.v }

func Lang(v string) Option { return &lang{v: v} }

func NewClient(opts ...Option) *Client {
	c := &Client{c: http.DefaultClient}
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

func (c *Client) defaultParams() url.Values {
	vs := make(url.Values)
	if c.lang != "" {
		vs.Add("hl", c.lang)
	}
	return vs
}

var errEmptyBundleID = &Error{message: "bundle id is empty"}

func (c *Client) getURL(bundleID string) (*url.URL, error) {
	if bundleID == "" {
		return nil, errEmptyBundleID
	}
	url, _ := url.Parse(endpoint)
	vs := c.defaultParams()
	vs.Add("id", bundleID)
	url.RawQuery = vs.Encode()
	return url, nil
}

func (c *Client) Get(ctx context.Context, bundleID string) (*Detail, error) {
	url, err := c.getURL(bundleID)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		io.Copy(ioutil.Discard, resp.Body)
		return nil, &Error{resp.StatusCode, "response status code error"}
	}

	return parseHTML(resp.Body)
}

type Detail struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	CoverArtURL   string `json:"coverArtUrl"`
	ContentRating string `json:"contentRating"`

	// https://play.google.com/store/apps/category/${GenreID}
	GenreID string `json:"genreId"`
	Genre   string `json:"genre"`

	// https://play.google.com/store/apps/dev?id=${DeveloperID}
	DeveloperID string `json:"developerId"`
	Developer   string `json:"developer"`

	// app-ads.txt
	DeveloperURL string `json:"developerUrl"`
	BundleID     string `json:"bundleId"`
	StoreID      string `json:"storeId"`
}

func parseHTML(r io.Reader) (d *Detail, err error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	ret := &Detail{}

	ret.Title = doc.Find(`head > title[id="main-title"]`).First().Text()
	if i := strings.LastIndex(ret.Title, "-"); i > 1 { // At least one character
		ret.Title = ret.Title[:i-1]
	}

	doc.Find("head > meta").Each(func(_ int, s *goquery.Selection) {
		attr, ok := s.Attr("name")
		if !ok {
			return
		}
		content, _ := s.Attr("content")
		switch attr {
		case "description":
			ret.Description = content
		case "appstore:developer_url":
			ret.DeveloperURL = content
		case "appstore:bundle_id":
			ret.BundleID = content
		case "appstore:store_id":
			ret.StoreID = content
		}
	})

	if coverArt, ok := doc.Find(`img[itemprop="image"]`).First().Attr("src"); ok {
		ret.CoverArtURL = coverArt
	}

	s := doc.Find(`a[href^="/store/apps/dev?id="]`).First()
	if href, ok := s.Attr("href"); ok {
		vs, _ := url.ParseQuery(href[strings.Index(href, "?")+1:])
		ret.DeveloperID = vs.Get("id")
	}
	ret.Developer = s.Text()

	s = doc.Find(`a[itemprop="genre"]`).First()
	if href, ok := s.Attr("href"); ok {
		ret.GenreID = strings.TrimPrefix(href, "/store/apps/category/")
	}
	ret.Genre = s.Text()

	// TODO: Avoid select by class
	// ja: img[alt$="歳以上"], en: imp[alt^="Rated for"]
	s = doc.Find(`.E1GfKc`).First()
	ret.ContentRating, _ = s.Attr("alt")

	return ret, nil
}

type Error struct {
	code    int
	message string
}

func (e *Error) Code() int { return e.code }

func (e *Error) Error() string { return fmt.Sprintf("playstore: %s", e.message) }
