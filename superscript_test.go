package superscript

import (
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/testutil"
	subscript "github.com/zmtcreative/gm-subscript"
)

type TestCase struct {
	desc string
	md   string
	html string
}

func TestGoldmarkOnly(t *testing.T) {
	// These tests are to show how Goldmark handles carats by default,
	// without our extension enabled.
	// Since Footnote also uses carats, using the [^id] syntax,
	// we include Footnote here to show that the two extensions do not conflict.
	mdTest := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			// NewSuperscript(),
		),
	)

	testCases := []TestCase{
		{
			desc: "Goldmark only: simple math expression with caret delimiters",
			md:   `a^2^ + b^2^ = c^2^`,
			html: `<p>a^2^ + b^2^ = c^2^</p>`,
		},
		{
			desc: "Goldmark only: simple math expression single caret",
			md:   `a^2 + b^2 = c^2`,
			html: `<p>a^2 + b^2 = c^2</p>`,
		},
		{
			desc: "Goldmark only: footnote using the caret character",
			md:   `Hi, Bob[^1]
[^1]: Close the airlock before removing your helmet!`,
			html: `<p>Hi, Bob<sup id="fnref:1"><a href="#fn:1" class="footnote-ref" role="doc-noteref">1</a></sup></p>
<div class="footnotes" role="doc-endnotes">
<hr>
<ol>
<li id="fn:1">
<p>Close the airlock before removing your helmet!&#160;<a href="#fnref:1" class="footnote-backref" role="doc-backlink">&#x21a9;&#xfe0e;</a></p>
</li>
</ol>
</div>`,
		},
		// {
		// 	desc: "",
		// 	md:   ``,
		// 	html: ``,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			testutil.DoTestCase(mdTest, testutil.MarkdownTestCase{
				Description: tc.desc,
				Markdown:    tc.md,
				Expected:    tc.html,
			}, t)
		})
	}

}

func TestSuperscriptCore(t *testing.T) {
	mdTest := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			NewSuperscript(),
		),
	)

	testCases := []TestCase{
		{
			desc: "Superscript: x squared example",
			md:   `x^2^`,
			html: `<p>x<sup>2</sup></p>`,
		},
		{
			desc: "Superscript: simple math expression with caret delimiters",
			md:   `a^2^ + b^2^ = c^2^`,
			html: `<p>a<sup>2</sup> + b<sup>2</sup> = c<sup>2</sup></p>`,
		},
		{
			desc: "Superscript: more advanced usage",
			md:   `x = y^6^ + z^n+1^`,
			html: `<p>x = y<sup>6</sup> + z<sup>n+1</sup></p>`,
		},
		{
			desc: "Superscript: other symbols inside superscript",
			md:   `a^2!^, b^2,1^, c^n+1^`,
			html: `<p>a<sup>2!</sup>, b<sup>2,1</sup>, c<sup>n+1</sup></p>`,
		},
		{
			desc: "Superscript: HTML entities inside superscript",
			md:   `a^2&times;n^, b^2&#x1f604;^, c^&#x215f;n^`,
			html: `<p>a<sup>2×n</sup>, b<sup>2😄</sup>, c<sup>⅟n</sup></p>`,
		},
		{
			desc: "Superscript: invalid syntax - no closing caret",
			md:   `a^2 + b^2 = c^2`,
			html: `<p>a^2 + b^2 = c^2</p>`,
		},
		{
			desc: "Superscript: invalid syntax - no leading or trailing spaces",
			md:   `a^2 ^ + b^ 2^ = c^ 2 ^`,
			html: `<p>a^2 ^ + b^ 2^ = c^ 2 ^</p>`,
		},
		{
			desc: "Superscript: invalid syntax - no interior spaces",
			md:   `a^2 a^ + b^b2^ = c^2 foo^`,
			html: `<p>a^2 a^ + b<sup>b2</sup> = c^2 foo^</p>`,
		},
		{
			desc: "Superscript: no nested superscripts",
			md:   `a^2^2^^ + b^2^ = c^2^`,
			html: `<p>a<sup>2</sup>2^^ + b<sup>2</sup> = c<sup>2</sup></p>`,
		},
		{
			desc: "Superscript: sequencial superscripts on same level",
			md:   `a^2^^2^ + b^2^ = c^2^`,
			html: `<p>a<sup>2</sup><sup>2</sup> + b<sup>2</sup> = c<sup>2</sup></p>`,
		},
		{
			desc: "Superscript: footnote with no superscript",
			md:   `Hi, Bob![^1]
[^1]: Close the airlock before removing your helmet!`,
			html: `<p>Hi, Bob!<sup id="fnref:1"><a href="#fn:1" class="footnote-ref" role="doc-noteref">1</a></sup></p>
<div class="footnotes" role="doc-endnotes">
<hr>
<ol>
<li id="fn:1">
<p>Close the airlock before removing your helmet!&#160;<a href="#fnref:1" class="footnote-backref" role="doc-backlink">&#x21a9;&#xfe0e;</a></p>
</li>
</ol>
</div>`,
		},
		{
			desc: "Superscript: footnote using a superscript in the footnote text",
			md:   `Hi, Albert![^1]
[^1]: E=mc^2^ is a famous equation.`,
			html: `<p>Hi, Albert!<sup id="fnref:1"><a href="#fn:1" class="footnote-ref" role="doc-noteref">1</a></sup></p>
<div class="footnotes" role="doc-endnotes">
<hr>
<ol>
<li id="fn:1">
<p>E=mc<sup>2</sup> is a famous equation.&#160;<a href="#fnref:1" class="footnote-backref" role="doc-backlink">&#x21a9;&#xfe0e;</a></p>
</li>
</ol>
</div>`,
		},
		{
			desc: "Superscript: superscript inside square brackets - NOT a footnote",
			md:   `Hi, Albert![^1^]
[^1]: E=mc^2^ is a famous equation.`,
			html: `<p>Hi, Albert![<sup>1</sup>]</p>`,
		},
		// {
		// 	desc: "",
		// 	md:   ``,
		// 	html: ``,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			testutil.DoTestCase(mdTest, testutil.MarkdownTestCase{
				Description: tc.desc,
				Markdown:    tc.md,
				Expected:    tc.html,
			}, t)
		})
	}

}

func TestSuperscriptAdvanced(t *testing.T) {
	mdTest := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			NewSuperscript(),
			subscript.NewSubscript(),
		),
	)

	testCases := []TestCase{
		{
			desc: "Superscript-Subscript: subscript following superscript",
			md:   `a^2^~2~ + b^2^ = c^2^`,
			html: `<p>a<sup>2</sup><sub>2</sub> + b<sup>2</sup> = c<sup>2</sup></p>`,
		},
		{
			desc: "Superscript-Subscript: no nested subscript inside superscript",
			md:   `a^2~2~^ + b^2^ = c^2^`,
			html: `<p>a<sup>2~2~</sup> + b<sup>2</sup> = c<sup>2</sup></p>`,
		},
		{
			desc: "Superscript-Subscript: no nested superscript inside subscript",
			md:   `a~2^2^~ + b^2^ = c^2^`,
			html: `<p>a<sub>2^2^</sub> + b<sup>2</sup> = c<sup>2</sup></p>`,
		},
		// {
		// 	desc: "",
		// 	md:   ``,
		// 	html: ``,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			testutil.DoTestCase(mdTest, testutil.MarkdownTestCase{
				Description: tc.desc,
				Markdown:    tc.md,
				Expected:    tc.html,
			}, t)
		})
	}

}