const CACHE_NAME = "localhuman-mail-v0.1.0";
const APP_SHELL = ["/localhuman-mail/", "/localhuman-mail/manifest.webmanifest", "/localhuman-mail/icon.svg"];

self.addEventListener("install", (event) => {
  event.waitUntil(caches.open(CACHE_NAME).then((cache) => cache.addAll(APP_SHELL)));
  self.skipWaiting();
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((keys) => Promise.all(keys.filter((key) => key !== CACHE_NAME).map((key) => caches.delete(key))))
  );
  self.clients.claim();
});

self.addEventListener("fetch", (event) => {
  const request = event.request;
  const url = new URL(request.url);
  if (request.method !== "GET" || !url.pathname.startsWith("/localhuman-mail/")) {
    return;
  }

  event.respondWith(fetch(request).catch(() => caches.match(request).then((cached) => cached ?? caches.match("/localhuman-mail/"))));
});

