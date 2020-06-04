Antidote
========
Antidote crawls a web page and 'cures' the HTML by loading assets (CSS, JS, images) directly into the DOM.

The effect is that the final HTML will have zero external HTTP calls on page load. 

## Usage

#### Getting cured HTML

```go
a := antidote.New()
a.Mix(&antidote.Ingredients{URL: "https://www.website.com"})

html, err := a.Cure()
if err != nil {
	log.Fatal(err)
}
```

#### Getting cured HTML and saving it for later

```go
a := antidote.New()
a.Mix(&antidote.Ingredients{URL: "https://www.website.com"})

if _, err := a.Cure() != nil {
	log.Fatal(err)
}

// `a.Html()` will contain the cured HTML. Use at your leisure.
```

#### Saving the HTML to a file

```go
a := antidote.New()
a.Mix(&antidote.Ingredients{URL: "http://www.website.com"})

html, err := a.Cure()
if err != nil {
	log.Fatal(err)
}

f, err := os.OpenFile("./website.html", os.O_WRONLY|os.O_CREATE, 0666)
if err != nil {
	log.Fatal(err)
}

f.Write([]byte(html))
```

## What works

- [x] **Convert CSS assets to raw source**

```html
<!-- This -->
<link rel="stylesheet" href="foo.com/bar.css" />

<!-- To this -->
<style>
    // Raw CSS here...
</style>
```

- [x] **Convert JS assets to raw source**

```html
<!-- This -->
<script src="foo.com/bar.js"></script>

<!-- To this -->
<script>
    // Raw JS here...
</script>
```

- [x] **Convert image src values to base64 data URL's**

```html
<!-- This -->
<img src="foo.com/bar.png" />

<!-- To this -->
<img src="data:image/png;base64,abcd..." />
```

## TODO

- [ ] **Convert CSS property URL's to base64 data URL's**

```html
<!-- This -->
<style>
    body {
        background-image: url(../foo/bar.png);
    }
</style>

<!-- To this -->
<style>
    body {
        background-image: url(data:image/png;base64,abcd...);
    }
</style>
```

Contributing
============

In general, this project follows the "fork-and-pull" Git workflow.

 1. **Fork** the repo on GitHub
 2. **Clone** the project to your own machine
 3. **Commit** changes to your own branch
 4. **Push** your work back up to your fork
 5. Submit a **Pull request** so that we can review your changes

NOTE: Be sure to merge the latest from "master" before making a pull request!

Licensing
=========
Antidote is licensed under the MIT License. See [LICENSE](https://github.com/lansana/antidote/blob/master/LICENSE) for the full license text.

