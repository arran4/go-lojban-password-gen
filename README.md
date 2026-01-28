# Go Lojban Password Generator

This project provides a tool to generate Lojban sentences that can be used as passwords or passphrases. It randomly combines gismu (root words) and cmavo (structure words) to create unique sequences.

## Installation

You can install the tool using `go install`:

```bash
go install github.com/arran4/go-lojban-password-gen/cmd/jbopwdgen@latest
```

## Usage

To run the password generator, you need the standard Lojban word lists: `gismu.txt` and `cmavo.txt`.

These files are expected to be in the current directory, or in the directory specified by the `DICTIONARY_DIR` environment variable.

```bash
jbopwdgen -gismu /path/to/gismu.txt -cmavo /path/to/cmavo.txt
```

### Options

- `-gismu`: Path to the `gismu.txt` file (default: `$DICTIONARY_DIR/gismu.txt` or `./gismu.txt`).
- `-cmavo`: Path to the `cmavo.txt` file (default: `$DICTIONARY_DIR/cmavo.txt` or `./cmavo.txt`).
- `-minsize`: Minimum number of words in the generated sentence (default: 5).

### Example Output

```text
Generated Random Lojban sequence:
prami mi 42 gerku

Sentence Components:
prami: love
mi: I/me
gerku: dog
```

## Dictionary Files

The tool expects the Lojban dictionary files in a specific fixed-width format (legacy format).

- **gismu.txt**: Lines must be at least 157 characters long. The word is expected at columns 1-6 (0-indexed).
- **cmavo.txt**: Lines must be at least 63 characters long. The word is expected at columns 0-11.

You can often find these files in Lojban archives (e.g., legacy definitions).

## License

See [LICENSE](LICENSE) file.
