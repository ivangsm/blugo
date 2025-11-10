#!/bin/bash
# Script to test different theme modes

echo "Testing Blugo Theme System"
echo "=========================="
echo ""

CONFIG_FILE="$HOME/.config/blugo/config.toml"

# Backup existing config if it exists
if [ -f "$CONFIG_FILE" ]; then
    cp "$CONFIG_FILE" "$CONFIG_FILE.backup"
    echo "✓ Backed up existing config to $CONFIG_FILE.backup"
fi

# Create config directory
mkdir -p "$HOME/.config/blugo"

# Test 1: ANSI Mode (recommended)
echo ""
echo "Test 1: ANSI Mode (terminal colors)"
echo "-----------------------------------"
cat > "$CONFIG_FILE" << 'EOF'
language = "en"
show_emojis = true
theme_mode = "ansi"
refresh_interval = 2
pairing_delay = 1000
disconnect_delay = 500
max_terminal_width = 140
show_rssi = true
show_battery = true
show_device_address = true
compact_mode = false
auto_trust_on_pair = true
auto_start_scanning = true
remember_language = true
battery_high_threshold = 60
battery_low_threshold = 30
hide_unnamed_devices = false
min_rssi_threshold = -100
device_timeout = 0
EOF
echo "✓ Created config with theme_mode = 'ansi'"
echo "  This mode uses your terminal's ANSI colors"
echo ""
echo "Press Enter to test (will show version and exit)..."
read
./blugo --version
echo ""

# Test 2: TrueColor Mode
echo ""
echo "Test 2: TrueColor Mode (original blugo colors)"
echo "----------------------------------------------"
sed -i 's/theme_mode = "ansi"/theme_mode = "truecolor"/' "$CONFIG_FILE"
echo "✓ Changed config to theme_mode = 'truecolor'"
echo "  This mode uses the original hardcoded Blugo colors"
echo ""
echo "Press Enter to test..."
read
./blugo --version
echo ""


# Restore backup if it exists
if [ -f "$CONFIG_FILE.backup" ]; then
    echo ""
    echo "Restore original config? (y/n)"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        mv "$CONFIG_FILE.backup" "$CONFIG_FILE"
        echo "✓ Restored original config"
    else
        rm "$CONFIG_FILE.backup"
        echo "✓ Kept new config (removed backup)"
    fi
fi

echo ""
echo "Theme testing complete!"
echo ""
echo "To use a specific theme mode, edit ~/.config/blugo/config.toml"
echo "and set theme_mode to one of: 'ansi', 'truecolor'"
