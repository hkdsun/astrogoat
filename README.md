# Astrogoat

Astrogoat is a MySQL load simulation tool. It provides a sophisticated testing rig based on docker containers.

# Components

* Replication cluster: simple master and slave percona mysql configuration
* Toxiproxy: configurable network condition between the master and slave
* [`pt-heartbeat`](https://www.percona.com/doc/percona-toolkit/LATEST/pt-heartbeat.html): used to monitor replication lag
* Prometheus: metric collection - currently configured to scrape replication lag from the pt-heartbeat table (`meta.heartbeat`)
* Grafana: frontend to Prometheus, allows for viewing replication lag in realtime

# Development

Start the test replication cluster:

```
docker-compose up -d
```

Install required gems and run the setup script:

```
bundle
./script/setup-repl
```

Toxiproxy settings can be tuned in `./script/setup-repl` and then updated with

```
./script/setup-repl toxiproxy
```
