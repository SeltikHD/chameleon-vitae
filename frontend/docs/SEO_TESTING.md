# Testing Your SEO Implementation

## Quick Test Commands

### 1. View Meta Tags

```bash
curl -s http://localhost:3000 | grep -E '<title>|<meta|<link rel='
```

### 2. View JSON-LD Structured Data

```bash
curl -s http://localhost:3000 | grep -A 50 'application/ld+json'
```

### 3. Check Robots.txt

```bash
curl http://localhost:3000/robots.txt
```

### 4. Check Sitemap

```bash
curl http://localhost:3000/sitemap.xml
```

### 5. View All Public Images

```bash
ls -lh public/*.png public/*.ico public/*.svg
```

## Online Validation Tools

### Google Rich Results Test

1. Open: <https://search.google.com/test/rich-results>
2. Enter: <http://localhost:3000> (or use your deployed URL)
3. Expected: All schemas validate ✅

### Facebook Sharing Debugger

1. Open: <https://developers.facebook.com/tools/debug/>
2. Enter your deployed URL
3. Expected: OG image loads, title and description correct ✅

### Twitter Card Validator

1. Open: <https://cards-dev.twitter.com/validator>
2. Enter your deployed URL
3. Expected: Summary large image card renders ✅

### Schema Markup Validator

1. Open: <https://validator.schema.org/>
2. Paste your page HTML or URL
3. Expected: All JSON-LD schemas validate ✅

### Lighthouse SEO Audit

```bash
lighthouse http://localhost:3000 --only-categories=seo --view
```

Expected: 95-100/100 score ✅

## Browser DevTools Testing

### 1. Open <http://localhost:3000> in Chrome/Firefox

### 2. Open DevTools (F12)

### 3. Check Console - Should have no errors

### 4. Network tab - Check analytics scripts load (only in production)

### 5. View Page Source (Ctrl+U) - Verify meta tags

## What to Look For

### ✅ Good Signs

- Title tag present and descriptive
- Meta description 150-160 characters
- All OG tags present (og:title, og:description, og:image, og:url)
- JSON-LD schemas with no errors
- Images load correctly (check Network tab)
- Canonical URL matches current page
- robots.txt allows crawling
- sitemap.xml accessible

### ❌ Red Flags

- Missing title or description
- Duplicate meta tags
- JSON-LD syntax errors
- 404 errors for images/icons
- robots.txt blocking important pages
- Sitemap returning 404

## Production Checklist

Before deploying:

- [ ] All images generated and optimized
- [ ] Sitemap accessible
- [ ] Robots.txt configured
- [ ] Meta tags on all pages
- [ ] JSON-LD schemas validate
- [ ] OG images render in social media debuggers
- [ ] Lighthouse SEO score 95+
- [ ] No console errors
- [ ] Analytics scripts load (production only)

## Post-Deployment

1. **Submit to Google Search Console**
   - Add property: <https://chameleon-vitae.com>
   - Verify ownership
   - Submit sitemap
   - Monitor index coverage

2. **Test Social Sharing**
   - Share on Facebook
   - Share on Twitter/X
   - Share on LinkedIn
   - Verify OG image displays

3. **Monitor Performance**
   - Google Search Console
   - Vercel Analytics
   - Vercel Speed Insights
   - Core Web Vitals

## Troubleshooting

### OG Image Not Loading

```bash
# Check if file exists
ls -lh public/og-image.png

# Check if accessible
curl -I http://localhost:3000/og-image.png
```

### Sitemap 404

```bash
# Check Nuxt config
cat nuxt.config.ts | grep sitemap

# Verify module installed
bun pm ls | grep sitemap
```

### JSON-LD Errors

```bash
# Extract and validate JSON
curl -s http://localhost:3000 | grep -A 50 'application/ld+json' | python -m json.tool
```

## Success Criteria

Your SEO is successfully implemented when:

1. ✅ Lighthouse SEO score is 95+
2. ✅ Google Rich Results Test shows no errors
3. ✅ Social media cards render correctly
4. ✅ All images load without 404s
5. ✅ Sitemap accessible and valid
6. ✅ robots.txt properly configured
7. ✅ No console errors
8. ✅ Analytics tracking works (production)

## Next Steps After Testing

1. Deploy to production
2. Submit sitemap to Google Search Console
3. Monitor index coverage
4. Track organic traffic growth
5. Optimize based on Search Console data
6. A/B test meta descriptions for better CTR
7. Monitor Core Web Vitals
8. Update content based on user search queries

---

**Last Updated**: January 23, 2026  
**Status**: Ready for Testing ✅
