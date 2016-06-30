package neuron

import (
	"math/rand"
)

const (
	MAX_OUTPUTS = 5
)

type Cell struct {
	base_p               float64
	excit_p              float64
	pool                 float64
	last_excit_timestamp int64
}

func (pl *Cell) Decrease() {
	pl.pool -= 0.1
}

func (pl *Cell) Recover(delta int64) {
	pl.pool = pl.pool + 0.1 * float64(delta)
}

type Transmission struct {
	p                    float64
	last_trans_timestamp int64
}

func (ts *Transmission) Decrease() {
	ts.p -= 0.1
}

func (ts *Transmission) Increase() {
	ts.p += 0.1
}

type Neuron struct {
	// collect each pointers of the predecessors(neurons)
	Pre_neurons []*Neuron
	// collect each pointers of the successors(neurons)
	Post_neurons []*Neuron

	Excited bool
	cell    Cell
	trans   Transmission
}

func (nn *Neuron) caculate_next_neuron_present_status() {
	det_step := step - nn.cell.last_excit_timestamp
	nn.cell.Recover(det_step)
}

func (nn *Neuron) in_blocking_period() bool {
	delta_step := step - nn.cell.last_excit_timestamp
	if delta_step < 5 {
		return true
	} else {
		return false
	}
}

func (nn *Neuron) in_activing_period() bool {
	delta_step := step - nn.trans.last_trans_timestamp
	if delta_step < 3 {
		return true
	} else {
		return false
	}
}

func (nn *Neuron) merge_probability(trans_p float64) (p float64) {
	p = trans_p + nn.cell.excit_p
	return
}

func (nn *Neuron) try_enough_energy() bool {
	if nn.cell.pool < 0.1 {
		return false
	} else {
		return true
	}
}

func (nn *Neuron) try_excite() bool {
	 r := rand.New(rand.NewSource(16))
	 p := r.Float64()
	 if (p < nn.cell.excit_p) {
		return true
	} else {
		return false
	}
}

func (nn *Neuron) change_state() {

}

func (nn *Neuron) try_avergy_pre_neurons() {

}

func (nn *Neuron) pass_potential() {
	this := nn
	next := caculate_next_neuron_present_status()
	if next.in_blocking_period() {
		// if next neuron need recovered, then this neuron should not trans its excited and its excited_p should decrease
		this.trans.Decrease()
	} else {

		if next.in_activing_period() {
			temp_p := next.merge_probability()
			if next.try_enough_energy() {
				if next.try_excite() { // if could be excited,
					next.cell.Decrease()
					this.trans.Increase()
					push_next_neuron_into_dequeue()
					next.change_state(blocking) // let next be into blocking_state
				} else { // want excite, but no engery
					next.cell.excit_p = temp_p
					// TODO: should there be a decrease of this.trans
				}
			} else { // not enough energy
				next.try_avergy_pre_neurons()
			}
		} else { // in scilent state
			temp_p := next.cell.base_p
			if next.try_enough_energy() {
				next.cell.pool.Decrease()
				this.trans.Increase()
				push_next_neuron_into_dequeue()
				next.change_state(blocking) // just to mark a timestamp
			} else {
				next.change_state(activing) // just to mark a timestamp
				// TODO: should there be a decrease of this.trans, maybe a distinguishing of growing up and mature is needed
			}
		}

	}
}
