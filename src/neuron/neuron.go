package neuron

import (
// "fmt"
)

const (
	MAX_OUTPUTS = 5
)

type Cell struct {
	base_p int64
	excit_p int64
	pool float64
	last_excit_timestamp int64
}

func (pl *Cell) Decrease() {
	pl.Pool -= 0.1
}


func (pl *Pool) Recover() {

}

type Transmission struct {
	p float64
	last_trans_timestamp int64
}

type Neuron struct {
	// collect each pointers of the predecessors(neurons)
	Pre_neurons []*Neuron
	// collect each pointers of the successors(neurons)
	Post_neurons []*Neuron

	Excited bool
	cell Cell
	trans Transmission
}

func (nn *Neuron) pass_potential() {
	this := nn
	next := caculate_next_neuron_present_status()
	if(next.in_resting_period()) {
		this.trans.Decrease()
	} else if(next.in_activing_period()) {
		merge_probability()
		next.try_excit()

	}
}
