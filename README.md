
# Lightweight Kafka benchmark tool

The goal of this tool is to efficiently saturate Kafka broker without being the bottleneck itself.

This tool is heavily based on [franz-go](https://github.com/twmb/franz-go/tree/master/examples/bench) benchmark example.

You would need [podman](https://podman.io/) installed to build images. See [images/Makefile](images/Makefile) for example configuration

The tool can be run using make

```
make build
make run
```

The command above will autmatically create the required topic (named `test`) and produce data to it. It expects that your broker is accessible as `localhost:9092`