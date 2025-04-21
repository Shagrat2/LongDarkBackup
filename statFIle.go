package main

import (
	_ "embed"
)

//go:embed web/favicon.png
var cFavIconPNG []byte

//go:embed web/style.css
var cStyleCSS []byte

//go:embed web/app.wasm
var cAppWASM []byte
