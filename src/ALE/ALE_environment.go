package ALE

import (
	"bufio"
	"fmt"
	// "github.com/willf/bitset"
	"io"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Mask struct {
	x      int64
	y      int64
	radius int64
}

type ALE struct {
	height int64
	width  int64

	mask Mask

	screen_list          []int64
	binarized_screen     []bool
	avaliable_controller []string
	config               map[string]int64

	extern_command *exec.Cmd
	stdin          io.WriteCloser
	stdout         io.ReadCloser
	reader         *bufio.Reader
	err            error
}

func (ale *ALE) connect_to_the_controller() (num_of_controller int64) {
	ale.config = map[string]int64{
		"A_NOOP":            0,
		"A_FIRE":            1,
		"A_UP":              2,
		"A_RIGHT":           3,
		"A_LEFT":            4,
		"A_DOWN":            5,
		"A_UP_RIGHT":        6,
		"A_UP_LEFT":         7,
		"A_DOWN_RIGHT":      8,
		"A_DOWN_LEFT":       9,
		"A_UP_FIRE":         10,
		"A_RIGHT_FIRE":      11,
		"A_LEFT_FIRE":       12,
		"A_DOWN_FIRE":       13,
		"A_UP_RIGHT_FIRE":   14,
		"A_UP_LEFT_FIRE":    15,
		"A_DOWM_RIGHT_FIRE": 16,
		"A_DOWN_LEFT_FIRE":  17,
		"B_NOOP":            18,
		"B_FIRE":            19,
		"B_UP":              20,
		"B_RIGHT":           21,
		"B_LEFT":            22,
		"B_DOWN":            23,
		"B_UP_RIGHT":        24,
		"B_UP_LEFT":         25,
		"B_DOWN_RIGHT":      26,
		"B_DOWN_LEFT":       27,
		"B_UP_FIRE":         28,
		"B_RIGHT_FIRE":      29,
		"B_LEFT_FIRE":       30,
		"B_DOWN_FIRE":       31,
		"B_UP_RIGHT_FIRE":   32,
		"B_UP_LEFT_FIRE":    33,
		"B_DOWM_RIGHT_FIRE": 34,
		"B_DOWN_LEFT_FIRE":  35,
		"RESET":             40,
		"SAVE_STATE":        43,
		"LOADE_STATE":       44,
		"SYSTEM_RESET":      45}

	ale.avaliable_controller = []string{
		"A_NOOP",
		"A_FIRE",
		"A_RIGHT",
		"A_LEFT",
		"SYSTEM_RESET"}
	/*
	   ale.Output_to_controller = make(map[*Neuron]int, len(avlaliabe_controller))
	   for i := 0; i <len(output_to_controller); i++ {
	       ale.Output_to_controller[outputs[i]] = config[avaliable_controller[i]]
	   }
	*/

	num_of_controller = int64(len(ale.avaliable_controller))

	return
}

func (ale *ALE) Init() (num_of_controller int64, num_of_state int64) {
	ale.mask.x = 0
	ale.mask.y = 0
	ale.mask.radius = 10
	// ale.extern_command = exec.Command("./ale", "-game_controller", "fifo", "-display_screen", "true", "Breakout.bin")
	ale.extern_command = exec.Command("./bin/ale", "-game_controller", "fifo", "Breakout.bin")

	ale.stdin, ale.err = ale.extern_command.StdinPipe()
	if ale.err != nil {
		fmt.Println(ale.err)
	}

	ale.stdout, ale.err = ale.extern_command.StdoutPipe()
	if ale.err != nil {
		fmt.Println("stdout: ", ale.err)
	}
	// fmt.Println("stdin: ", ale.stdin)
	// fmt.Println("stdout: ", ale.stdout)

	ale.reader = bufio.NewReader(ale.stdout)

	ale.extern_command.Start()

	line, _, err := ale.reader.ReadLine()

	if err != nil {
		fmt.Println(err)
	}

	temp := strings.Split(string(line), "-")
	_ = "breakpoint"

	ale.height, _ = strconv.ParseInt(temp[0], 10, 64)
	ale.width, _ = strconv.ParseInt(temp[1], 10, 64)

	_, err = ale.stdin.Write([]byte("1,0,0,1\n"))

	fmt.Println("height: ", ale.height)
	fmt.Println("width: ", ale.width)

	ale.screen_list = make([]int64, ale.height*ale.width)
	fmt.Println("len of screen_list", len(ale.screen_list))

	num_of_controller = ale.connect_to_the_controller()
	num_of_state = 8*ale.height*ale.width + 2 //all the screen pixels, 8bits each pixels,

	return

}

func (ale *ALE) Final() {
	// ale.extern_command.Wait()
	ale.stdin.Close()
	ale.stdout.Close()
}

func submatrix(origin []int64, col, row, i, j, p, q int64) (result []int64) {
	result = make([]int64, p*q)
	k := 0
	for m := int64(0); m < p; m++ {
		for n := int64(0); n < q; n++ {
			if (i+m) < 0 || (i+m) >= row || (j+n) < 0 || (j+n) >= col {
				result[k] = 0
			} else {
				// fmt.Println(k, (i+m)*col+j+n)
				result[k] = origin[(i+m)*col+j+n]
			}
			k++
		}
	}
	return
}

func (ale *ALE) auxiliary_lens(image []int64, x int64, y int64, radius int64) (subimage []int64) {

	x_0, y_0 := x-radius, y-radius
	col := 2 * radius
	row := 2 * radius

	subimage = submatrix(image, ale.width, ale.height, x_0, y_0, col, row)

	_ = "breakpoint"

	return
}

func (ale *ALE) binarize(screen_list []int64) (binarized_screen []bool) {
	lenth := ale.height * ale.width
	ale.binarized_screen = make([]bool, lenth)

	for i := int64(0); i < lenth; i++ {
		if ale.screen_list[i] > 0 {
			ale.binarized_screen[i] = true
		} else {
			ale.binarized_screen[i] = false
		}
	}
	binarized_screen = ale.binarized_screen
	return
}

func (ale *ALE) Read_state() (screen_list []int64, is_terminated int64, is_scored int64) {
	line, _, _ := ale.reader.ReadLine()
	temp := strings.Split(string(line), ":")

	ptr := int64(0)
	lenth := len(temp[0])
	for i := 0; i < lenth; i += 4 {
		colour, _ := strconv.ParseInt((temp[0][i : i+2]), 16, 64)
		length, _ := strconv.ParseInt((temp[0][i+2 : i+4]), 16, 64)

		screen_size := ale.height * ale.width
		for j := ptr; (ptr < screen_size) && (ptr < j+length); ptr++ {
			ale.screen_list[ptr] = int64(colour)
		}
	}

	// screen_list = ale.binarize(ale.screen_list)
	screen_list = ale.auxiliary_lens(ale.screen_list, ale.mask.x, ale.mask.y, ale.mask.radius)

	episode_string := strings.Split(string(temp[1]), ",")
	// screen_list = ale.screen_list
	is_terminated, _ = strconv.ParseInt(episode_string[0], 10, 64)
	is_scored, _ = strconv.ParseInt(episode_string[1], 10, 64)

	return
}

func (ale *ALE) central_mask_point(mask_influences []int64) (is_changed bool) {
	x_t, y_t := ale.mask.x, ale.mask.y
	for _, v := range mask_influences {
		switch v {
		case 0:
			ale.mask.x += 1
		case 1:
			ale.mask.x -= 1
		case 2:
			ale.mask.y += 1
		case 3:
			ale.mask.y -= 1
		}
	}

	if ale.mask.x < 0 {
		ale.mask.x = 0
	} else if ale.mask.x > ale.height {
		ale.mask.x = ale.height - 1
	} else if ale.mask.y < 0 {
		ale.mask.y = 0
	} else if ale.mask.y > ale.width {
		ale.mask.y = ale.width - 1
	}

	_ = "breakpoint"

	if x_t == ale.mask.x && y_t == ale.mask.y {
		return false
	} else {
		return true
	}

}

func (ale *ALE) write_to_game(game_operator []int64) {
	fmt.Println("n: ", len(game_operator))
	n := len(game_operator)
	j := int64(rand.Intn(n))
	idx := game_operator[j]
	fmt.Println("avaliable_controller: ", ale.avaliable_controller[idx])
	out := ale.config[ale.avaliable_controller[idx]]
	result := fmt.Sprintf("%d, %d\n", out, 18)
	fmt.Println("result: ", result)

	_, err := ale.stdin.Write([]byte(result))
	if err != nil {
		fmt.Println(err)
	}
}

func (ale *ALE) Write_action(game_operator []int64, mask_influences []int64) {
	fmt.Println("game_operator: ", game_operator)
	rand.Seed(int64(time.Now().Nanosecond()))

	if len(game_operator) > 0 {
		ale.write_to_game(game_operator)
	} else {
		rand := int64(rand.Intn(len(ale.avaliable_controller)))
		ale.write_to_game(append(game_operator, rand))
	}

	if len(mask_influences) > 0 {
		ale.central_mask_point(mask_influences)
	} else {
		rand := int64(rand.Intn(5))
		ale.central_mask_point(append(mask_influences, rand))
	}

}
