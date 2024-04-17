# Palworld Exporter

Exports Palworld dedicated server metrics for prometheus.

## Prerequisites

Palworld exporter relies on [Palworld REST API](https://tech.palworldgame.com/ja/category/rest-api).
Set `RESTAPIEnabled=True` on your `PalWorldSettings.ini`.

## Usage

```
usage: palworld_exporter [<flags>]


Flags:
  -h, --[no-]help             Show context-sensitive help (also try --help-long and --help-man).
      --web.listen-address=":18212"  
                              Address to expose metrics. ($LISTEN_ADDRESS)
      --scrape_uri="http://localhost:8212/v1/api/metrics"  
                              URI to Palworld REST API metrics endpoint. ($SCRAPE_URI)
      --http_user="admin"     Username for Palworld REST API basic authentication. ($HTTP_USER)
      --http_password=HTTP_PASSWORD  
                              Password for Palworld REST API basic authentication. ($HTTP_PASSWORD)
      --log.level="info"      Only log messages with the given severity or above. If log.format is set 'default', this option is ignored. Valid levels: [debug,
                              info, warn, error] ($LOG_LEVEL)
      --log.format="default"  Output format of log messages. Valid formats: [default, text, json] ($LOG_FORMAT)
      --[no-]version          Show application version.
```

Options can be set both as command-line arguments and environments variables.

## Docker

```
docker run -e SCRAPE_URI=<your server endpoint> -e HTTP_PASSWORD=<your server password> -p 18212:18212 --rm jimmysharp/palworld_exporter
```

## Metrics

|Name|Description|Type|
|---|---|---|
|`palworld_up`|The status of the last scrape: `1` for success and `0` for failed|Gauge|
|`palworld_server_fps`|The server FPS|Gauge|
|`palworld_current_player_num`|The number of current players|Gauge|
|`palworld_server_frame_time`|Server frame time (ms)|Gauge|
|`palworld_max_player_num`|The maximum number of players|Gauge|
|`palworld_uptime`|The server uptime of seconds|Counter|

## License

[MIT](./LICENSE)