import { copyFileSync, existsSync, readFileSync, writeFileSync } from "node:fs";
import { resolve } from "node:path";

const docs = resolve(import.meta.dirname, "../../docs");
const indexPath = resolve(docs, "index.html");
const fallbackPath = resolve(docs, "404.html");

if (existsSync(indexPath)) {
  copyFileSync(indexPath, fallbackPath);
}

const marker = {
  version: process.env.npm_package_version ?? "0.1.0",
  commit: "unknown",
  buildTime: new Date().toISOString()
};

try {
  marker.commit = readFileSync(resolve(import.meta.dirname, "../../.git/HEAD"), "utf8").trim();
} catch {
  marker.commit = "unknown";
}

writeFileSync(resolve(docs, "build-meta.json"), JSON.stringify(marker, null, 2));
