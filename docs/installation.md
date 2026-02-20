# Installation

## 1. Install Go

Hashdrop requires an up-to-date Go toolchain. It is designed to run on **Linux**, **macOS**, or **Windows via WSL**.

**Option 1 — Webi (Linux / WSL / macOS, recommended)**
```bash
curl -sS https://webi.sh/golang | sh
```

Read the output and follow any printed instructions before continuing.

**Option 2 — Official installer (any platform)**

Follow the [official Go installation instructions](https://go.dev/doc/install). On Windows this means downloading and running a `.msi` package.

After installing, open a new shell and confirm everything works:
```bash
go version
```

If you see a version string, move on to step 2.

<details>
<summary>Troubleshooting</summary>

- **Already have Go via Webi?** Re-run the same Webi command to update it.
- **"command not found" after install?** The Go binary's directory is probably not in your `PATH`. First, confirm where `go` lives — common locations are `~/.local/opt/go/bin` (Webi) or `/usr/local/go/bin` (official). Test with the full path, e.g. `~/.local/opt/go/bin/go version`, then add the directory to your `PATH`:
```bash
# Linux / WSL
echo 'export PATH=$PATH:$HOME/.local/opt/go/bin' >> ~/.bashrc
source ~/.bashrc
```
```bash
# macOS
echo 'export PATH=$PATH:$HOME/.local/opt/go/bin' >> ~/.zshrc
source ~/.zshrc
```

</details>

---

## 2. Install Hashdrop

Run the following command to download, build, and install the `hashdrop` binary:
```bash
go install github.com/anxhukumar/hashdrop/cli/cmd/hashdrop@latest
```

Verify the installation:
```bash
hashdrop --help
```

<details>
<summary>Troubleshooting</summary>

If you get a "command not found" error, Go's install directory (`$GOBIN`, which defaults to `$HOME/go/bin`) is likely not in your `PATH`. Add it:
```bash
# Linux / WSL
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
```
```bash
# macOS
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
source ~/.zshrc
```

Then try `hashdrop --help` again.

</details>

---

## 3. Create an account and log in


Run `hashdrop auth register` to create an account, then `hashdrop auth login` to authenticate.

You're all set — head over to the [Usage](./usage.md) to upload your first file.