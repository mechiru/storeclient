# storeclient

[![ci](https://github.com/mechiru/storeclient/workflows/ci/badge.svg)](https://github.com/mechiru/storeclient/actions?query=workflow:ci)

This library provides a client to get app information.

## Example

### appstore
Use itunes api to get the information.

```go
import "github.com/mechiru/storeclient/appstore"

c := appstore.NewClient(appstore.Lang("ja_jp"), appstore.Country("JP"))
resp, err := c.Lookup(context.Background(), appstore.StoreID(340368403))
if err != nil {
	// TODO: handle error
}
fmt.Printf("response: %#v\n", resp)
```

### playstore
Scraping the app details page to get the information.

```go
import "github.com/mechiru/storeclient/playstore"

c := playstore.NewClient(playstore.Lang("ja"))
resp, err := c.Get(context.Background(), "com.cookpad.android.activities")
if err != nil {
	// TODO: handle error
}
fmt.Printf("response: %#v\n", resp)
```

**LICENSE**<br>
[MIT](./LICENSE)
