# EasyWeather /Ecowitt weather station  prometheus exporter

Simple webserver to export EasyWeather / Ecowitt weather
station data into prometheus so you hook up a prometheus and grafana ontop of it.

Example dashboard: https://weather.linkmenu58.lt/

Note: Ecowitt weather station can "send" data only to port `80` and no TLS is
possible. So when configuring your station to `POST` data, make sure its deployed
with non-TLS URL! I know... no TLS in 202x...

Deployment using [synpse.NET](https://synpse.net):
```
name: easyweather
description: easyweather
scheduling:
  type: Conditional
  selectors:
    weather: "true"
spec:
  containers:
    - name: easyweather
      image: ghcr.io/mjudeikis/ecowitt-easyweather:0.0.9
      ports:
        - 9080:9080
      env:
        - name: SERVER_URI
          value: :9080
      restartPolicy: {}

```

Grafana & prometheus deployment using [synpse.NET](https://synpse.net):
```
name: weather-prometheus
description: Prometheus metrics weahter station
scheduling:
  type: Conditional
  selectors:
    monitoring: infra
spec:
  containers:
    - name: prometheus
      image: prom/prometheus
      user: root
      ports:
        - 19090:9090
      volumes:
        - /synpse/prometheus-weather/:/prometheus
      secrets:
        - name: prometheus-config
          filepath: /etc/prometheus/prometheus.yml
      restartPolicy: {}
    - name: grafana
      image: grafana/grafana
      user: root
      ports:
        - 13000:3000
      volumes:
        - /synpse/grafana-weather:/var/lib/grafana
      env:
        - name: GF_PATHS_CONFIG
          value: /etc/grafana/grafana2.ini
      secrets:
        - name: grafana-config
          filepath: /etc/grafana/grafana2.ini
      restartPolicy: {}
```

Where `prometheus-config` is file based secret:
```
# my global config
global:
  scrape_interval:     15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
  - static_configs:
    - targets:
      # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  # - "first_rules.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:

  - job_name: 'weather-api'
    scrape_interval: 30s
    scrape_timeout: 10s
    static_configs:
     - targets: ['weather-ingestor.example.lt']
    metrics_path: "/api/metrics"
    scheme: https
```

`grafana-config`:
```
[security]
admin_user=admin
admin_password=admin_password
disable_login_form=true

[auth.anonymous]
enabled=true
hide_version=true

[auth.basic]
enabled = true
```

Stack is exposed using cloudflare tunnels - https://synpse.net/blog/cloudflare/setting-up-cloudflare-argo-tunnels/
