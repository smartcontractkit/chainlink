# Tools for configuring OCR2DR

Designed for OCR2DR project, to be able to fetch keys and IDs for different key bundles needed for OCR2DR, such as eth, p2p, ocr2.

## Usage

1. Create a hidden file containing remotes: hosts, logins and passwords. It is a simple text file, where each row has the following structure:

```
boot_node login password
node0 login0 password0
node1 login1 password1
...
```

Note: *bootstrap* node must go first.

Preparing such a file is the only "manual" step that is required. 

2. Review `templates` sub-folder

2. Run the tool: `go run . .myremotes`

The tool creates artefacts under `artefacts` sub-directory.

TBD
