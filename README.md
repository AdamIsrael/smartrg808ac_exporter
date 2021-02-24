[![CircleCI](https://circleci.com/gh/AdamIsrael/smartrg808ac_exporter.svg?style=svg)](https://circleci.com/gh/adamisrael/smartrg808ac_exporter)
[![](https://goreportcard.com/badge/github.com/AdamIsrael/smartrg808ac_exporter)](https://goreportcard.com/report/github.com/AdamIsrael/smartrg808ac_exporter)
# SmartRG 808AC exporter for Prometheus

Work in progress


## Usage

```bash
SMARTRG_HOSTNAME=192.168.0.1 SMARTRG_USERNAME=admin SMARTRG_PASSWORD=<your password> ./smartrg808ac_exporter
```

## Development

You'll need to have [Go](https://golang.org/) installed to build the binary.

```bash
$ make build
go build -o smartrg808ac_exporter
```

## Environment Variables

- `SMARTRG_HOSTNAME` The hostname or IP address of your cable modem.
- `SMARTRG_USERNAME` The username to authenticate as. You can find this on the bottom of the modem.
- `SMARTRG_PASSWORD` The password to authenticate with. You can find this on the bottom of the modem.

## Systemd

```bash
# Create a user to run the service under
sudo useradd --no-create-home --shell /bin/false smartrg808ac_exporter

# Create the systemd service unit
sudo tee /etc/systemd/system/smartrg808ac_exporter.service <<"EOF"
[Unit]
Description=SmartRG 808AC Exporter

[Service]
User=smartrg808ac_exporter
Group=smartrg808ac_exporter
EnvironmentFile=-/etc/default/smartrg808ac_exporter
ExecStart=/usr/local/bin/smartrg808ac_exporter $OPTIONS
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and start the smartrg808ac_exporter service
sudo systemctl daemon-reload && \
sudo systemctl start smartrg808ac_exporter && \
sudo systemctl status smartrg808ac_exporter && \
sudo systemctl enable smartrg808ac_exporter
```

## History

| Description | Date | Version |
| ----------- | ---- | ------- |
| Initial release | 23 Feb 2020 | 0.1 |