package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFontFamily(t *testing.T) {
	t.Parallel()
	f := FontFamily("Arial Italic")
	assert.Equal(t, "Arial_Italic", f)
}

func TestSave(t *testing.T) {
	t.Parallel()
	buf := bytes.Buffer{}
	fmt.Fprintln(&buf, "test")
	name := filepath.Join(os.TempDir(), "retrotxt-fonts.test.txt")
	err := Save(&buf, name)
	defer os.Remove(name)
	require.NoError(t, err)
}

func TestCSS_String(t *testing.T) {
	t.Parallel()
	css := CSS{
		ID:         "Web_Arial",
		FontFamily: "Arial",
		Size:       "12px",
	}
	x, err := css.String()
	require.NoError(t, err)
	assert.Contains(t, x, "Web_Arial")
}

func TestHeader(t *testing.T) {
	t.Parallel()
	head := Header{
		Origin: "Font API",
		Usage:  "Some information on the usage of the font",
	}
	x, err := head.String()
	require.NoError(t, err)
	assert.Contains(t, x, "Some information")
}

func TestRadio_String(t *testing.T) {
	t.Parallel()
	r := Radio{
		Name:       "font",
		ID:         "Web_Arial",
		FontFamily: "Arial",
		For:        "arial",
		Label:      "Arial",
		Underline:  true,
	}
	x, err := r.String()
	require.NoError(t, err)
	assert.Contains(t, x, "Web_Arial")
}

func TestTitle(t *testing.T) {
	t.Parallel()
	n := Title("IBM incl. on-board video")
	assert.Equal(t, `IBM includes <span class="has-text-weight-normal">on-board video</span>`, n)
}

func TestUsage(t *testing.T) {
	t.Parallel()
	u := Usage("Some w/information[?]")
	assert.Equal(t, "Some with information", u)
}

func TestVariant(t *testing.T) {
	t.Parallel()
	v := Variant("regular")
	assert.False(t, v)
	v = Variant("ibm-pc")
	assert.False(t, v)
	v = Variant("ibm-pc-2x")
	assert.True(t, v)
	v = Variant("ibm-pc-2y")
	assert.True(t, v)
	v = Variant("")
	assert.False(t, v)
}
