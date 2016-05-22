# betaHao

betaHao is a program to implement a new idea on how the neuron works

## Installation

```bash
git clone https://github.com/duminhui/betaHao
```

## Dependencies & Installation

`golang 1.3.3`: golang compilier  
[`lane`](https://github.com/oleiade/lane): an implementation of queue struct which is used to record dynamic working neurons  
[`Arcade Learning Environment(ALE)`](https://github.com/mgbellemare/Arcade-Learning-Environment): a platform for simulating AI environment  
[`breakout`](https://atariage.com/2600/roms/Breakout.zip): a game ROM running on Atari 2600, for more other game available on ALE, please visit [AtariAge](https://atariage.com/system_items.html?SystemID=2600&ItemTypeID=ROM)


```bash
$ sudo apt-get install golang
$ cd betaHao/ && export GOPATH=`pwd`
$ go get github.com/oleiade/lane
$ mkdir temp && cd temp
$ sudo apt-get install libsdl1.2-dev libsdl-gfx1.2-dev libsdl-image1.2-dev cmake
$ git clone https://github.com/mgbellemare/Arcade-Learning-Environment
$ cd Arcade-Learning-Environment/
$ mkdir build && cd build
$ cmake -DUSE_SDL=ON -DUSE_RLGLUE=ON -DBUILD_EXAMPLES=ON ..
$ make -j 4
$ cd ../..
$ git clone https://github.com/davidljung/rl-glue
$ cd rl-glue/
$ ./configure
$ make
$ sudo make install
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
