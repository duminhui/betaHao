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

//TODO: however pseudocode now
func (nn *Neuron) pass_potential() {
	this := nn
	next := caculate_next_neuron_present_status()
	if(next.in_resting_period()) {
		// if next neuron need recovered, then this neuron should not trans its excited and its excited_p should decrease
		this.trans.Decrease()
	} else {
		if(next.in_activing_period()) {
			temp_p = merge_probability()
			if(next.try_enough_energy()) {
				if (next.try_excite()) { // if could be excited,
					next.cell.pool.Decrease() // decrease next energy
					this.trans.Increase()
					push_next_neuron_into_dequeue() 
					change_next_neuron_state(resing_state) // let next be into resting_state
				} else {
					next.cell.excit_p = temp_p
					// TODO: should there be a decrease of this.trans
				}
			} else { // not enough energy,
				next.try_avergy_pre_neurons()
			}
		} else { // in scilent state
			temp_p = next.cell.base_p
			if(next.try_enough_energy() { 
				next.cell.pool.Decrease()
				this.trans.Increase()
				push_next_neuron_into_dequeue()
				change_next_neuron_state(resting_state) // just to mark a timestamp
			} else {
				change_next_neuron_state(activing_state) // just to mark a timestamp
				// TODO: should there be a decrease of this.trans, maybe a distinguishing of growing up and mature is needed
			}
		}
	}
}
