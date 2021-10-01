# retrotxt-fonts

retrotxt-fonts is a tool to dynamically generate the HTML and CSS in use by [RetroTxt](https://github.com/bengarrett/RetroTxt) for [The Oldschool PC](https://int10h.org/oldschool-pc-fonts/) fonts.

The tool requires [Go](https://golang.org/) and the RetroTxt repository.

### Usage

```bash
# clone the repositories if they don't exist
git clone git@github.com:bengarrett/retrotxt-fonts.git
git clone git@github.com:bengarrett/RetroTxt.git

cd retrotxt-fonts
go run .
```

### Update or restore The Oldschool PC fonts

The process relies on [Python 3](https://www.python.org/) scripts and libraries that need to be installed.

```bash
# change directory to the font sources
cd RetroTxt/fonts

# unzip webfonts
unzip -j oldschool_pc_font_pack_v2.2_web.zip "woff - Web (webfonts)/*.woff" -d ../ext/fonts

# change directory to the extension fonts
cd RetroTxt/ext/fonts

# if missing, install python libraries
pip3 install fontTools[woff]
pip3 install Brotli

# convert woff fonts into better compressed woff2 format
python3 woff-to-woff2.py

# cleanup and remove woff fonts
python3 woff-cleanup.py
```