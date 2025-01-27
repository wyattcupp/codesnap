# codesnap

`codesnap` is a command-line tool (written in Go) for aggregating all code files in a project into a single Markdown block, suitable for providing “context” to an LLM (Large Language Model). This helps you easily copy & paste your entire project’s relevant source files into a prompt or a separate context window. I hated the process of doing this when my LLM use was heavy so I built this.

## Contents
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Command Options](#command-options)
- [Ignore Patterns](#ignore-patterns)
  - [Default Ignores](#default-ignores)
  - [Custom Ignores](#custom-ignores)
- [Output Format](#output-format)
- [Common Use Cases](#common-use-cases)
- [License](#license)

## Features

- 🚀 Fast recursive codebase scanning
- 📝 Markdown-formatted output
- 📋 Direct clipboard support
- 🎯 Smart file filtering:
  - Skips binary files
  - Ignores sensitive content (env files, keys, etc.)
  - Excludes common build artifacts
- 🎛 Flexible ignore patterns:
  - Built-in defaults for common ignore patterns
  - Support for `.codesnap_ignore` file
  - Additional CLI ignore flags
- Token count (Note: This is based off tokenization used for gpt-4o, adjust accordingly for other models)

## Installation

```bash
# Install directly with Go
go install github.com/wyattcupp/codesnap@latest

# Or clone and build
git clone https://github.com/wyattcupp/codesnap.git
cd codesnap
go build
```

## Usage
```bash
# Basic usage (outputs to clipboard)
codesnap -c

# Save to file instead
codesnap -o codebase.md

# Scan specific directory
codesnap -d /path/to/project -c

# Add additional ignore patterns
codesnap -i "tests/" -i "*.tmp" -c
```

### Command Options
```bash
Flags:
  -c, --clipboard       Copy output to clipboard
  -d, --dir string      Target directory to scan (default: current directory)
  -h, --help           Help for codesnap
  -i, --ignore strings  Additional patterns to ignore
  -o, --out string     Output file path for markdown
  ```

## Ignore Patterns
### Default Ignores
The tool comes with sensible defaults to ignore (to see the complete list, see the [code](/internal/collector/collector.go)):

- Version control directories (.git/, .svn/)
- Build artifacts (dist/, build/, *.exe)
- Dependencies (node_modules/, vendor/)
- Sensitive files (.env, *.key, credentials)
- Binary and media files
- Common metadata files

### Custom Ignores
Create a `.codesnap_ignore` file in your project root:
```.gitignore
# Ignore specific files
secret.txt
*.log

# Ignore directories
temp/
local/
```

The syntax follows `.gitignore` conventions:

- Use `#` for comments
- `*` matches any string except `/`
- `**` matches zero or more directories
- Trailing `/` matches directories only
- Leading `!` negates the pattern

## Output Format
The tool generates markdown with each file, either copied to clipboard or to an output file, with format:

````
<codebase>

**path/to/file.ext**

```

// your file contents

```

</codebase>
````

The corresponding terminal output shows a success message, along with token count:
```
Token count: {total tokens in output}
```

## License
MIT, use away


