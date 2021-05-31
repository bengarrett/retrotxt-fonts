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
	"time"
)

const (
	alert = "\u274c"
)

// Paths for named files and directory locations.
type Paths struct {
	RetroTxt  string // path to the local repo: github.com/bengarrett/RetroTxt.
	WoffFonts string // path to woff web fonts.
	DataJSON  string // path to the font json data store.
	SaveCSS   string // output css filename.
	SaveHTML  string // output html filename.
}

func (p *Paths) init(root string) {
	p.RetroTxt = root
	p.WoffFonts = filepath.Join(p.RetroTxt, "ext", "fonts")
	p.DataJSON = filepath.Join(p.RetroTxt, "ext", "json", "font_info.json")
	p.SaveCSS = `fonts.css`
	p.SaveHTML = `fonts.html`
}

// CSS rules.
type CSS struct {
	FontFamily string
	ID         string
	Size       string // font size in pixel or em
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

// Fonts from The Ultimate Oldschool PC Font Pack v2.x.
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
	For        string // label for and input id
	Label      string // the font title displayed to the end user
	Underline  bool   // underline the font label
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
  font-size: {{.Size}};
  line-height: {{.Size}};
}`

// radio input template.
const radioTpl = `<label class="radio" for="{{.For}}">
  <input type="radio" name="{{.Name}}" id="{{.For}}" value="{{.ID}}"> {{if .Underline}}<u>{{.Label}}</u>{{else}}{{.Label}}{{end}}
  <a href="https://int10h.org/oldschool-pc-fonts/fontlist/font?{{.ID}}" target="_blank">
  <svg role="img" class="material-icons has-text-info is-size-7"><use xlink:href="../assets/svg/material-icons.svg#info"></use></svg>
</a></label>`

// header template.
const hTpl = `<hr><h2 class="title has-text-dark is-size-6 {{if not .Usage}} mb-2{{end}}">{{.Origin}}</h2>{{if .Usage}}
<p class="subtitle has-text-dark is-size-7 mb-2">{{.Usage}}</p>{{end}}`

func main() {
	var (
		css   bytes.Buffer
		html  bytes.Buffer
		fonts Fonts
		h     string
		name  Paths
	)
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	name.init(filepath.Join(usr.HomeDir, "github/RetroTxt"))
	raw, err := ioutil.ReadFile(name.DataJSON)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err = json.Unmarshal(raw, &fonts); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	/* Example HTML:
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
	now := time.Now().UTC().Format(time.RFC822Z)
	const start, msdos, video, semi = 0, 59, 130, 184
	const box, h1, h10, hr = `<div class="box">`, `<h1 class="title is-size-3 has-text-dark mb-2">`, `</h1>`, `<hr>`
	const info = `<p class="is-size-7">Fonts support the original IBM PC, 256 character encoding (codepage 437); <u>marked</u> fonts expands support to some 780 characters</p>`
	fmt.Fprintf(&html, "<!-- automatic generation begin (%s) -->\n<div>\n", now)
	cnt, errs := 0, 0
	for i := range fonts.FontInfo {
		f := &fonts.FontInfo[i]
		n := f.WebSafeName
		if variant(n) {
			continue
		}
		cnt++
		switch f.Index {
		case start:
			fmt.Fprintln(&html, "<!-- IBM PC -->")
			fmt.Fprintln(&html, box+"<a id=\"ibmpc\"></a>"+h1+"IBM PC &amp; family"+h10+info)
		case msdos:
			fmt.Fprintln(&html, "</div>\n<!-- MS-DOS -->")
			fmt.Fprintln(&html, box+"<a id=\"msdos\"></a>"+h1+"MS-DOS compatibles"+h10+info)
		case video:
			fmt.Fprintln(&html, "</div>\n<!-- Video hardware -->")
			fmt.Fprintln(&html, box+"<a id=\"video\"></a>"+h1+"Video hardware"+h10+info)
		case semi:
			fmt.Fprintln(&html, "</div>\n<!-- Semi-compatible -->")
			fmt.Fprintln(&html, box+"<a id=\"semico\"></a>"+h1+"Semi-compatibles"+h10+info)
		}
		if f.InfotxtOrigins != h {
			head := Header{
				Origin: title(f.InfotxtOrigins),
				Usage:  usage(f.InfotxtUsage),
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
		filename := filepath.Join(name.WoffFonts, fmt.Sprintf("%s.woff", ff))
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			errs++
			fmt.Printf("\n%s Skipped %s, file not found: %q\nDebug: %+v\n", alert, n, filename, f)
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
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Fprintln(&html, s)
	}
	if errs > 0 {
		fmt.Printf("\nScanned through %d records and %d woff files were missing!\n", cnt, errs)
	}
	fmt.Fprintf(&html, "</div></div>\n<!-- automatic generation end (%s) -->\n", now)
	if err := save(&html, name.SaveHTML); err != nil {
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
			Size:       fmt.Sprintf("%dpx", f.FonWoffSzPx),
		}
		s, err := c.String()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Fprintln(&css, s)
	}
	if err := save(&css, name.SaveCSS); err != nil {
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

func title(n string) string {
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

func usage(n string) string {
	s := strings.ReplaceAll(n, "[?]", "")
	s = strings.ReplaceAll(s, "w/", "with ")
	s = strings.ReplaceAll(s, "chars", "characters")
	s = strings.ReplaceAll(s, "char", "character")
	return s
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
