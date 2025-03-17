# Benchmark

## Required tools

- [tcpkali](https://github.com/machinezone/tcpkali) 

## Usage 

### Install
```shell
bo build -o benchmark
```

### TCPKALI
choose a server machine to run, host such as 192.168.100.120.

```shell
./benchmark --mode=server
```

choose a client machine to run, host such as 192.168.100.121.

```shell
./benchmark --mode=tcpkali --host=192.168.100.120 --time=10s --count=50 --out=.
```
```shell
./benchmark --mode=tcpkali --host=192.168.100.120 --repeat=5000 --count=50 --out=.
```

### Local

```shell
./benchmark --mode=local
```