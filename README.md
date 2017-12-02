# Astrogoat

MySQL load simulation tool

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
