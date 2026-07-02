# Security Policy

## Supported versions

Only the latest release receives security fixes.

## Reporting a vulnerability

Please **do not** open a public issue for security problems. Instead, use
GitHub's private vulnerability reporting: on this repository, go to the
**Security** tab → **Report a vulnerability**. We will respond in the advisory
thread.

If the vulnerability is in PDF parsing, rendering, decryption, or another part
of the underlying library rather than the MCP layer, please report it to
[aspose-pdf-foss-for-go](https://github.com/aspose-pdf-foss/aspose-pdf-foss-for-go)
via the same private-reporting mechanism there.

## Scope notes

- This server executes **local file operations on behalf of an AI client**: it
  reads and writes files at the paths the client supplies. Run it only with
  MCP clients you trust, under a user account whose filesystem access matches
  what you are willing to grant the client.
- PDFs are untrusted input. The underlying library is designed to parse
  malformed files without hanging or crashing; a reproducible crash, hang, or
  memory blow-up on a crafted PDF is a security-relevant bug — please report
  it privately as described above.
