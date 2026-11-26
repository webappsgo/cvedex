# Translating cvedex

cvedex embeds its translation catalog in `src/common/i18n/locales/`. The
supported languages are English (`en`), Spanish (`es`), French (`fr`), German
(`de`), Italian (`it`), and Arabic (`ar`). Arabic uses right-to-left layout.

Add new messages to every locale JSON file using the same nested key in each
file. English is the fallback source locale. Keep locale metadata aligned with
the registry in `src/common/i18n/i18n.go`, preserve placeholders exactly, and
add all CLDR plural categories declared for that locale. The i18n tests check
message-key parity, metadata, placeholders, and plural categories.

Server templates use `{{t .Lang "namespace.message"}}`. Resolve language from
the supported `lang` query choice, consented preference cookie, and
`Accept-Language` header, in that order. Unsupported query choices fall back to
English. Preference cookies are persisted only when the visitor has granted
preference-cookie consent.

Preserve CVE identifiers, DNS names, API field names, CVSS vectors, and upstream
source names exactly. Translate user-facing interface strings while leaving
machine-readable response fields and upstream data unchanged. Review
translations for terminology consistency, accessibility, and text expansion
before release; credit translators by their preferred name only with consent.
