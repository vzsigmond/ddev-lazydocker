# LazyDocker DDEV Addon

This addon integrates [LazyDocker](https://github.com/jesseduffield/lazydocker) into your DDEV project, providing a simple and interactive TUI (Terminal User Interface) for managing Docker containers.

[![DDEV Add-on](https://img.shields.io/badge/DDEV-Add--on-blueviolet)](https://ddev.readthedocs.io/en/stable/users/extend/addon/)
![Repo Size](https://img.shields.io/github/repo-size/vzsigmond/ddev-lazydocker)
![Latest Tag](https://img.shields.io/github/v/tag/vzsigmond/ddev-lazydocker)
![CI](https://github.com/vzsigmond/ddev-lazydocker/actions/workflows/tests.yml/badge.svg)

## 📚 Table of Contents
- [Features](#-features)
- [Requirements](#-requirements)
- [Installation](#-installation)
- [Usage](#-usage)
- [Uninstallation](#-uninstallation)
- [Contributing](#-contributing)
- [Credits](#-credits)

## 📦 Features
- Easy-to-use interface for monitoring Docker containers and resources.
- Runs seamlessly within the DDEV ecosystem.
- Project-specific data: only shows containers, images, volumes, and networks related to your current DDEV project.
- One-liner to open LazyDocker via `ddev lazydocker`

## 🧰 Requirements
**DDEV** installed and set up on your system.

## 🚀 Installation

From your project directory:

```bash
ddev add-on get vzsigmond/ddev-lazydocker
```

```bash
ddev restart
```

## 💻 Usage:

```bash
ddev lazydocker
```

## 🧼 Uninstallation

```bash
ddev add-on remove vzsigmond/ddev-lazydocker
```

## 🤝 Contributing

Feel free to submit issues or pull requests if you encounter bugs or have suggestions for improvement. Contributions are always welcome!

## 🙌 Credits
- [Jesse Duffield](https://github.com/jesseduffield/) – Creator of LazyDocker
- DDEV Team – For the excellent local development tool

Enjoy a smoother Docker workflow! 🐳✨