# Permissions

The `fix-chainlink-tests` can run into lots of permissions issues when running in a sandboxed environment. Here are some common ones and fixes to implement.

## Sandbox Settings References

- [Claude Code](https://code.claude.com/docs/en/settings#sandbox-settings)
- [Gemini CLI](https://geminicli.com/docs/cli/sandbox/#overview-of-sandboxing)
- [Cursor](https://cursor.com/docs/reference/sandbox)
- [Codex](https://developers.openai.com/codex/concepts/sandboxing)

## Unable to Write/Read Go Cache

Make sure your filesystem can read/write to these paths:

- `~/Library/Application Support/go`
- `~/Library/Caches/go-build`
- `~/.asdf/installs/golang`

## Connection Denied for `localhost` or `::1`

Add these domains to allowed network connections:

- `localhost`
- `127.0.0.1`
- `::1`
