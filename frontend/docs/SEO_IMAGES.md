# SEO Images Guide

This document provides instructions for generating the required SEO images for Chameleon Vitae.

## Required Images

### Favicon Set

- `favicon.ico` (48x48)
- `favicon-16x16.png`
- `favicon-32x32.png`
- `apple-touch-icon.png` (180x180)
- `android-chrome-192x192.png`
- `android-chrome-512x512.png`

### Open Graph / Social Media

- `og-image.png` (1200x630) - Main Open Graph image
- `logo.png` (512x512) - Logo for structured data

## Design Guidelines

### Color Palette

- Primary: Emerald (#10b981)
- Secondary: Violet (#8b5cf6)
- Background: Zinc (#09090b, #18181b)
- Text: White/Zinc shades

### OG Image Design

The Open Graph image should include:

1. **Chameleon Vitae** branding
2. Tagline: "The CV That Adapts"
3. Key benefit: "AI-Powered Resume Engineering"
4. Chameleon icon/illustration
5. Gradient background (emerald/violet)

### Logo Design

The logo should:

1. Feature a chameleon icon
2. Be recognizable at small sizes
3. Work on both light and dark backgrounds
4. Include the emerald-to-violet gradient

## Generation Tools

You can use these tools to create the images:

1. **Figma** - Professional design tool
2. **Canva** - Quick and easy designs
3. **GIMP** - Free alternative to Photoshop
4. **Favicon Generator** - <https://realfavicongenerator.net/>

## Quick Setup

If you don't have custom images yet, you can use placeholders:

1. Create a simple SVG logo
2. Use ImageMagick to convert to different sizes:

   ```bash
   convert logo.svg -resize 192x192 android-chrome-192x192.png
   convert logo.svg -resize 512x512 android-chrome-512x512.png
   ```

3. For OG image, create a 1200x630 PNG with:
   - Background gradient
   - Logo centered
   - Text overlay with title and tagline

## Optimization

After creating images, optimize them:

```bash
# PNG optimization
optipng *.png
pngquant --quality=65-80 *.png

# JPEG optimization (if using JPG)
jpegoptim --max=85 *.jpg
```

## Placement

All images should be placed in:

```text
frontend/public/
```

The Nuxt configuration will automatically serve them.
