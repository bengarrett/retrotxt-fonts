package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	// path to woff web fonts.
	fontDir = `github/RetroTxt-daily/ext/fonts`
	// path to the json font data store.
	dataFile = `github/RetroTxt-daily/ext/json/font_info.json`
	// output css filename.
	fnHTML = `fonts.html`
	// output html filename.
	fnCSS = `fonts.css`
)

// CSS rules.
type CSS struct {
	FontFamily string
	ID         string
}

func (c CSS) String() (string, error) {
	buf := &bytes.Buffer{}
	t, err := template.New("css").Parse(cssTpl)
	if err != nil {
		return "", fmt.Errorf("input element: %w", err)
	}
	err = t.Execute(buf, c)
	if err != nil {
		return "", fmt.Errorf("input template: %w", err)
	}
	return buf.String(), nil
}

// Fonts from The Ultimate Oldschool PC Font Pack v2.0.
type Fonts struct {
	// autogenerated at https://mholt.github.io/json-to-go
	FontInfo []struct {
		Index          int    `json:"index"`
		WebSafeName    string `json:"web_safe_name"`
		HasPlus        bool   `json:"has_plus"`
		BaseName       string `json:"base_name"`
		HasAspect      bool   `json:"has_aspect"`
		SqAspect       string `json:"sq_aspect"`
		AcAspect       string `json:"ac_aspect"`
		OrigW          int    `json:"orig_w"`
		OrigH          int    `json:"orig_h"`
		FonWoffSzPx    int    `json:"fon_woff_sz_px"`
		TtfSzPx        int    `json:"ttf_sz_px"`
		TtfSzPt        int    `json:"ttf_sz_pt"`
		InfotxtOrigins string `json:"infotxt_origins"`
		InfotxtUsage   string `json:"infotxt_usage"`
	} `json:"font_info"`
}

// Header for groups of similar fonts.
type Header struct {
	Origin string
	Usage  string
}

func (h Header) String() (string, error) {
	buf := &bytes.Buffer{}
	t, err := template.New("header").Parse(hTpl)
	if err != nil {
		return "", fmt.Errorf("header element: %w", err)
	}
	err = t.Execute(buf, h)
	if err != nil {
		return "", fmt.Errorf("header template: %w", err)
	}
	return buf.String(), nil
}

// Radio HTML element values.
type Radio struct {
	Name       string // form name that should be the shared for all radio input elements
	ID         string // unique ID used for JS and CSS assignments
	FontFamily string // assigned CSS font-family value
	Label      string // the font title displayed to the end user
}

func (r Radio) String() (string, error) {
	buf := &bytes.Buffer{}
	t, err := template.New("webpage").Parse(radioTpl)
	if err != nil {
		return "", fmt.Errorf("input element: %w", err)
	}
	err = t.Execute(buf, r)
	if err != nil {
		return "", fmt.Errorf("input template: %w", err)
	}
	return buf.String(), nil
}

// css rule template.
const cssTpl = `@font-face {
	font-family: "{{.ID}}";
	src: url("../fonts/{{.FontFamily}}.woff") format("woff");
	font-display: swap;
  }
  .font-{{.ID}} {
	font-family: {{.ID}};
  }`

// radio input template.
const radioTpl = `<label for="{{.ID}}">
  <input type="radio" name="{{.Name}}" id="{{.ID}}" value="{{.FontFamily}}">{{.Label}}
</label>`

// header template.
const hTpl = `<h2 class="title">{{.Origin}}</h2>{{if .Usage}}
<h3 class="subtitle">{{.Usage}}</h3>{{end}}`

func main() {
	var (
		css   bytes.Buffer
		html  bytes.Buffer
		fonts Fonts
		h     string
	)
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	raw, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, dataFile))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err = json.Unmarshal(raw, &fonts); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for i := range fonts.FontInfo {
		f := &fonts.FontInfo[i]
		n := f.WebSafeName
		if variant(n) {
			continue
		}
		const start, msdos, video, semi = 0, 59, 130, 184
		switch f.Index {
		case start:
			fmt.Fprintln(&html, "<h1 class=\"title\">IBM PC &amp; family</h1>")
		case msdos:
			fmt.Fprintln(&html, "<h1 class=\"title\">MS-DOS compatibles</h1>")
		case video:
			fmt.Fprintln(&html, "<h1 class=\"title\">Video hardware</h1>")
		case semi:
			fmt.Fprintln(&html, "<h1 class=\"title\">Semi-compatibles</h1>")
		}
		if f.InfotxtOrigins != h {
			head := Header{
				Origin: f.InfotxtOrigins,
				Usage:  f.InfotxtUsage,
			}
			s, err := head.String()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			h = f.InfotxtOrigins
			fmt.Fprintln(&html, s)
		}
		ff := fontFamily(f.BaseName)
		filename := filepath.Join(usr.HomeDir, fontDir, fmt.Sprintf("%s.woff", ff))
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Println("! Font file not found:", ff)
			fmt.Println(n)
			continue
		}
		r := Radio{
			Name:       "font",
			ID:         f.WebSafeName,
			FontFamily: ff,
			Label:      f.BaseName,
		}
		s, err := r.String()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Fprintln(&html, s)
	}
	if err := save(&html, fnHTML); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for i := range fonts.FontInfo {
		f := &fonts.FontInfo[i]
		n := f.WebSafeName
		if variant(n) {
			continue
		}
		c := CSS{
			ID:         f.WebSafeName,
			FontFamily: fontFamily(f.BaseName),
		}
		s, err := c.String()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Fprintln(&css, s)
	}
	if err := save(&css, fnCSS); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func fontFamily(b string) string {
	s := strings.ReplaceAll(b, " ", "_")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "_re.", "_re")
	s = strings.ReplaceAll(s, "_:", "_")
	s = strings.ReplaceAll(s, "AT&T", "ATT")
	return fmt.Sprintf("Web_%s", s)
}

func save(b io.WriterTo, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("save create %q: %w", name, err)
	}
	w := bufio.NewWriter(f)
	_, err = b.WriteTo(w)
	if err != nil {
		return fmt.Errorf("save write to %q: %w", name, err)
	}
	return nil
}

func variant(n string) bool {
	const tail = 3
	var end = ""
	if len(n) > tail {
		end = n[len(n)-tail:]
	}
	switch end {
	case "-2x", "-2y":
		return true
	}
	return false
}
