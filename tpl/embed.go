package tpl

import "embed"

// content holds our static web server content.
//
//go:embed *.tmpl
var Content embed.FS
