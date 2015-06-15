# docker-compose-graphviz

Turn a `docker-compose.yml` into a Graphviz `.dot` file. Currently in prototype state. For an example of the output,
check out https://github.com/abesto/abesto-net-docker.

## Installation

```sh
go install github.com/abesto/docker-compose-graphviz
```

## Usage

`cd` into a directory that has a `docker-compose.yml` file.

```sh
docker-compose-graphviz | dot -odocker-compose.jpg -Tjpg && open docker-compose.jpg
```