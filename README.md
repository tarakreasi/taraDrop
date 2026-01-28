# taraDrop

**taraDrop** is a simple, single-binary local file sharing tool. It was born from a specific need: **transferring video from an iPhone to a Linux laptop without a cable.**

It is now a universal tool that works for **Android**, **iOS**, **Windows**, **Linux**, and **MacOS** users. Because it uses standard HTTP, you can transfer files from **any device** (Phone or Laptop) to the host computer just by opening a web browser.

## The Story

I needed to copy a video from an iPhone to my Linux laptop.
*   I am not an iPhone user.
*   I didn't have a cable with me.
*   Bluetooth transfer relies on specific stacks.
*   I just wanted it done.

I initially solved this with a Python script. However, I wanted a tool that I could easily share with Windows users something that "just works" without requiring them to install Python, manage dependencies, or open a terminal.

**Why switch to Golang?**
The goal was a **single, portable executable**.
*   **Windows Experience**: Users just double-click `taraDrop.exe`. No terminal window (silent mode), no configuration.
*   **Zero Dependencies**: Everything is compiled in.

## Key Features

*   **Smart Port Selection (User Friendly)**:
    *   The app anticipates that non-technical users shouldn't have to type port numbers.
    *   It defaults to **Port 80** so you can just type the IP (e.g., `http://192.168.1.15`).
    *   If Port 80 is unavailable or permission is denied, it automatically falls back to **Port 5000**.
*   **Smart IP Detection**: The app intelligently scans your network interfaces to display the correct LAN IP address you are actually using.
*   **Simple GUI**: A minimal control panel to see the URL, open the uploads folder, or stop the server.

## Installation & Usage

### Windows
1.  Download `taraDrop.exe`.
2.  Double-click the file.
3.  A window will appear showing your address (e.g., `http://192.168.1.5`).
4.  Open that address on your mobile phone browser to start uploading.

### Linux
1.  Download or build `taraDrop`.
2.  Run the binary:
    ```bash
    ./taraDrop
    ```
    *Note: To potentially use Port 80, you may need `sudo`, otherwise it will default to Port 5000.*

## Development

**Prerequisites:** Go 1.21+

### Build from Source

```bash
# Clone the repository
git clone https://github.com/tarakreasi/taraDrop.git
cd taraDrop

# Build for Linux
make linux

# Build for Windows (from Linux using MinGW)
make windows
```

## Tech Stack

*   **Language**: Golang (selected for cross-compilation and single-binary distribution).
*   **GUI**: Fyne (a lightweight Go toolkit for building native apps).

## Author & Repository

*   **GitHub**: [tarakreasi/taraDrop](https://github.com/tarakreasi/taraDrop)
*   **Maintainer**: Tri Wantoro (ajarsinau@gmail.com)
*   **Web**: [tarakreasi.com](https://tarakreasi.com)
