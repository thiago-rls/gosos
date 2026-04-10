package utils

const HelpText = `gosos (Save Our Services) - A simple cli tool and API status monitor

Usage: gosos <command> [options]

Commands:
  add <url>              Add a URL to the monitoring list
  remove <url|index>     Remove a URL from the monitoring list, by full URL
                         or by its index as shown in 'gosos list'
  list                   Display all registered URLs
  run                    Check the status of all registered URLs once
  live [interval]        Start monitoring all URLs in real-time
                         [interval]: Optional check interval in seconds (default: 60)
  help                   Show this help message

Examples:
  gosos add https://example.com
  gosos remove https://example.com
  gosos remove 0         (Remove the first URL in the list)
  gosos list
  gosos run
  gosos live 30          (Check every 30 seconds)

For more information, visit: https://git.thrls.net/thiagorls/gosos
`
