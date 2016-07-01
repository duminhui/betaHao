package neuron

import (
	"math/rand"
)

const (
	MAX_OUTPUTS = 5
)

const (
	Quiet int64 = iota // be quiet when initialized
	Active
	Blocked
)

type Cell struct {
	base_p               float64
	excit_p              float64
	pool                 float64
	last_excit_timestamp int64
}

func (pl *Cell) Decrease() {
	pl.pool -= 0.2
}

func (pl *Cell) Recover(delta int64) {
	pl.pool = pl.pool + 0.1*float64(delta)
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

	state int64
	cell  Cell
	trans Transmission

	Excited bool // run-time tag, for inputs to mark
}

func (nn *Neuron) caculate_present_state() {
	// recover pool energy
	det_excit_step := step - nn.cell.last_excit_timestamp
	nn.cell.Recover(det_excit_step)
	// judge state
	det_trans_step := step - nn.trans.last_trans_timestamp

	if det_excit_step < 5 {
		nn.state = Blocked
	} else if det_trans_step < 3 {
		nn.state = Active
	} else {
		nn.state = Quiet
	}
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
	if p < nn.cell.excit_p {
		return true
	} else {
		return false
	}
}

func (nn *Neuron) change_state(state int64) {
	nn.state = state
}

func (nn *Neuron) push_into_dequeue(nn *Neuron) {

}

func (this *Neuron) pass_potential(next *Neuron) bool {
	next.caculate_present_state()
	if next.in_blocking_period() {
		// if next neuron need recovered, then this neuron should not trans its excited and its excited_p should decrease
		this.trans.Decrease()
	} else {

		if next.in_activing_period() {
			temp_p := next.merge_probability(this.cell.excit_p)
			if next.try_enough_energy() {
				if next.try_excite() { // if could be excited,
					next.cell.Decrease() // equals try_avergy_pre_neurons
					this.trans.Increase()
					next.change_state(Blocked) // let next be into blocking_state
					return true
				} else { // want excite, but no engery
					next.cell.excit_p = temp_p
					// TODO: should there be a decrease of this.trans
				}
			} else { // not enough energy
				next.cell.Decrease()
			}
		} else { // in scilent state
			temp_p := next.cell.base_p
			if next.try_enough_energy() {
				next.cell.Decrease()
				this.trans.Increase()
				next.change_state(Blocked) // just to mark a timestamp
				return true
			} else {
				next.change_state(Active) // just to mark a timestamp
				// TODO: should there be a decrease of this.trans, maybe a distinguishing of growing up and mature is needed
			}
		}

	}
	return false
}
