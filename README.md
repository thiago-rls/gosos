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
- `remove <url>`: Remove a URL from the monitoring list
- `list`: Display all registered URLs
- `run`: Check the status of all registered URLs once
- `live [interval]`: Start monitoring all URLs in real-time
    - `[interval]`: Optional check interval in seconds (default: 30)
- `help`: Show the help message

### Examples

```
gosos add https://example.com
gosos remove https://example.com
gosos list
gosos run
gosos live 60  # Check every 60 seconds
```

## Configuration

Gosos stores the list of URLs in a JSON file located at `~/.gosos-urls.json`. This file is automatically created by the tool.

## Development

To contribute to gosos or set up the development environment:

1. Clone the repository:
   ```
   git clone https://git.thrls.net/thiagorls/gosos.git
   ```
2. Navigate to the project directory:
   ```
   cd gosos
   ```
3. Install dependencies:
   ```
   go mod tidy
   ```
4. Build the project:
   ```
   go build gosos.go
   ```
5. Run the project:
   ```
   go run ./gosos
   ```

## Dependencies

- [pterm](https://github.com/pterm/pterm): For terminal output styling and live updates
- Standard Go libraries

## Contributing

This is a project I built with the set of features I need for my personal use, but any contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any problems or have any questions, please open an issue on the project repository.

## Todo
- [ ] Create tests
- [ ] Fix some glitches with pterm terminal output
- [ ] Add options to remove from list by index
