package main
import (
    "bufio"
    "strconv"
    "strings"
    "fmt"
    "os/exec"
)

type ALE struct {
    weight int64
    height int64
    screen_list
    extern_command
    reader
    stdin
    stdout
}

func (ale *ALE) Init() {
    ale_exec_file := []string{ "./ale", "-game_controller", "fifo", "-display_screen", "true", "Breakout.bin" }
    ale.extern_command = exec.Command(ale_exec_file)

    ale.stdin, err = extern_command.StdinPipe()
    if err != nil {
        fmt.Println(err)
    }

    ale.stdout, err = extern_command.StdoutPipe()
    if err != nil {
        fmt.Println(err)
    }

    ale.reader = bufio.NewReader(stdout)

    ale.exetern_command.Start()
    ale.pull_command.Wait()

    line, _, err := ale.reader.ReadLine()

    if err != nil {
        fmt.Println(err)
    }

    temp := strings.Split(string(line), "-")
    ale.height, _ = strconv.ParseInt(temp[0], 10, 64)
    ale.width, _ = strconv.ParseInt(temp[1], 10, 64)

    _, err = stdin.Write([]byte("1,0,0,1\n"))

    ale.screen_list = make([]int64, ale.height*width)

}

func (ale *ALE) Final() {
    ale.stdin.Close()
    ale.stdout.Close()
}

func (ale *ALE) read_screen_state() (screen_list []int64, is_terminated int64, is_scored int64) {
    line, _, err  = ale.ReadLine()
    temp = strings.Split(string(line), ":")

    ptr := int64(0)
    lenth := len(temp[0])
    for i:=0; i<lenth; i+=4 {
        colour, _ := strconv.ParseInt((temp[0][i:i+2]), 16, 64)
        length, _ := strconv.ParseInt((temp[0][i+2:i+4]), 16, 64)

        screen_size := ale.height*ale.width
        for j:=ptr; (ptr<screen_size)&&(ptr<j+length); ptr++ {
            screen_list[ptr] = colour
        }
    }

    episode_string := srings.Split(string(temp[1]), ",")
    is_terminated, _ := strconv.ParseInt(episode_string[0], 10, 64)
    is_scored, _ := strconv.ParseInt(episode_string[1], 10, 64)

    return
}
