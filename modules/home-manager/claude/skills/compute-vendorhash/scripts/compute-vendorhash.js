#!/usr/bin/env node
const fs = require("fs");
const { spawnSync } = require("child_process");

const PLACEHOLDER_PATTERNS = [
  /sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=/g,
  /sha256-0000000000000000000000000000000000000000000=/g,
];

function findPlaceholder(content) {
  for (const pattern of PLACEHOLDER_PATTERNS) {
    const match = content.match(pattern);
    if (match) return match[0];
  }
  return null;
}

function extractHash(stderr) {
  const gotMatch = stderr.match(/got:\s+(sha256-[A-Za-z0-9+/=]+)/);
  if (gotMatch) return gotMatch[1];

  const sriMatch = stderr.match(/sha256-[A-Za-z0-9+/=]{43,}=/);
  if (sriMatch) return sriMatch[0];

  return null;
}

function main() {
  const [filePath, attrName] = process.argv.slice(2);

  if (!filePath) {
    console.error("Usage: compute-vendorhash.js <path-to-nix-file> [flake-attribute]");
    console.error("");
    console.error("Examples:");
    console.error('  compute-vendorhash.js packages/gitnexus/default.nix "#myGoTool"');
    console.error('  compute-vendorhash.js ./package.nix');
    process.exit(1);
  }

  if (!fs.existsSync(filePath)) {
    console.error(`File not found: ${filePath}`);
    process.exit(1);
  }

  let content = fs.readFileSync(filePath, "utf-8");
  const placeholder = findPlaceholder(content);

  if (!placeholder) {
    console.error("No placeholder hash found in file.");
    console.error("Expected one of:");
    console.error('  sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=');
    console.error('  sha256-0000000000000000000000000000000000000000000=');
    process.exit(1);
  }

  console.error(`Found placeholder: ${placeholder}`);
  console.error("Building to extract correct hash...");

  const buildTarget = attrName || filePath;
  const result = spawnSync("nix", ["build", buildTarget, "--no-link"], {
    encoding: "utf-8",
    maxBuffer: 10 * 1024 * 1024,
  });

  const output = result.stderr || result.stdout || "";
  const correctHash = extractHash(output);

  if (!correctHash) {
    console.error("Could not extract hash from build output.");
    console.error("Build stderr:");
    console.error(output.slice(0, 5000));
    process.exit(1);
  }

  console.error(`Extracted hash: ${correctHash}`);

  const newContent = content.replace(placeholder, correctHash);
  fs.writeFileSync(filePath, newContent);

  console.error(`Replaced placeholder in ${filePath}`);
  console.log(correctHash);
}

main();
