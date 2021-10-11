import {setCacheNameDetails, skipWaiting} from "workbox-core";
import {enable} from "workbox-navigation-preload";
import {cleanupOutdatedCaches, createHandlerBoundToURL, precacheAndRoute} from "workbox-precaching";
import {NavigationRoute, registerRoute} from "workbox-routing";
import {CacheFirst} from "workbox-strategies";

setCacheNameDetails({
    prefix: "app",
});

enable();
skipWaiting();
cleanupOutdatedCaches();

// @ts-ignore
export const manifests = __manifests

precacheAndRoute(
    manifests.concat({
        revision: `${Date.now()}`,
        url: "/",
    }),
);

registerRoute(new NavigationRoute(createHandlerBoundToURL("/")));

registerRoute(
    /\.(?:chunk\.js)$/,
    new CacheFirst({
        cacheName: "chunk-cache",
    }),
    "GET",
);

registerRoute(
    /\.(?:png|jpg|jpeg|svg|webp)$/,
    new CacheFirst({
        cacheName: "image-cache",
    }),
    "GET",
);
