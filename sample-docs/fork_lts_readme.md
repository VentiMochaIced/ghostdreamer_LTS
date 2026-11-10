# Ghostdreamer LTS — Go Host & Stream Engine (Fork Variant)

Welcome to the **Ghostdreamer LTS** (Long Term Stable) repository fork! This branch serves as an isolated, high-performance testing ground and standalone execution environment for the Ghostdreamer audio visualizer and stream control engine.

---

## ⚡ Architecture Overview

Ghostdreamer is engineered around a **Single-File Binary & Embedded Template** paradigm (`//go:embed`). Unlike legacy Content Management Systems that rely on heavy server-side databases and persistent user session trackers, Ghostdreamer provides absolute user privacy:
* **Zero-Persistence Local State:** All configuration parameters, custom themes, and visualizer presets run strictly in-memory or locally via `127.0.0.1:8080`.
* **Air-Gapped & Ephemeral:** No user telemetry, browsing history, or private feed configurations ever leave your local machine.
* **Multi-Theme & Layout Engine:** Switch seamlessly between modern dark mode, VLC orange interfaces, and YouTube-inspired aesthetics.

---

## 🚀 Quick Start & Local Execution

### Prerequisites
* Go 1.20 or higher installed on your machine.

### 1. Clone & Navigate
```bash
git clone -b ghostdreamer_LTS https://github.com/YOUR_ALT_USERNAME/ghostdreamer.git
cd ghostdreamer
```

### 2. Run Locally (`go run`)
```bash
go run main.go
```
Open your browser and navigate to: **`http://127.0.0.1:8080`**

---

## 🛠️ Compiling Standalone Binaries

To compile a 100% self-contained executable binary (with zero external template dependencies required at runtime):

### For Windows (Headless/GUI Optimized):
```powershell
go build -ldflags="-H windowsgui -s -w" -o ghostdreamerv1b.exe main.go
```

### For Linux / Servers:
```bash
go build -ldflags="-s -w" -o ghostdreamerv1b-linux main.go
```

---

## 🎨 Features & Customization
* **Selectable Layouts:** Toggle instantly between `default` modern dashboard and `vlc` orange control modes.
* **Live Audio Visualizer:** Cybernetic frequency bars, radial waves, and oscilloscope rendering loops powered by HTML5 Canvas & Web Audio API.
* **Auto-Update Check:** Built-in lightweight endpoint (`/api/version`) checking upstream GitHub releases for updates without intrusive tracking.
* **PWA Support:** Install directly as a native desktop application via Edge or Chrome (`manifest.json` included).

---

## 📄 License & Contributing
Maintained as an independent LTS fork variant. Contributions, bug reports, and theme variants are welcome via pull requests to this branch before upstream review.