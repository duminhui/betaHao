# betaHao

betaHao is a program to implement a new idea on how the neuron works

## Installation

```bash
git clone https://github.com/duminhui/betaHao
```

## Dependencies & Installation

`golang 1.3.3`: golang compilier

> `gccgo` may be also available. I'm not sure.

[`lane`](https://github.com/oleiade/lane): an implementation of queue struct which is used to record dynamic working neurons

```bash
$ sudo apt-get install golang
$ cd betaHao/ && export GOPATH=`pwd`
$ go get github.com/oleiade/lane
```

## Compliling

```bash 
$ go build
```

## Usage

```bash
$ ./betaHao
```

## Contributing

Any help will be appreciated, including translating, coding, wiking, and so on.

More explain about the idea and the algorithm is still working on [duminhui.cn](http://duminhui.cn/post/article/neuron-simulation)(Chinese Only, and unfinished)
