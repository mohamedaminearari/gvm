<p align="center">
    <img src="assets/GVMLogo.png" alt="GVM logo" width="250"/>
</p>

<h1 align="center">

Godot Version Manager

</h1>

<p align="center">
A cross-platform command-line tool for installing and managing multiple versions of the <a href="https://godotengine.org/">Godot game engine</a>.
</p>

<p align="center">
<span>
    <img src="https://img.shields.io/badge/Version%20-%201.26.4%20-%20%23478cbf?logo=go">
</span>
<span>
    <img src="https://goreportcard.com/badge/github.com/mohamedaminearari/gvm">
</span>
</p>


## Features
- **Cross-platform** — works on Windows, Linux, and macOS
- **Install any version** — download and install any stable or pre-release version of Godot directly from GitHub
- **Multiple versions side by side** — keep as many Godot versions installed as you need without conflicts
- **Easy version switching** — switch your active Godot version with a single command
- **Mono / C# support** — install and manage Godot Mono builds alongside standard builds
- **List installed versions** — see all locally installed versions at a glance
- **List remote versions** — browse all available Godot releases without leaving your terminal
- **No admin rights required** — everything lives under `~/.gvm`, no system-wide changes needed
- **Lightweight** — a single binary with no runtime dependencies
- **Fast** — written in Go for minimal overhead and quick installswd

## Installation
Download the Latest version from the [releases page](https://github.com/mohamedaminearari/gvm/releases), then extract it where ever you want, add the directory where you extracted it to the PATH as well as the ~/.gvm/bin directory.

## Usage
```
# List all available Godot versions for your OS
$ gvm ls-remote

# Install a specific Godot version
$ gvm install 4.3-stable

# Switch to an installed version
$ gvm use 4.3-stable

#Show the currently active version
$ gvm current

#List all locally installed versions
$ gvm ls

#Remove an installed version
$ gvm uninstall 4.3-stable

#Create a named alias for a version
$ gvm alias myproject 4.3-stable

#Show the path to a version's executable
$ gvm which 4.3-stable

#Launch a specific version directly
$ gvm run 4.3-stable
```

## Contributions
Pull requests are welcomed. For major changes, please open and issue first to discuss what you would like to change.

## License
[MIT](https://choosealicense.com/licenses/mit/)
