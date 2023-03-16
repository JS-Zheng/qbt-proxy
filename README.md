# qbt-proxy

qbt-proxy is a lightweight proxy server for qBittorrent that simplifies the process of connecting your web applications to your qBittorrent Web API. It is designed to handle cross-origin requests more easily by returning the SID (session ID) as JSON in the response body during login, instead of using the `Set-cookie` HTTP header. qbt-proxy also supports accessing other APIs by including the SID as a query parameter in the request.

By using qbt-proxy, you can add an extra layer of authentication, provide a unified access point for multiple qBittorrent instances, and benefit from support for both HTTP and HTTPS with optional TLS configuration.

## Features

- Proxy requests from your web application to your qBittorrent Web API.
- Support for both HTTP and HTTPS with optional TLS configuration.
- Return SID (session ID) as JSON in the response body during login, instead of using the `Set-cookie` HTTP header, to facilitate cross-origin requests.
- Support accessing other APIs by including the SID as a query parameter in the request.

## Requirements

- Go 1.16 or later
- qBittorrent 4.x with Web UI enabled

## Installation

1. Clone the repository:

```sh
git clone https://github.com/JS-Zheng/qbt-proxy.git
```

2. Build the project:

```sh
cd qbt-proxy
go build
```

## Usage

1. Create a configuration file named `config.yml` in the project directory or at `$XDG_CONFIG_HOME/qbt-proxy`. Refer to the `config.example.yml` file for an example.

2. Start the qbt-proxy:

```sh
./qbt-proxy
```

3. Access the qBittorrent Web UI through the qbt-proxy using the configured HTTP or HTTPS address.

## Configuration

The following configuration options are available:

- `debug`: Enables or disables debug mode for the entire application. Default: `false`.
- `http_port`: The HTTP port number to use for the server. Default: `9487`.
- `https_port`: The HTTPS port number to use for the server. Optional. Default: `0` (disabled).
- `tls_cert`: The path to the TLS certificate file. Optional.
- `tls_key`: The path to the TLS private key file. Optional.
- `base_url`: The base URL of the qBittorrent instance. Required.

### Environment Variables

You can also configure qbt-proxy using environment variables. The following variables can be used:

- `QBP_DEBUG`: Equivalent to the `debug` configuration option.
- `QBP_HTTP_PORT`: Equivalent to the `http_port` configuration option.
- `QBP_HTTPS_PORT`: Equivalent to the `https_port` configuration option.
- `QBP_TLS_CERT`: Equivalent to the `tls_cert` configuration option.
- `QBP_TLS_KEY`: Equivalent to the `tls_key` configuration option.
- `QBP_BASE_URL`: Equivalent to the `base_url` configuration option.

## Docker

A `Dockerfile` is provided for running qbt-proxy in a Docker container. To build the Docker image, run:

```sh
docker build -t qbt-proxy .
```

To run the qbt-proxy in a Docker container, use the following command:
```sh
docker run -d -p 9487:9487 -v /path/to/config.yml:/app/config.yml qbt-proxy
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.

