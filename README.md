<img width="400" src="./assets/kubesnap_logo.svg"/>

### Kubesnap: Improved kubernetes context & namespace management tool.

<a href="https://pkg.go.dev/github.com/hunsy9/kubesnap"><img src="https://pkg.go.dev/badge/github.com/hunsy9/kubesnap.svg" alt="Go Reference"></a>
<a href="https://github.com/hunsy9/kubesnap/releases/latest"><img src="https://img.shields.io/github/v/release/hunsy9/kubesnap?label=release" alt="Latest Release"></a>
<a href="https://github.com/hunsy9/kubesnap/actions"><img src="https://github.com/hunsy9/kubesnap/actions/workflows/release.yml/badge.svg" alt="Build Status"></a>

**kubesnap** is a TUI tool that provides an **instant cluster status overview** and **fast switching** between cluster contexts/namespaces, offering a streamlined way to manage your Kubernetes environment.

<img width="800" src="./assets/main.gif" alt="kubesnap demo"/>

## Features

- **Cluster Dashboard**: Real-time overview of current connection and resource status (Nodes, Pods, Events).
- **Context Switching**: Fast, fuzzy-searchable cluster context selector.
- **Edit Contexts**: Rename or Delete contexts directly within the TUI.
- **Namespace Switching**: Interactive namespace switcher with a `kubesnap ns ~` shortcut for default namespace.

## Get Started

### Homebrew (macOS & Linux)
```bash
brew tap hunsy9/kubesnap
brew install kubesnap
```

### Scoop (Windows)
```powershell
scoop bucket add hunsy9 https://github.com/hunsy9/scoop-bucket
scoop install kubesnap
```

### Shell Script (Unix-like)
```bash
curl -sfL https://raw.githubusercontent.com/hunsy9/kubesnap/main/install.sh | sh
```

## Usage

> **💡 Tip:** We recommend aliasing `kubesnap` to `ks` for speed.
> Add this to your shell config (`.bashrc`, `.zshrc`, etc.):
> ```bash
> alias ks='kubesnap'
> ```

### 1. Cluster Overview
View current cluster connection and health overview.

<img width="800" src="./assets/ks-demo.gif" alt="overview demo"/>

```bash
kubesnap
```

### 2. Switch/Rename/Delete Context
Open the interactive context switcher and edit your context.

<img width="800" src="./assets/ks-ctx-cmd-demo.gif" alt="context switch demo"/>

```bash
kubesnap ctx
```
- **`Enter`**: Switch to the selected context.
- **`r` key**: Rename the selected context.
- **`d` key**: Move to context deletion mode.

### 3. Switch Namespace

<img width="800" src="./assets/ks-ns-cmd-demo.gif" alt="namespace switch demo"/>

```bash
kubesnap ns
# or switch to default
kubesnap ns ~
```

### Key Bindings
- `↑ / ↓`: Navigate list
- `/`: Search / Filter
- `Enter`: Select item
- `r`: Rename context
- `d`: Delete context
- `q / Esc`: Quit / Cancel

## Maintenance

### Update
Run the install script again or use the following commands:

- **Homebrew**: `brew upgrade kubesnap`
- **Scoop**: `scoop update kubesnap`

### Uninstallation
```bash
# Via Homebrew
brew uninstall kubesnap

# Via Scoop
scoop uninstall kubesnap

# Via Script
curl -sfL https://raw.githubusercontent.com/hunsy9/kubesnap/main/uninstall.sh | sh
```

---
*Created by [Seunghun](https://github.com/hunsy9) using [Bubble Tea](https://github.com/charmbracelet/bubbletea).*