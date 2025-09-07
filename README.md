# Goldmark Superscript Extension

<!-- markdownlint-disable MD033 -->

[![Go Reference](https://pkg.go.dev/badge/github.com/zmtcreative/gm-superscript.svg)](https://pkg.go.dev/github.com/zmtcreative/gm-superscript)
[![Go version](https://img.shields.io/github/go-mod/go-version/zmtcreative/gm-superscript)](https://github.com/zmtcreative/gm-superscript)
[![License](https://img.shields.io/github/license/zmtcreative/gm-superscript)](./LICENSE.md)
![GitHub Tag](https://img.shields.io/github/v/tag/zmtcreative/gm-superscript?include_prereleases&sort=semver)

A [Goldmark](https://github.com/yuin/goldmark) extension that adds superscript support using single-caret syntax (`x^2^`). This extension allows you to render superscripts in your Markdown documents as HTML `<sup>` elements.

## Installation

```bash
go get github.com/zmtcreative/gm-superscript
```

## Configuration

### Basic Usage

```go
package main

import (
    "bytes"
    "fmt"

    "github.com/yuin/goldmark"
    "github.com/zmtcreative/gm-superscript"
)

func main() {
    md := goldmark.New(
        goldmark.WithExtensions(
            superscript.Superscript, // Use the pre-configured instance
        ),
    )

    var buf bytes.Buffer
    if err := md.Convert([]byte("x^2^"), &buf); err != nil {
        panic(err)
    }
    fmt.Print(buf.String()) // Output: <p>x<sup>2</sup></p>
}
```

### Alternative Configuration

```go
md := goldmark.New(
    goldmark.WithExtensions(
        superscript.NewSuperscript(), // Create a new instance
    ),
)
```

### With Other Extensions

This extension works seamlessly with other Goldmark extensions, including the built-in footnote extension:

```go
import (
    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/extension"
    "github.com/zmtcreative/gm-superscript"
)

md := goldmark.New(
    goldmark.WithExtensions(
        extension.GFM,
        extension.DefinitionList,
        extension.Footnote,             // Uses [^id] syntax for footnotes
        superscript.Superscript,        // Add superscript support
    ),
)
```

## Basic Examples

### Simple Mathematical Expressions

```markdown
x^2^
```

Renders as: x<sup>2</sup>

```markdown
a^2^ + b^2^ = c^2^
```

Renders as: a<sup>2</sup> + b<sup>2</sup> = c<sup>2</sup>

```markdown
x = y^6^ + z^n+1^
```

Renders as: x = y<sup>6</sup> + z<sup>n+1</sup>

### Advanced Examples

```markdown
a^2!^, b^2,1^, c^n+1^
```

Renders as: a<sup>2!</sup>, b<sup>2,1</sup>, c<sup>n+1</sup>

```markdown
a^2&times;n^, b^2&#x1f604;^, c^&#x215f;n^
```

Renders as: a<sup>2×n</sup>, b<sup>2😄</sup>, c<sup>⅟n</sup>

## Compatibility

### Footnote Extension Compatibility

This extension is fully compatible with Goldmark's built-in `extension.Footnote` and shares the use of the `^` character without conflicts:

- **Footnotes** use the syntax `[^id]` for references and `[^id]: content` for definitions
- **Superscripts** use the syntax `^content^` for inline superscript text

Both extensions can be used together without interference.

### Syntax Rules

The superscript extension follows strict parsing rules to ensure compatibility and prevent conflicts:

1. **No whitespace allowed**: Superscripts cannot contain spaces or any whitespace characters between the carets
   - ✅ Valid: `x^2^`, `a^n+1^`
   - ❌ Invalid: `x^2 ^`, `a^ 2^`, `x^ 2 ^`

2. **No line-start or whitespace-preceded superscripts**: Superscripts cannot start at the beginning of a line or immediately after whitespace
   - ✅ Valid: `x^2^` (preceded by 'x')
   - ❌ Invalid: `^2^` (at line start), `x ^2^` (after space)

3. **No nested markdown or HTML**: Content between carets is treated as literal text - no other markdown or HTML tags are processed inside superscripts
   - ✅ Valid: `x^2^`, `a^**bold**^` (renders as a<sup>**bold**</sup>)
   - ❌ The `**bold**` will not be processed as markdown inside the superscript

4. **No empty superscripts**: Empty carets `^^` are not processed as superscripts

5. **No nested carets**: Additional caret characters inside superscripts are not allowed
   - ✅ Valid: `x^2^`, `y^3^`
   - ❌ Invalid: `x^2^2^^` (would be parsed as `x<sup>2</sup>2^^`)

### Use Cases and Limitations

**Best used for:**

- Simple mathematical expressions: `E=mc^2^`
- Chemical formulas: `H^2^O` (though H₂O with subscripts might be more appropriate)
- Ordinal numbers: `1^st^`, `2^nd^`, `3^rd^`
- Simple annotations: `trademark^TM^`

**Not recommended for:**

- Complex mathematical expressions with multiple levels or nested elements
- Content requiring other markdown formatting (bold, italic, links, etc.)
- Multi-word superscripts with spaces

**For complex scenarios, consider:**

- Using HTML `<sup>` tags directly for more complex superscripts
- Using **KaTeX** or **MathJax** for advanced mathematical typesetting
- Using dedicated chemical formula renderers for scientific notation

## License

This project is licensed under the MIT License. See the [LICENSE.md](LICENSE.md) file for details.
