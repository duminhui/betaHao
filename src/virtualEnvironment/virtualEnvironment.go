package virtualEnvironment

import (
	"math/rand"
	"time"
)

type VirtualENV struct {
	num_of_controller int64
	num_of_state int64
}

func (ve *VirtualENV) Init() (num_of_controller, num_of_state int64) {
	ve.num_of_controller = 5
	ve.num_of_state = 10

	num_of_state = ve.num_of_state + 2
	num_of_controller = ve.num_of_controller

	return
}

func (ve *VirtualENV) Read_state() (screen_list []int64, is_terminated int64, is_scored int64) {
	screen_list = make([]int64, 0, ve.num_of_state)
	rand.Seed(int64(time.Now().Nanosecond()))
	// for random state
	var r float64
	for i:=int64(0); i<ve.num_of_state-int64(3); i++ {
		r = rand.Float64()
		if r < 0.5 {
			screen_list = append(screen_list, 1)
		} else {
			screen_list = append(screen_list, 0)
		}
	}
	// for fixed state
	for i:=0; i<3; i++ {
		screen_list = append(screen_list, 1)
	}

	is_terminated = 1
	is_scored = 1
	return
}

func (ve *VirtualENV) Write_action(game_controller []int64) {

}

func (ve *VirtualENV) Final() {

}
