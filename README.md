# Astrogoat

Astrogoat is a MySQL replication lag throttling simulation tool. It provides a sophisticated testing rig based on docker containers so that you can easily experiment with different throttling strategies (e.g. tune a PID controller).

# Components

* Goat: the core logic. See `main.go` and `goat/` package
* Replication cluster: simple master and slave percona mysql configuration
* [Toxiproxy](https://github.com/shopify/toxiproxy): configurable network condition between the master and slave
* [`pt-heartbeat`](https://www.percona.com/doc/percona-toolkit/LATEST/pt-heartbeat.html): used to monitor replication lag
* [Prometheus](https://prometheus.io): metric collection - currently configured to scrape replication lag from the pt-heartbeat table (`meta.heartbeat`)
* [Grafana](https://grafana.com/): frontend to Prometheus, allows for viewing replication lag in realtime

# Usage

* Start the test replication cluster:

```
docker-compose up -d
```

* Install required gems and run the setup script:

```
bundle
./script/setup-repl
```

* Toxiproxy settings can be tuned in `./script/setup-repl` and then updated with

```
./script/setup-repl toxiproxy
```

* Try putting some load on the database

```
go run main.go
```

* Monitor replication lag with Grafana at `localhost:3000`

![](https://i.imgur.com/5W1XAR3.png)
