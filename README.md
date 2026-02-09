<img width="400" src="./assets/kubesnap_logo.svg"/>

### Kubesnap: Improved kubernetes context & namespace management tool.

[![Latest Release](https://img.shields.io/github/v/release/hunsy9/kubesnap?label=release)](https://github.com/hunsy9/kubesnap/releases/latest)
[![License](https://img.shields.io/github/license/hunsy9/kubesnap)](https://github.com/hunsy9/kubesnap/blob/main/LICENSE)

**kubesnap** (`ks`) is a TUI tool that provides an **instant cluster status overview** and fast switching between cluster contexts/namespaces, offering a streamlined way to manage your Kubernetes environment.

<img width="800" src="./assets/kubesnap.gif" alt="kubesnap demo"/>

## Features

- **Context Switching**: Fast, fuzzy-searchable cluster context selector.
- **Edit Contexts**: Rename or Delete contexts directly within the TUI.
- **Namespace Switching**: Interactive namespace switcher with a `ks ns ~` shortcut for default namespace.
- **Cluster Dashboard**: Real-time overview of current connection and resource status (Nodes, Pods, Events).


## Installation

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

### Commands
- `ks`: View current cluster connection and health overview.
- `ks ctx`: Open the context switcher.
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