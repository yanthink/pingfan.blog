package markdown

import (
	"blog/app/markdown/katex"
	"bytes"
	_ "embed"
	"fmt"
	chromaHtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	rendererHtml "github.com/yuin/goldmark/renderer/html"
)

//go:embed markdown.css
var css string

func Parse(source []byte, opts ...parser.ParseOption) (html string, err error) {
	p := bluemonday.UGCPolicy()
	p.AllowStyling()
	source = p.SanitizeBytes(source)

	var buf bytes.Buffer

	err = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.NewCJK(extension.WithEastAsianLineBreaks(extension.EastAsianLineBreaksCSS3Draft), extension.WithEscapedSpace()),
			emoji.Emoji,
			&katex.Extender{},
			highlighting.NewHighlighting(
				highlighting.WithStyle("doom-one2"), // https://xyproto.github.io/splash/docs/all.html
				highlighting.WithFormatOptions(
					// chromaHtml.WithLineNumbers(true),
					chromaHtml.TabWidth(4),
				),
			),
		),
		goldmark.WithRendererOptions(
			rendererHtml.WithHardWraps(),
			rendererHtml.WithXHTML(),
			rendererHtml.WithUnsafe(),
		),
	).Convert(source, &buf, opts...)

	if err == nil {
		html = fmt.Sprintf(`<div class="markdown-body"><style>%s</style>%s</div>`, css, buf.Bytes())
	}

	return
}
