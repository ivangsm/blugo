# Theme System Implementation

## Summary

Added a flexible theme system to Blugo that allows the application to automatically adapt to your terminal's color scheme.

## Changes

### New Files Created

1. **`internal/ui/theme.go`** - Core theme system implementation
   - Two theme modes: ANSI and TrueColor
   - ANSI mode: Uses terminal's standard colors (0-15)
   - TrueColor mode: Original hardcoded Blugo colors
   - Dynamic style updates based on selected theme

2. **`THEMES.md`** - Comprehensive theme documentation
   - Detailed explanation of each theme mode
   - Configuration examples
   - Troubleshooting guide
   - Recommended practices

3. **`config.example.toml`** - Example configuration file
   - All theme options documented
   - Practical examples for each theme mode

4. **`test_theme.sh`** - Theme testing script
   - Interactive testing of both theme modes

### Modified Files

1. **`internal/config/config.go`**
   - Added `ThemeMode` field (string): "ansi" or "truecolor"
   - Updated default config to use "ansi" mode
   - Added theme documentation to config header

2. **`internal/ui/styles.go`**
   - Changed hardcoded color values to variables
   - Colors now initialized by theme system
   - Maintains backward compatibility

3. **`cmd/blugo/main.go`**
   - Added theme initialization on startup
   - Loads theme based on config settings
   - Fallback to default theme on error

4. **`README.md`** and **`README.es.md`**
   - Added "Flexible theming" feature
   - New "Theming" section with configuration examples
   - Link to detailed THEMES.md documentation

## Features

### 1. ANSI Mode (Recommended)
- Uses terminal's standard ANSI colors (0-15)
- Automatically respects terminal theme
- Works with any color scheme (Catppuccin, Gruvbox, Dracula, Nord, etc.)
- No configuration needed
- Best compatibility

**Configuration:**
```toml
theme_mode = "ansi"
```

**Benefits:**
- Automatic matching with terminal theme
- Works with ALL terminal emulators (Alacritty, Kitty, WezTerm, iTerm2, GNOME Terminal, etc.)
- Respects accessibility settings
- Zero configuration

### 2. TrueColor Mode
- Original Blugo color scheme
- Uses 256-color palette
- Consistent appearance across terminals
- Hardcoded colors

**Configuration:**
```toml
theme_mode = "truecolor"
```

**Benefits:**
- Predictable colors
- Original Blugo appearance
- Well-tested color combinations

---

## Color Mapping

### ANSI Mode
- Primary: ANSI 5 (Magenta)
- Secondary: ANSI 6 (Cyan)
- Accent: ANSI 4 (Blue)
- Success: ANSI 2 (Green)
- Warning: ANSI 3 (Yellow)
- Error: ANSI 1 (Red)
- Muted: ANSI 8 (Bright Black/Gray)
- Highlight: ANSI 15 (Bright White)
- Info: ANSI 6 (Cyan)

### TrueColor Mode
Uses the original 256-color palette codes:
- Primary: Color 205 (Magenta)
- Secondary: Color 86 (Cyan)
- Success: Color 82 (Green)
- Warning: Color 214 (Orange)
- Error: Color 196 (Red)
- And more...

---

## Implementation Details

### Theme Initialization Flow

1. Application starts (`main.go`)
2. Configuration loaded (`config.Init()`)
3. Theme mode read from config
4. Theme initialized (`ui.InitializeTheme(mode)`)
5. All style variables updated
6. UI rendered with selected theme

### Theme Switching

To switch themes, users simply:
1. Edit `~/.config/blugo/config.toml`
2. Change `theme_mode` value
3. Restart Blugo

The theme is applied immediately on startup.

---

## Benefits

### 1. User Experience
- Blugo now respects user's terminal theme preferences
- No more jarring color mismatches
- Seamless integration with terminal environment
- Better accessibility support

### 2. Simplicity
- Two simple, well-defined modes
- Easy to understand and configure
- No external file dependencies
- Reliable and predictable behavior

### 3. Compatibility
- Works with any terminal emulator
- No special requirements
- Backward compatible with existing installations
- Graceful fallback if config missing

### 4. Accessibility
- Users can use their preferred high-contrast themes
- Better for users with vision impairments
- Respects system-wide theme settings
- ANSI mode works with accessibility tools

---

## Migration Guide

### For Existing Users

No action required! The default configuration uses ANSI mode, which automatically adapts to your terminal's colors.

### To Use Original Blugo Colors

Edit `~/.config/blugo/config.toml`:
```toml
theme_mode = "truecolor"
```

Then restart Blugo.

---

## Testing

Build and test:
```bash
make build
./blugo
```

Test different modes:
```bash
./test_theme.sh
```

Or manually:
```bash
# Test ANSI mode
echo 'theme_mode = "ansi"' > ~/.config/blugo/config.toml
./blugo

# Test TrueColor mode
echo 'theme_mode = "truecolor"' > ~/.config/blugo/config.toml
./blugo
```

---

## Technical Implementation

### File Structure
```
internal/ui/
├── theme.go       # Theme system (85 lines)
│   ├── ThemeMode enum
│   ├── ColorScheme struct
│   ├── InitializeTheme()
│   ├── getANSIColorScheme()
│   ├── getDefaultTrueColorScheme()
│   └── updateStyles()
└── styles.go      # Style definitions (uses theme colors)
```

### Key Design Decisions

1. **No External Dependencies**: ANSI mode doesn't require reading external files
2. **Simple API**: Single function call to initialize theme
3. **Type Safety**: ThemeMode is a string type for safety
4. **Graceful Degradation**: Falls back to ANSI mode on error
5. **Minimal Changes**: Existing style code mostly unchanged

### Performance

- Theme initialization happens once at startup
- No runtime overhead after initialization
- No file I/O during operation (both modes)
- Negligible memory footprint

---

## Future Enhancements

Potential improvements:
- [ ] Live theme switching (without restart)
- [ ] Theme preview command
- [ ] Custom color definitions in config file
- [ ] Dark/light mode detection
- [ ] More predefined color schemes

---

## Conclusion

The theme system provides a simple yet powerful way to customize Blugo's appearance while maintaining excellent compatibility across different terminal emulators and color schemes. The recommended ANSI mode ensures Blugo looks great in any environment with zero configuration.
