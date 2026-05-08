import { rmSync } from "node:fs";
import { resolve } from "node:path";

const docs = resolve(import.meta.dirname, "../../docs");
for (const item of [
  "assets",
  "index.html",
  "404.html",
  "manifest.webmanifest",
  "icon.svg",
  "sw.js"
]) {
  rmSync(resolve(docs, item), { force: true, recursive: true });
}

