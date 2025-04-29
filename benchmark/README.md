# Benchmark

## Required tools

- [tcpkali](https://github.com/ssinyagin/tcpkali-debian) 
- [k6](https://grafana.com/docs/k6/latest/)

## Usage 

### Install
```shell
bo build -o benchmark
```

### TCP
choose a server machine to run, host such as 192.168.100.120.

```shell
./benchmark --kind=tcp --mode=server
```

choose a client machine to run, host such as 192.168.100.121.

```shell
./benchmark --kind=tcp --mode=client --host=192.168.100.120 --out=.
```

### HTTP
choose a server machine to run, host such as 192.168.100.120.

```shell
./benchmark --kind=http --mode=server
```

choose a client machine to run, host such as 192.168.100.121.

```shell
./benchmark --kind=http --mode=client --host=192.168.100.120 --out=.
```