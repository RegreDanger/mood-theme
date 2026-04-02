# MoodTheme 🎨🎵

A high-performance local daemon written in Go that dynamically synchronizes your VS Code color theme based on the song currently playing in Spotify. Utilizing a Fan-In HTTP architecture and atomic file mutation, MoodTheme connects Spicetify with your local OS to create an immersive, music-driven coding environment.

## 📖 Table of Contents

- [MoodTheme 🎨🎵](#moodtheme-)
  - [📖 Table of Contents](#-table-of-contents)
  - [🏗️ Architecture Diagram](#️-architecture-diagram)
  - [🌟 Features](#-features)
  - [📦 Installation](#-installation)
    - [The Engine (Go Daemon)](#the-engine-go-daemon)
    - [The Client (Spicetify Extension)](#the-client-spicetify-extension)
  - [🚀 Quick Start](#-quick-start)
    - [1. Configure your Themes](#1-configure-your-themes)
    - [2. Start the Engine](#2-start-the-engine)
  - [📚 Documentation](#-documentation)
    - [CLI Flags](#cli-flags)
    - [JSON Structure](#json-structure)
  - [🏗️ Architecture](#️-architecture)
    - [Atomic Mutation \& RWMutex](#atomic-mutation--rwmutex)
    - [Profile Broadcasting](#profile-broadcasting)
  - [🐛 Troubleshooting](#-troubleshooting)
  - [🤝 Contributing](#-contributing)
  - [📄 License](#-license)
  - [📞 Support](#-support)

## 🏗️ Architecture Diagram

``` mermaid
flowchart TD
    %% Client Layer
    subgraph Client [Client Side]
        S[Spotify + Spicetify] -- "songchange event" --> J(JS Hook)
    end

    %% Backend Layer
    subgraph Daemon [MoodTheme Go Daemon]
        J -- "HTTP POST /theme" --> API{Go REST API}
        API <-->|"Thread-Safe Query"| Cache((RAM Cache<br>w/ RWMutex))
        Cache -. "Loads at startup" .-> DB[(themes.json)]
    end

    %% Output Layer
    subgraph Target [OS Environment]
        API -- "JSON Mutation Broadcast" --> VSC[VS Code Profiles]
    end

    %% Theming
    classDef default fill:#1e1e1e,stroke:#333,stroke-width:2px,color:#d4d4d4
    classDef spotify fill:#1db954,stroke:#191414,stroke-width:2px,color:#ffffff,font-weight:bold
    classDef goapi fill:#00add8,stroke:#000000,stroke-width:2px,color:#ffffff,font-weight:bold
    classDef vscode fill:#007acc,stroke:#000000,stroke-width:2px,color:#ffffff,font-weight:bold
    classDef db fill:#f39c12,stroke:#000000,stroke-width:2px,color:#ffffff,font-weight:bold

    %% Apply Styles
    class S spotify
    class API goapi
    class VSC vscode
    class DB,Cache db
```

## 🌟 Features

- **Zero-Friction Daemon**: Runs silently in the background capturing HTTP webhooks.
- **Atomic Data Mutation**: Protects your disk by keeping theme configurations in RAM using `sync.RWMutex` for thread-safe reads/writes.
- **Hot-Reload**: Update your song mappings on the fly without restarting the server.
- **Universal VS Code Profile Broadcast**: Automatically discovers and updates the `settings.json` across all your isolated VS Code profiles.
- **Built-in CORS & Anti-Spam Middleware**: Secures the local endpoint from browser preflight panic and spam requests.

## 📦 Installation

### The Engine (Go Daemon)

Download the pre-compiled binary for your operating system:

1. Go to the [Releases](https://www.google.com/search?q=https://github.com/RegreDanger/MoodTheme/releases "null") page.
2. Download the binary matching your OS (`moodtheme-windows.exe`, `moodtheme-mac`, `moodtheme-linux`).
3. Place it in a dedicated folder (e.g., `Documents/MoodTheme`).

### The Client (Spicetify Extension)

_Requires_ [_Spicetify_](https://spicetify.app/) _installed._

1. Copy the `mood-theme.js` file from this repository.
2. Paste it into your Spicetify extensions folder:
   - **Windows:** `%appdata%\spicetify\Extensions`
   - **Linux/Mac:** `~/.config/spicetify/Extensions`
   Run the following commands in your terminal:

```bash
spicetify config extensions mood-theme.js
spicetify apply
```

## 🚀 Quick Start

### 1. Configure your Themes

Create a `themes.json` file next to your downloaded binary. Map the exact Spotify song titles to your installed VS Code themes, here's an example and the project contains one.

```json
{
    "Themes": [
        {
            "theme_name": "Activate SCARLET protocol (beta)",
            "songs": [
                "After Dark",
                "Spoiler"
            ]
        },
        {
            "theme_name": "Kryptonite",
            "songs": [
                "Hearing Damage"
            ]
        }
    ]
}

```

### 2. Start the Engine

Run the binary from your terminal:

```bash
# Windows
./moodtheme-windows.exe

# Linux/Mac
./moodtheme-linux
```

Play a mapped song on Spotify, and watch your VS Code theme change instantly!

## 📚 Documentation

### CLI Flags

You can inject dependencies and configuration paths at startup using CLI flags:

| Flag | Default | Description |
| :--- | :--- | :--- |
| `-themes` | `./themes.json` | Absolute or relative path to your JSON mappings file. |

### JSON Structure

- `theme_name`: Must match the **exact ID** of the theme installed in VS Code.
- `songs`: Array of strings. Must match the **exact output** of the Spotify track metadata (case-sensitive).

## 🏗️ Architecture

### Atomic Mutation & RWMutex

To prevent Disk I/O bottlenecks when skipping tracks quickly, MoodTheme loads your `themes.json` into RAM at startup. It utilizes Go's `sync.RWMutex` to ensure thread-safe operations, allowing Spicetify to send hundreds of concurrent requests without race conditions or memory corruption.

### Profile Broadcasting

VS Code isolates settings across different profiles (`%APPDATA%\Code\User\profiles`). MoodTheme acts as a state synchronizer, scanning the OS directory tree to resolve the hashed profile paths and mutating every `settings.json` found to maintain a consistent environment.

## 🐛 Troubleshooting

| Error | Cause & Solution |
| :--- | :--- |
| **`OPTIONS /theme 404` (CORS)** | The Go daemon is not running. Start the MoodTheme engine _before_ opening Spotify. |
| **Theme doesn't change** | 1. Ensure the song name in `temas.json` matches Spotify exactly and ensure the `theme_name` is installed in VS Code. |
| **"Song doesn't have any theme"** | Working as intended. The current track is not mapped in your JSON file. Check the Spotify DevTools Console (`Ctrl + Shift + I`) for the exact song name. |

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the **PolyForm Noncommercial License 1.0.0.**

You may use, modify, and distribute this software for noncommercial purposes only.

Any commercial use, including but not limited to selling, sublicensing, or integrating this software into a paid product or service, is prohibited without explicit permission.

For commercial licensing inquiries, please contact me: [my email](mailto:carlosemiliogranadapererz@gmail.com)

See the `LICENSE` file for full legal details.

## 📞 Support

Built with ❤️ by [RegreDanger](https://www.google.com/search?q=https://github.com/RegreDanger "null")
