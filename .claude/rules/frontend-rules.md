# Frontend Rules (PART 16, 17)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- ❌ Client-side rendering (React, Vue, Angular, Svelte, etc.)
- ❌ Require JavaScript for core functionality
- ❌ Client-side routing (SPA)
- ❌ Business logic in JavaScript
- ❌ Desktop-first CSS (use mobile-first)
- ❌ Inline CSS or JavaScript
- ❌ JavaScript alerts (use toast notifications)
- ❌ Use npm, webpack, or any JS build toolchain
- ❌ Generic placeholder content in /server/about or /server/help
- ❌ Stub templates or "coming soon" pages
- ❌ Hardcode colors — use CSS custom properties
- ❌ Put templates or static files outside src/server/static/ and src/server/template/

## CRITICAL - ALWAYS DO
- ✅ Server-side rendering with Go templates only
- ✅ Progressive enhancement (works without JS)
- ✅ Mobile-first responsive CSS
- ✅ Embed all assets via go:embed
- ✅ CSS word-break: break-all for long strings (IPv6, .onion, hashes, tokens)
- ✅ Full admin panel with ALL settings accessible
- ✅ WCAG 2.1 AA accessibility
- ✅ Touch targets minimum 44x44px
- ✅ Dark/light/auto theme via CSS custom properties
- ✅ /server/about content sourced from IDEA.md
- ✅ /server/help content with real endpoints and curl examples

## PUBLIC PAGES (cvedex-specific)
| Page | Description |
|------|-------------|
| / | CVE search homepage |
| /cve/{id} | Per-CVE detail page |
| /cves | CVE browse/list page |
| /vendors | Vendor browse page |
| /products | Product browse page |
| /kev | CISA KEV list |
| /dns | DNS lookup demo |
| /dnssec | DNSSEC DS record display |
| /zone | Zone file download |
| /server/about | App info (from IDEA.md) |
| /server/help | API docs with real examples |

## ADMIN PANEL (/admin/*)
| Page | Description |
|------|-------------|
| /admin/ | Dashboard |
| /admin/sync | Data sync trigger |
| /admin/dnssec | DNSSEC key management |
| /admin/zone | Zone config |
| /admin/logs | Access logs |
| /admin/config | Config viewer |
| /admin/setup | First-run setup wizard |

Session auth required for all /admin/* routes.

## THEMING
```css
:root {
  --color-bg: #1a1a1a;
  --color-text: #e0e0e0;
  /* all colors as CSS custom properties */
}
[data-theme="light"] { /* light overrides */ }
@media (prefers-color-scheme: light) { /* auto mode */ }
```

## LONG STRINGS (REQUIRED CSS)
```css
.long-string, .ip-address, .onion-address, .api-token, .hash {
  word-break: break-all;
  overflow-wrap: break-word;
  font-family: monospace;
}
```

## BREAKPOINTS (mobile-first)
| Target | CSS |
|--------|-----|
| Mobile (base) | No media query |
| Tablet+ | @media (min-width: 768px) |
| Desktop+ | @media (min-width: 1024px) |

---
For complete details, see AI.md PART 16, 17
