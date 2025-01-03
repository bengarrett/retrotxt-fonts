package main

/*
   <!-- Modern fonts -->
   <div class="box mt-4">
     <a id="modern"></a>
     <h1 class="title is-size-3 has-text-dark mb-2">Modern</h1>
     <p class="is-size-7"><u>Marked</u> fonts support a large range of Unicode glyphs and
       languages</p>
     <hr>
     <h2 class="title has-text-dark is-size-6">
       IBM Plex
       <a href="https://www.ibm.com/plex/specs" target="_blank">
         <svg role="img" class="material-icons has-text-info is-size-7">
           <use xlink:href="../assets/svg/material-icons.svg#info"></use>
         </svg>
       </a>
     </h2>
     <p class="subtitle has-text-dark is-size-7 mb-2">Plex was designed to capture IBM’s spirit and
       history
     </p>
     <div class="control">
       <label class="radio" for="ibmplexmono">
         <input type="radio" name="font" id="ibmplexmono" value="ibmplexmono"> <u>Mono Regular</u>
       </label>
       <label class="radio" for="ibmplextlight">
         <input type="radio" name="font" id="ibmplextlight" value="ibmplextlight"> <u>Mono Light</u>
       </label>
       <label class="radio" for="ibmplextmedium">
         <input type="radio" name="font" id="ibmplextmedium" value="ibmplextmedium"> <u>Mono Medium</u>
       </label>
     </div>
     <hr>
*/

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// CSS rules.
type CSS struct {
	FontFamily string // Font-family name
	ID         string // Unique ID for the font-family
	Size       string // Size of the font as pixel or em values
}

func (c CSS) String() (string, error) {
	buf := &bytes.Buffer{}
	t, err := template.New("css").Parse(cssTmpl)
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
	Origin string // font origin or group name
	Usage  string // font usage or description
}

func (h Header) String() (string, error) {
	buf := &bytes.Buffer{}
	t, err := template.New("header").Parse(headTmpl)
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
	For        string // label for and input id
	Label      string // the font title displayed to the end user
	Underline  bool   // underline the font label
}

func (r Radio) String() (string, error) {
	buf := &bytes.Buffer{}
	t, err := template.New("webpage").Parse(radioTmpl)
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
const cssTmpl = `@font-face {
  font-family: "{{.ID}}";
  src: url("../fonts/{{.FontFamily}}.woff") format("woff");
  font-display: swap;
}
.font-{{.ID}} {
  font-family: {{.ID}};
  font-size: {{.Size}};
  line-height: {{.Size}};
}`

// radio input template.
const radioTmpl = `<a href="https://int10h.org/oldschool-pc-fonts/fontlist/font?{{.ID}}" target="_blank">
  <svg role="img" class="material-icons has-text-dark">` +
	`<use xlink:href="../assets/svg/material-icons.svg#info"></use></svg>
</a>
<label for="{{.For}}">
  <input type="radio" name="{{.Name}}" id="{{.For}}" value="{{.ID}}"> ` +
	`{{if .Underline}}<u>{{.Label}}</u>{{else}}{{.Label}}{{end}}
</label>`

// header template.
const headTmpl = `<h2 class="title has-text-dark is-size-6 mt-4` +
	`{{if not .Usage}} mb-2{{end}}">{{.Origin}}</h2>{{if .Usage}}
<h3 class="subtitle has-text-dark is-size-7 mb-2">{{.Usage}}</h3>{{end}}`

func main() {
	const (
		data = `github/RetroTxt/ext/json/font_info.json`
		font = `github/RetroTxt/ext/fonts`
	)
	if err := Generate(data, font); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// Generate creates the HTML and CSS files for the font selection page.
//
// The "data" is the file path to the json font data store.
// The "font" is the directory path to woff web fonts.
func Generate(data, font string) error {
	const fnameHTML = `fonts.html` // output html filename
	const fnameCSS = `fonts.css`   // output css filename
	var (
		css    bytes.Buffer
		html   bytes.Buffer
		fonts  Fonts
		header string
	)
	usr, err := user.Current()
	if err != nil {
		return err
	}
	raw, err := os.ReadFile(filepath.Join(usr.HomeDir, data))
	if err != nil {
		return err
	}
	if err = json.Unmarshal(raw, &fonts); err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC822Z)
	const start, msdos, video, semi = 0, 59, 130, 184
	const h1, h10, hr = `<div class="box mt-4"><h1 class="title is-size-3 has-text-dark mb-2">`, `</h1></div>`, `<hr>`
	const info = `<p class="is-size-7">Fonts support the original IBM PC, 256 character encoding (codepage 437);` +
		` <u>marked</u> fonts expands support to some 780 characters</p>`
	fmt.Fprintf(&html, "<!-- automatic generation begin (%s) -->\n<div>\n", now)
	for i := range fonts.FontInfo {
		f := &fonts.FontInfo[i]
		n := f.WebSafeName
		if Variant(n) {
			continue
		}
		switch f.Index {
		case start:
			// <a id="ibmpc"></a>
			fmt.Fprintln(&html, hr+h1+"IBM PC &amp; family"+h10)
			fmt.Fprintln(&html, info)
		case msdos:
			fmt.Fprintln(&html, hr+h1+"MS-DOS compatibles"+h10)
			fmt.Fprintln(&html, info)
		case video:
			fmt.Fprintln(&html, hr+h1+"Video hardware"+h10)
			fmt.Fprintln(&html, info)
		case semi:
			fmt.Fprintln(&html, hr+h1+"Semi-compatibles"+h10)
			fmt.Fprintln(&html, info)
		}
		if f.InfotxtOrigins != header {
			head := Header{
				Origin: Title(f.InfotxtOrigins),
				Usage:  Usage(f.InfotxtUsage),
			}
			s, err := head.String()
			if err != nil {
				return err
			}
			header = f.InfotxtOrigins
			fmt.Fprintln(&html, s)
		}
		ff := FontFamily(f.BaseName)
		filename := filepath.Join(usr.HomeDir, font, ff+".woff")
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "! Font file not found:", ff)
			fmt.Fprintln(os.Stderr, n)
			continue
		}
		r := Radio{
			Name:       "font",
			ID:         f.WebSafeName,
			FontFamily: ff,
			For:        strings.ToLower(ff),
			Label:      f.BaseName,
			Underline:  f.HasPlus,
		}
		s, err := r.String()
		if err != nil {
			return err
		}
		fmt.Fprintln(&html, s)
	}
	fmt.Fprintf(&html, "</div>\n<!-- automatic generation end (%s) -->\n", now)
	if err := Save(&html, fnameHTML); err != nil {
		return err
	}
	for i := range fonts.FontInfo {
		f := &fonts.FontInfo[i]
		n := f.WebSafeName
		if Variant(n) {
			continue
		}
		c := CSS{
			ID:         f.WebSafeName,
			FontFamily: FontFamily(f.BaseName),
			Size:       fmt.Sprintf("%dpx", f.FonWoffSzPx),
		}
		s, err := c.String()
		if err != nil {
			return err
		}
		fmt.Fprintln(&css, s)
	}
	if err := Save(&css, fnameCSS); err != nil {
		return err
	}
	return nil
}

// FontFamily returns a cleaned up string for the font-family name.
func FontFamily(b string) string {
	s := strings.ReplaceAll(b, " ", "_")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "_re.", "_re")
	s = strings.ReplaceAll(s, "_:", "_")
	s = strings.ReplaceAll(s, "AT&T", "ATT")
	return "Web_" + s
}

// Save writes the buffer to a named file.
func Save(b io.WriterTo, name string) error {
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

// Title returns a cleaned up string for the font title.
func Title(n string) string {
	const (
		span = "<span class=\"has-text-weight-normal\">"
		cls  = "</span>"
	)
	s := n
	if strings.ContainsAny(s, "(") {
		s = strings.Replace(s, "(", span+"(", 1)
		s += cls
	} else if strings.Contains(s, "Multimode Graphics Adapter") {
		s = strings.Replace(s, "Multimode Graphics Adapter", span+"Multimode Graphics Adapter", 1)
		s += cls
	}
	s = strings.ReplaceAll(s, "incl.", "includes")
	s = strings.Replace(s, "Adapter Interface drivers for", span+"Adapter Interface drivers for"+cls, 1)
	s = strings.Replace(s, "series video BIOS", span+"series video BIOS"+cls, 1)
	s = strings.Replace(s, "on-board video", span+"on-board video"+cls, 1)
	s = strings.Replace(s, "system font", span+"system font"+cls, 1)
	s = strings.Replace(s, "system-loaded font", span+"system-loaded font"+cls, 1)
	s = strings.Replace(s, "firmware and system", span+"firmware and system"+cls, 1)
	return s
}

// Usage returns a cleaned up string for the font usage.
func Usage(n string) string {
	s := strings.ReplaceAll(n, "[?]", "")
	s = strings.ReplaceAll(s, "w/", "with ")
	s = strings.ReplaceAll(s, "chars", "characters")
	s = strings.ReplaceAll(s, "char", "character")
	return s
}

// Variant returns true if the font name is a variant.
func Variant(n string) bool {
	const tail = 3
	end := ""
	if len(n) > tail {
		end = n[len(n)-tail:]
	}
	switch end {
	case "-2x", "-2y":
		return true
	}
	return false
}
