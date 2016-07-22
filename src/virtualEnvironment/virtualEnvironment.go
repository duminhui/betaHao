package virtualENV

import ()

type virtualENV struct {
	num_of_controller int64
	num_of_state int64
}

func (ve *virtualENV) Init() (ve.num_of_controller int64, ve.num_of_state int64) {
	ve.num_of_controller = 5
	ve.num_of_state = 10
	return
}

func (ve *virtualENV) Read_state() (screen_list []int64, is_terminated int64, is_scored int64) {
	
