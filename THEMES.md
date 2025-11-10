# Blugo Theme System

Blugo supports two color theme modes to match your terminal environment.

## Theme Modes

### 1. ANSI Colors (Recommended) ⭐

**Mode:** `theme_mode = "ansi"`

Uses your terminal's standard ANSI color palette (colors 0-15). This automatically respects whatever terminal theme you're using.

**Advantages:**
- Automatically matches your terminal theme
- Works with any color scheme (Catppuccin, Gruvbox, Dracula, Nord, Solarized, etc.)
- No configuration needed
- Consistent with other terminal applications
- Most compatible option

**Configuration:**
```toml
theme_mode = "ansi"
```

**Color Mapping:**
- Primary: ANSI 5 (Magenta)
- Secondary: ANSI 6 (Cyan)
- Accent: ANSI 4 (Blue)
- Success: ANSI 2 (Green)
- Warning: ANSI 3 (Yellow)
- Error: ANSI 1 (Red)
- Muted: ANSI 8 (Bright Black/Gray)
- Border: ANSI 8 (Bright Black/Gray)
- Highlight: ANSI 15 (Bright White)
- Info: ANSI 6 (Cyan)
- Background: ANSI 0 (Black)
- Foreground: ANSI 7 (White)

---

### 2. TrueColor (Original)

**Mode:** `theme_mode = "truecolor"`

Uses the original hardcoded Blugo color scheme (256-color palette).

**Advantages:**
- Consistent appearance across all terminals
- Original Blugo look and feel
- Well-tested color combinations

**Configuration:**
```toml
theme_mode = "truecolor"
```

**Colors Used:**
- Primary: #ff5fd7 (Magenta)
- Secondary: #5fd7ff (Cyan)
- Success: #5fff00 (Green)
- Warning: #ffaf00 (Orange)
- Error: #ff0000 (Red)
- And more...

---

## How to Configure

Edit your configuration file at `~/.config/blugo/config.toml`:

```toml
# Use your terminal's colors (recommended)
theme_mode = "ansi"

# Or use original Blugo colors
# theme_mode = "truecolor"
```

Then restart Blugo for changes to take effect.

---

## Examples

### Example 1: Use terminal's color scheme (most users)

```toml
# ~/.config/blugo/config.toml
theme_mode = "ansi"
```

This will automatically use whatever colors your terminal is configured with. If you change your terminal theme, Blugo will automatically match it.

**Popular Terminal Themes That Work:**
- Catppuccin (Mocha, Latte, Frappé, Macchiato)
- Gruvbox (Dark, Light)
- Dracula
- Nord
- Solarized (Dark, Light)
- Tokyo Night
- One Dark / One Light
- And many more!

---

### Example 2: Use original Blugo colors

```toml
# ~/.config/blugo/config.toml
theme_mode = "truecolor"
```

This gives you the original Blugo appearance with carefully selected colors.

---

## Testing Different Themes

You can quickly test different theme modes by editing `~/.config/blugo/config.toml` and restarting Blugo:

```bash
# Try ANSI mode (terminal colors)
echo 'theme_mode = "ansi"' >> ~/.config/blugo/config.toml
blugo

# Try TrueColor mode (original colors)
echo 'theme_mode = "truecolor"' >> ~/.config/blugo/config.toml
blugo
```

---

## Why ANSI Mode is Recommended

ANSI mode is the best choice for most users because:

1. **Automatic Matching**: Changes with your terminal theme automatically
2. **Accessibility**: Respects your color preferences and accessibility needs
3. **Consistency**: Matches other terminal applications
4. **No Configuration**: Works out of the box
5. **Universal**: Works with all terminal emulators

Whether you're using Alacritty, Kitty, WezTerm, iTerm2, GNOME Terminal, Konsole, or any other terminal, ANSI mode will make Blugo look great.

---

## Troubleshooting

### Colors look wrong

If colors appear incorrect:

1. Make sure your terminal supports 256 colors or true color
2. Check that your `$TERM` variable is set correctly:
   ```bash
   echo $TERM
   # Should be something like: xterm-256color, alacritty, tmux-256color
   ```
3. Try switching to `truecolor` mode if ANSI mode doesn't look good

### Want consistent colors across terminals

If you want Blugo to look the same regardless of your terminal theme, use:
```toml
theme_mode = "truecolor"
```

### Terminal doesn't support colors

If you're on a very old terminal that doesn't support colors, try:
```toml
theme_mode = "ansi"
```

ANSI mode uses basic 16 colors which are supported by virtually all terminals.

---

## Changing Your Terminal Theme

When using `theme_mode = "ansi"`, Blugo automatically adapts when you change your terminal's color scheme. No need to configure anything in Blugo!

**Examples:**
- Change Alacritty theme → Blugo updates automatically
- Switch Kitty theme → Blugo updates automatically
- Change iTerm2 profile → Blugo updates automatically
