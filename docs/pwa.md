# Progressive Web App

cvedex exposes an installable web app manifest, icons, iOS startup images, and a
service worker. The worker caches versioned static assets, localized offline
messages, and server-approved public informational pages. It checks navigation
requests against the server first and displays a translated offline page when
the server cannot be reached and no matching public page copy is available.

Public page copies are keyed by the URL and the rendered theme, language, cookie
consent, and CCPA preference. The server marks only static informational pages
as eligible, and only when the request has no authentication or unknown
personalization cookie. Responses remain `no-store` for HTTP caches. Contact,
security-report, preference, admin, health, API, and other form or data pages
are not stored for offline use. If a response is not explicitly approved for
offline use, the worker removes older copies for that URL.

The service worker does not intercept API requests or queue form submissions.
Forms require a live connection and must be submitted again after connectivity
returns.

On supported browsers, use the browser's install control to add cvedex to the
home screen or desktop. Safari on iOS does not expose the browser install prompt;
use **Share** and then **Add to Home Screen**. The application remains usable
without installing the worker or enabling JavaScript.

When an update is ready, the page offers an update prompt. Choosing **Update
Now** activates the waiting worker and reloads the page; choosing to wait leaves
the current version active until the next visit or reload.
