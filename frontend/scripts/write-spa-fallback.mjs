import { copyFileSync, existsSync, writeFileSync } from "node:fs";
import { resolve } from "node:path";

const docs = resolve(import.meta.dirname, "../../docs");
const indexPath = resolve(docs, "index.html");
const fallbackPath = resolve(docs, "404.html");

if (existsSync(indexPath)) {
  copyFileSync(indexPath, fallbackPath);
}

const marker = {
  version: process.env.npm_package_version ?? "0.1.0",
  commit: process.env.LOCALHUMAN_BUILD_COMMIT ?? "pages",
  buildTime: process.env.LOCALHUMAN_BUILD_TIME ?? "static"
};

writeFileSync(resolve(docs, "build-meta.json"), JSON.stringify(marker, null, 2));
