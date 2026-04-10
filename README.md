# Go-sos (Save Our Services)

Gosos is a simple CLI tool for monitoring website and API statuses.

I needed something I could use to check the multiple self-hosted services I use.

## Features

- Add and remove URLs for monitoring
- List all registered URLs
- Check the status of all URLs at once
- Real-time monitoring with customizable intervals
- User-friendly command-line interface

## Installation

To install gosos, make sure you have Go installed on your system, then run:

```
go install git.thrls.net/thiagorls/gosos@latest
```

## Usage

Gosos provides several commands for managing and monitoring URLs:

```
gosos <command> [options]
```

### Commands

- `add <url>`: Add a URL to the monitoring list
- `remove <url|index>`: Remove a URL from the monitoring list, by full URL or by its index as shown in `gosos list`
- `list`: Display all registered URLs
- `run`: Check the status of all registered URLs once
- `live [interval]`: Start monitoring all URLs in real-time
    - `[interval]`: Optional check interval in seconds (default: 30)
- `help`: Show the help message

### Examples

```
gosos add https://example.com
gosos remove https://example.com
gosos remove 0       # Remove the first URL in the list
gosos list
gosos run
gosos live 60        # Check every 60 seconds
```

## Configuration

Gosos stores the list of URLs in a JSON file located at `~/.gosos-urls.json`. This file is automatically created by the tool.

## Dependencies

- [pterm](https://github.com/pterm/pterm): For terminal output styling and live updates

## Contributing

This is a project I built with the set of features I need for my personal use, but any contributions are welcome! Please feel free to submit a Pull Request.

