package ALE
import (
    "neuron"
    "bufio"
    "strconv"
    "strings"
    "fmt"
    "os/exec"
    "math/rand"
    "time"
)

type ALE struct {
    weight int64
    height int64
    screen_list
    extern_command
    Output_to_controller
    reader
    stdin
    stdout
}

func (ale *ALE) connect_to_the_controller(outputs []*Neuron) {
    config := map[string]int{
        "A_NOOP": 0,
        "A_FIRE": 1,
        "A_UP": 2,
        "A_RIGHT": 3,
        "A_LEFT": 4,
        "A_DOWN": 5,
        "A_UP_RIGHT": 6,
        "A_UP_LEFT": 7,
        "A_DOWN_RIGHT": 8,
        "A_DOWN_LEFT": 9,
        "A_UP_FIRE": 10,
        "A_RIGHT_FIRE": 11,
        "A_LEFT_FIRE": 12,
        "A_DOWN_FIRE": 13,
        "A_UP_RIGHT_FIRE": 14,
        "A_LEFT_FIRE": 15,
        "A_DOWM_RIGHT_FIRE": 16,
        "A_DOWN_LEFT_FIRE": 17,
        "B_NOOP": 18,
        "B_FIRE": 19,
        "B_UP": 20,
        "B_RIGHT": 21,
        "B_LEFT": 22,
        "B_DOWN": 23,
        "B_UP_RIGHT": 24,
        "B_UP_LEFT": 25,
        "B_DOWN_RIGHT": 26,
        "B_DOWN_LEFT": 27,
        "B_UP_FIRE": 28,
        "B_RIGHT_FIRE": 29,
        "B_LEFT_FIRE": 30,
        "B_DOWN_FIRE": 31,
        "B_UP_RIGHT_FIRE": 32,
        "B_LEFT_FIRE": 33,
        "B_DOWM_RIGHT_FIRE": 34,
        "B_DOWN_LEFT_FIRE": 35,
        "RESET": 40,
        "SAVE_STATE": 43,
        "LOADE_STATE": 44,
        "SYSTEM_RESET": 45,}

    avaliable_controller := []string{
        "A_NOOP",
        "A_FIRE",
        "A_RIGHT",
        "A_LEFT",
        "SYSTEM_RESET"}

    ale.Output_to_controller = make(map[*Neuron]int, len(avlaliabe_controller))
    for i := 0; i <len(output_to_controller); i++ {
        ale.Output_to_controller[outputs[i]] = config[avaliable_controller[i]]
    }
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

func (ale *ALE) write_action(excited_outputs_list []*Neuron) {
    num := len(excited_outputs_list)    
    rand.Seed(int64(time.Now().Nanosecond()))
    i := rand.Intn(num)
    result := ale.Output_to_controller[i] + ",18"

    _, err = stdin.Write([]byte(result))
}
