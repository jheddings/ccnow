# Style

Every segment node accepts a `style` object that controls how its output looks.
Style is applied to the final rendered value — segments produce raw data, the
renderer handles presentation.

## Attributes

| Attribute | Type    | Description                                              |
| --------- | ------- | -------------------------------------------------------- |
| `color`   | string  | Text (foreground) color                                  |
| `bgcolor` | string  | Background color                                         |
| `bold`    | boolean | Bold text                                                |
| `italic`  | boolean | Italic text                                              |
| `prefix`  | string  | Text prepended before the segment value                  |
| `suffix`  | string  | Text appended after the segment value                    |

`prefix` and `suffix` are included in the styled output — they inherit the
same color, bold, and italic settings as the segment value.

### Example

```json
{
  "segment": "git.branch",
  "style": {
    "color": "whiteBright",
    "bold": true,
    "prefix": "\ue0a0 ",
    "suffix": " "
  }
}
```

This renders the branch name in bold bright white, preceded by a Nerd Font
branch glyph and followed by a space.

## Colors

ccglow supports three color tiers. Use whichever fits your terminal:

### Named Colors (16 ANSI)

The standard terminal palette. These work everywhere.

| Color          | Bright Variant   |
| -------------- | ---------------- |
| `black`        | `blackBright`    |
| `red`          | `redBright`      |
| `green`        | `greenBright`    |
| `yellow`       | `yellowBright`   |
| `blue`         | `blueBright`     |
| `magenta`      | `magentaBright`  |
| `cyan`         | `cyanBright`     |
| `white`        | `whiteBright`    |

### 256-Color Palette

Pass a numeric string from `"0"` to `"255"` to access the extended palette.
Useful for subtle grays and specific tones without going full truecolor.

```json
{ "color": "240" }
```

### Truecolor (24-bit)

Full RGB via hex string. Requires a truecolor-capable terminal.

```json
{ "color": "#00afff" }
```

### Background Colors

`bgcolor` accepts the same three formats as `color`:

```json
{
  "style": {
    "color": "white",
    "bgcolor": "#DC0000"
  }
}
```

## Putting It Together

A fully styled segment with foreground, background, and decorations:

```json
{
  "segment": "git.branch",
  "style": {
    "color": "white",
    "bgcolor": "#333333",
    "bold": true,
    "prefix": " \ue0a0 ",
    "suffix": " "
  }
}
```

For more on which segments are available and what data they render, see
[SEGMENTS.md](SEGMENTS.md).
