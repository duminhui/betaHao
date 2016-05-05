# betaHao

betaHao is a program to implement a new idea on how the neuron works

## Dependencies

`go 1.3.3`

[`assert`](https://github.com/bmizerany/assert) is needed for developing

For Debian User:

```bash
$ sudo apt-get install golang
```
> `gccgo` may be also available. I'm not sure.

## Compliling

```bash 
$ git clone https://github.com/duminhui/betaHao
$ cd betaHao/
$ export GOPATH=`pwd`
$ go build
```

For developer

```bash
$ cd betaHao/
$ export GOPATH=`pwd`
$ go get github.com/bmizerany/assert
```

## Usage

```bash
$ ./betaHao
```

## Contributing

Any help will be appreciated, including translating, coding, wiking, and so on.

More explain about the idea and the algorithm is still working on [duminhui.cn](http://duminhui.cn/post/article/neuron-simulation)(Chinese Only, and unfinished)
