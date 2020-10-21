const fs = require("fs");
const path = require("path");
const packageJson = require("../package.json");

const CATALOG_PATH = path.normalize("./CATALOG.md");

const snippetsPaths = getRegisteredSnippets(packageJson);
const buffer = fs.createWriteStream(CATALOG_PATH);

buffer.write("# Catalog\n");
buffer.write("\n");
buffer.write("Here's a list of snippets available at the moment:\n");
buffer.write("\n");

for (const filepath of snippetsPaths) {
  addSnippetSection(filepath, buffer);
}

buffer.close();

/**
 * Get registered snippets on package.json
 */
function getRegisteredSnippets(json) {
  return json.contributes.snippets
    .map((snippet) => path.normalize(snippet.path))
    .sort(sortAlphabetically);
}

/**
 * Add a new snippet section for the output buffer
 */
function addSnippetSection(snippetFilePath, buffer) {
  buffer.write(`## ${path.basename(snippetFilePath, ".code-snippets")}\n\n`);
  addMarkdownTable(snippetFilePath, buffer);
  buffer.write("\n---\n\n");
}

/**
 * Generate a Markdown table with all snippets prefixes
 * and respective descriptions
 */
function addMarkdownTable(snippetFilePath, buffer) {
  const content = JSON.parse(fs.readFileSync(snippetFilePath));

  buffer.write("| Prefix | Description |\n");
  buffer.write("|--------|-------------|\n");

  Object.values(content)
    .sort((a, b) => sortAlphabetically(a.prefix, b.prefix))
    .forEach((element) => {
      buffer.write(`| ${element.prefix} | ${element.description} |\n`);
    });
}

/**
 * Sort "a" and "b" alphabetically
 */
function sortAlphabetically(a, b) {
  return path.basename(a).localeCompare(path.basename(b));
}
