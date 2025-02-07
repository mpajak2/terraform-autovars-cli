# Terraform-autovars-cli

terraform-autovars-cli is a simple CLI tool written in Go that helps your team run Terraform commands (`plan`, `apply`, `output`) from within a stack directory. It automatically detects the current stack based on your working directory and gathers Terraform variable files from environment-specific directories. It also supports decryption of secured variable files using `sops`.

## Features

- **Automatic Stack Detection**: Detects the stack name from the current working directory.
- **Environment Support**: Accepts an environment parameter (e.g., `dev`, `prod`) to locate variable files.
- **Variable File Collection**: Searches for JSON variable files in `ENVVARS/<env>` and `ENVVARS/<env>/secured`.
- **Secured File Decryption**: Uses `sops` to decrypt secured variable files.
- **Terraform Integration**: Passes the variable files as `-var-file` arguments to Terraform commands.

## Project Structure

An example project structure:

```
project/
├── stacks/
│   ├── 05-datadog/    <-- Run pops from this directory
│   └── another-stack/
└── ENVVARS/
    ├── dev/
    │   ├── 05-datadog.json
    │   ├── another-stack.json
    │   └── secured/
    │       ├── 05-datadog.json    (encrypted with sops)
    └── prod/
        ├── 05-datadog.json
        └── another-stack.json
```

## Prerequisites

- Go (for building the tool)
- Terraform
- `sops` (for decrypting secured variable files)

## Installation

### Build from Source

Clone the repository:

```bash
git clone https://github.com/yourusername/pops.git
cd pops
```

Build the executable:

#### On Linux/Mac:
```bash
go build -o pops main.go
```

#### On Windows (using PowerShell):
```powershell
go build -o pops.exe main.go
```
## Optionally get executable from release
- there is pops.exe file in the release to download

### Add to your PATH (optional):

- **Windows**: Move `pops.exe` to a folder already in your `PATH` or add its folder to the `PATH` environment variable.
- **Linux/Mac**: Move the `pops` binary to a directory like `/usr/local/bin` or add its folder to your `PATH`.

## Usage

Navigate to a stack directory. For example, if you are working in `project/stacks/05-datadog`, open your terminal in that directory.

Run the command:

```bash
pops plan dev
```

Where:
- `plan`: Terraform action (can be `plan`, `apply`, or `output`).
- `dev`: Environment (e.g., `dev` or `prod`).

### The tool will:

1. Detect the stack name (in this case, `05-datadog`).
2. Search for variable files in `ENVVARS/dev` and `ENVVARS/dev/secured` (decrypting any secured files using `sops`).
3. Pass the found files as `-var-file` arguments to the Terraform command.

