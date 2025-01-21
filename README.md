# codebase

`codebase` is a command-line tool (written in Go) for aggregating all code files in a project into a single Markdown block, suitable for providing “context” to an LLM (Large Language Model). This helps you easily copy & paste your entire project’s relevant source files into a prompt or a separate context window.

## Key Features

- **Recursive Collection**: Walks through the target directory and gathers all files (relative paths + code blocks).
- **Ignore Rules**: Respects a `.codebase_ignore` file plus additional ignores passed via `--ignore`.
- **Markdown Output**: Combines everything into a single Markdown-format string, labeling each file with its relative path.
- **No Code Dump to Console**: By design, the tool never dumps the entire code to stdout. Instead, you can:
  - Save to a specified `.md` file  
  - Copy directly to your system clipboard.

## Installation
TODO
