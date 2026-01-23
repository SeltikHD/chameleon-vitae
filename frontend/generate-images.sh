#!/bin/bash

# Image generation script for Chameleon Vitae SEO
# This script converts SVG logos to various PNG formats required for SEO

set -e

echo "🎨 Generating SEO images for Chameleon Vitae..."

# Check if we're in the frontend directory
if [ ! -f "public/logo.svg" ]; then
    echo "❌ Error: Run this script from the frontend directory"
    exit 1
fi

cd public

# Check for required tools
if command -v convert &> /dev/null; then
    CONVERTER="imagemagick"
    echo "✓ Using ImageMagick"
elif command -v inkscape &> /dev/null; then
    CONVERTER="inkscape"
    echo "✓ Using Inkscape"
else
    echo "❌ Error: Neither ImageMagick nor Inkscape found"
    echo "Install one of them:"
    echo "  - ImageMagick: sudo apt install imagemagick (Debian/Ubuntu)"
    echo "  - Inkscape: sudo apt install inkscape"
    exit 1
fi

# Function to convert SVG to PNG
convert_svg() {
    local input=$1
    local output=$2
    local size=$3

    if [ "$CONVERTER" = "imagemagick" ]; then
        convert -background none -resize "${size}x${size}" "$input" "$output"
    else
        inkscape "$input" --export-filename="$output" --export-width="$size" --export-height="$size"
    fi

    echo "  ✓ Created $output"
}

# Generate favicon sizes from logo.svg
echo ""
echo "📱 Generating favicons..."
convert_svg "logo.svg" "favicon-16x16.png" 16
convert_svg "logo.svg" "favicon-32x32.png" 32
convert_svg "logo.svg" "apple-touch-icon.png" 180
convert_svg "logo.svg" "android-chrome-192x192.png" 192
convert_svg "logo.svg" "android-chrome-512x512.png" 512
convert_svg "logo.svg" "logo.png" 1024

# Generate .ico file (requires ImageMagick)
if [ "$CONVERTER" = "imagemagick" ]; then
    echo ""
    echo "🔷 Generating favicon.ico..."
    convert -background none logo.svg -define icon:auto-resize=48,32,16 favicon.ico
    echo "  ✓ Created favicon.ico"
else
    echo ""
    echo "⚠️  Skipping favicon.ico (requires ImageMagick)"
    echo "   Convert manually: convert logo.svg -define icon:auto-resize=48,32,16 favicon.ico"
fi

# Generate OG image (1200x630)
echo ""
echo "🖼️  Generating Open Graph image..."
if [ "$CONVERTER" = "imagemagick" ]; then
    convert -background none og-image.svg og-image.png
    echo "  ✓ Created og-image.png"
else
    inkscape og-image.svg --export-filename=og-image.png --export-width=1200 --export-height=630
    echo "  ✓ Created og-image.png"
fi

# Optimize PNGs if optipng is available
if command -v optipng &> /dev/null; then
    echo ""
    echo "🗜️  Optimizing PNG files..."
    optipng -quiet *.png
    echo "  ✓ All PNGs optimized"
else
    echo ""
    echo "💡 Tip: Install optipng to optimize PNG file sizes"
    echo "   sudo apt install optipng"
fi

echo ""
echo "✅ All SEO images generated successfully!"
echo ""
echo "📦 Generated files:"
ls -lh favicon*.png apple-touch-icon.png android-chrome-*.png og-image.png 2>/dev/null || true
ls -lh favicon.ico 2>/dev/null || true

echo ""
echo "🎉 Done! Your images are ready for production."
