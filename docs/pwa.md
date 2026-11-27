# Progressive Web App

cvedex exposes an installable web app manifest, icons, iOS startup images, and a
service worker. The worker caches versioned static assets and localized offline
messages. It checks navigation requests against the server first and displays a
translated offline page when the server cannot be reached. It does not cache
HTML pages because their rendered language and theme depend on visitor settings.

The service worker does not intercept API requests or queue form submissions.
Contact, security, and preference forms require a live connection and must be
submitted again after connectivity returns.

On supported browsers, use the browser's install control to add cvedex to the
home screen or desktop. Safari on iOS does not expose the browser install prompt;
use **Share** and then **Add to Home Screen**. The application remains usable
without installing the worker or enabling JavaScript.

When an update is ready, the page offers an update prompt. Choosing **Update
Now** activates the waiting worker and reloads the page; choosing to wait leaves
the current version active until the next visit or reload.
