import { existsSync, readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { dirname, resolve } from "node:path";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const specPath = resolve(root, "docs/api/openapi.yaml");

if (!existsSync(specPath)) {
  console.error("docs/api/openapi.yaml is missing");
  process.exit(1);
}

const spec = readFileSync(specPath, "utf8");
const required = [
  "openapi: 3.1.0",
  "/health:",
  "/auth/register:",
  "/groups:",
  "/groups/{groupId}/participants:",
  "/groups/{groupId}/expenses:",
];
const missing = required.filter((needle) => !spec.includes(needle));

if (missing.length > 0) {
  console.error(`OpenAPI contract is missing required entries: ${missing.join(", ")}`);
  process.exit(1);
}

console.log("OpenAPI contract contains the first-slice endpoints.");
