<img width="400" src="./assets/kubesnap_logo.svg"/>

### Kubesnap: Improved kubernetes context & namespace management tool.

<a href="https://pkg.go.dev/github.com/hunsy9/kubesnap"><img src="https://pkg.go.dev/badge/github.com/hunsy9/kubesnap.svg" alt="Go Reference"></a>
<a href="https://github.com/hunsy9/kubesnap/releases/latest"><img src="https://img.shields.io/github/v/release/hunsy9/kubesnap?label=release" alt="Latest Release"></a>
<a href="https://github.com/hunsy9/kubesnap/actions"><img src="https://github.com/hunsy9/kubesnap/actions/workflows/release.yml/badge.svg" alt="Build Status"></a>

**kubesnap** (`ks`) is a TUI tool that provides an **instant cluster status overview** and **fast switching** between cluster contexts/namespaces, offering a streamlined way to manage your Kubernetes environment.

<img width="800" src="./assets/main.gif" alt="kubesnap demo"/>

## Features

- **Cluster Dashboard**: Real-time overview of current connection and resource status (Nodes, Pods, Events).
- **Context Switching**: Fast, fuzzy-searchable cluster context selector.
- **Edit Contexts**: Rename or Delete contexts directly within the TUI.
- **Namespace Switching**: Interactive namespace switcher with a `ks ns ~` shortcut for default namespace.

## Get Started

### Homebrew (macOS & Linux)
```bash
brew tap hunsy9/kubesnap
brew install kubesnap
```

### Shell Script (Unix-like)
```bash
curl -sfL https://raw.githubusercontent.com/hunsy9/kubesnap/main/install.sh | sh
```

## Usage

### View Cluster Status

<img width="800" src="./assets/ks-cmd-demo.gif" alt="kubesnap demo"/>

- `ks`: View current cluster connection and health overview.

### Switch/Rename/Delete Contexts

<img width="800" src="./assets/ks-ctx-cmd-demo.gif" alt="kubesnap demo"/>

- `ks ctx`: Open the context switcher.
    - `'r' key`: Context Rename mode
    - `'d' key`: Context Delete mode

### Switch Namespace

<img width="800" src="./assets/ks-ns-cmd-demo.gif" alt="kubesnap demo"/>

- `ks ns`: Open the namespace switcher.
- `ks ns ~`: Quickly switch to the `default` namespace.

### Key Bindings
- `↑ / ↓`: Navigate list
- `/`: Search / Filter
- `Enter`: Select item
- `r`: Rename context
- `d`: Delete context
- `q / Esc`: Quit / Cancel

## Maintenance

### Update
Run the install script again or use `brew upgrade kubesnap`.

### Uninstallation
```bash
# Via Script
curl -sfL https://raw.githubusercontent.com/hunsy9/kubesnap/main/uninstall.sh | sh

# Via Homebrew
brew uninstall kubesnap
```

---
*Created by [Seunghun](https://github.com/hunsy9) using [Bubble Tea](https://github.com/charmbracelet/bubbletea).*