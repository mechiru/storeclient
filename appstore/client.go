/*
https://affiliate.itunes.apple.com/resources/documentation/itunes-store-web-service-search-api/
https://stackoverflow.com/questions/8839328/itunes-api-lookup-by-bundle-id
*/

package appstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const endpoint = "https://itunes.apple.com"

type Client struct {
	c             *http.Client
	lang, country string
}

type Option interface{ apply(c *Client) }

// HttpClient

type httpClient struct{ c *http.Client }

func (hc httpClient) apply(c *Client) { c.c = hc.c }

func HttpClient(c *http.Client) Option { return &httpClient{c: c} }

// Language

type lang struct{ v string }

func Lang(v string) Option { return &lang{v: v} }

func (l lang) apply(c *Client) { c.lang = l.v }

// Country
// http://en.wikipedia.org/wiki/ISO_3166-1_alpha-2

type country struct{ v string }

func (co country) apply(c *Client) { c.country = co.v }

func Country(v string) Option { return &country{v: v} }

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
		vs.Add("lang", c.lang)
	}
	if c.country != "" {
		vs.Add("country", c.country)
	}
	return vs
}

func (c *Client) lookupURL(k LookupKey) (*url.URL, error) {
	url, err := url.Parse(endpoint + "/lookup")
	if err != nil {
		return nil, err
	}

	vs := c.defaultParams()
	key, val, err := k.toParams()
	if err != nil {
		return nil, err
	}
	vs.Add(key, val)
	url.RawQuery = vs.Encode()

	return url, nil
}

func (c *Client) Lookup(ctx context.Context, k LookupKey) (*LookupResponse, error) {
	url, err := c.lookupURL(k)
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

	if resp.StatusCode < 200 || 300 < resp.StatusCode {
		io.Copy(ioutil.Discard, resp.Body)
		return nil, &Error{resp.StatusCode, "response status code error"}
	}

	var ret LookupResponse
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

type LookupKey struct{ bundleID, storeID string }

func StoreID(v int) LookupKey { return LookupKey{storeID: strconv.Itoa(v)} }

func BundleID(v string) LookupKey { return LookupKey{bundleID: v} }

var errEmptyLookupKey = errors.New("appstore: both `storeID` and `bundleID` are empty")

func (lk LookupKey) toParams() (key, value string, err error) {
	switch {
	case lk.storeID != "":
		key, value = "id", lk.storeID
	case lk.bundleID != "":
		key, value = "bundleId", lk.bundleID
	default:
		err = errEmptyLookupKey
	}
	return
}

type LookupResponse struct {
	ResultCount int             `json:"resultCount"`
	Results     []*LookupResult `json:"results"`
}

type LookupResult struct {
	ScreenshotURLs                     []string  `json:"screenshotUrls"`
	IPADScreenshotURLs                 []string  `json:"ipadScreenshotUrls"`
	AppletvScreenshotURLs              []string  `json:"appletvScreenshotUrls"`
	ArtworkURL60                       string    `json:"artworkUrl60"`
	ArtworkURL100                      string    `json:"artworkUrl100"`
	ArtworkURL512                      string    `json:"artworkUrl512"`
	ArtistViewURL                      string    `json:"artistViewUrl"`
	SupportedDevices                   []string  `json:"supportedDevices"`
	Advisories                         []string  `json:"advisories"`
	IsGameCenterEnabled                bool      `json:"isGameCenterEnabled"`
	Features                           []string  `json:"features"`
	Kind                               string    `json:"kind"`
	TrackCensoredName                  string    `json:"trackCensoredName"`
	LanguageCodesISO2A                 []string  `json:"languageCodesISO2A"`
	FileSizeBytes                      string    `json:"fileSizeBytes"`
	SellerURL                          string    `json:"sellerUrl"`
	ContentAdvisoryRating              string    `json:"contentAdvisoryRating"`
	AverageUserRatingForCurrentVersion float64   `json:"averageUserRatingForCurrentVersion"`
	UserRatingCountForCurrentVersion   int       `json:"userRatingCountForCurrentVersion"`
	AverageUserRating                  float64   `json:"averageUserRating"`
	TrackViewURL                       string    `json:"trackViewUrl"`
	TrackContentRating                 string    `json:"trackContentRating"`
	TrackName                          string    `json:"trackName"`
	TrackID                            int       `json:"trackId"`
	GenreIDs                           []string  `json:"genreIds"`
	ReleaseDate                        time.Time `json:"releaseDate"`
	FormattedPrice                     string    `json:"formattedPrice"`
	PrimaryGenreName                   string    `json:"primaryGenreName"`
	IsVppDeviceBasedLicensingEnabled   bool      `json:"isVppDeviceBasedLicensingEnabled"`
	CurrentVersionReleaseDate          time.Time `json:"currentVersionReleaseDate"`
	ReleaseNotes                       string    `json:"releaseNotes"`
	PrimaryGenreID                     int       `json:"primaryGenreId"`
	SellerName                         string    `json:"sellerName"`
	MinimumOsVersion                   string    `json:"minimumOsVersion"`
	Currency                           string    `json:"currency"`
	Description                        string    `json:"description"`
	ArtistID                           int       `json:"artistId"`
	ArtistName                         string    `json:"artistName"`
	Genres                             []string  `json:"genres"`
	Price                              float64   `json:"price"`
	BundleID                           string    `json:"bundleId"`
	Version                            string    `json:"version"`
	WrapperType                        string    `json:"wrapperType"`
	UserRatingCount                    int       `json:"userRatingCount"`
}

type Error struct {
	code    int
	message string
}

func (e *Error) Code() int { return e.code }

func (e *Error) Error() string { return fmt.Sprintf("appstore: %s", e.message) }
