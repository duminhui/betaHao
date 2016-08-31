package neuron

import (
	// "fmt"
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
	Base_p               float64
	excit_p              float64
	pool                 float64
	last_excit_timestamp int64
}

func (pl *Cell) Decrease() {
	pl.pool -= 0.2
}

func (pl *Cell) Recover(delta int64) {
	pl.pool = pl.pool + 0.01*float64(delta)
	if pl.pool > 1 {
		pl.pool = 1
	}
}

type Transmission struct {
	post_neurons []*Neuron
	Trans_p      []float64
}

func (trans Transmission) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	for i := 0; i < len(trans.post_neurons); i++ {
		b = append([]byte("conn:"))
	}
	return b, nil
}

type Axon struct {
	P                    float64
	last_trans_timestamp int64
	Trans                Transmission
}

func (ts *Axon) Decrease() {
	ts.P -= 0.01
	if ts.P < -1 {
		ts.P = -1
	}
}

func (ts *Axon) Increase() {
	ts.P += 0.01
	if ts.P > 1 {
		ts.P = 1
	}
}

type Neuron struct {
	reversal_tag bool
	Key          int64
	// collect each pointers of the predecessors(neurons)
	pre_neurons []*Neuron
	// collect each pointers of the successors(neurons)
	// post_neurons []*Neuron

	state int64
	Cell  Cell
	Axon  Axon

	Excited bool // run-time tag, for inputs to mark
}

func (nn *Neuron) Init() {
	nn.Cell.Base_p = 0.5
	nn.Cell.pool = 1
}

func (nn *Neuron) recover_energy() {
	det_excit_step := step - nn.Cell.last_excit_timestamp
	nn.Cell.Recover(det_excit_step)
}

func (nn *Neuron) in_blocking_period() bool {
	if nn.state == Blocked {
		return true
	} else {
		return false
	}
}

func (nn *Neuron) in_activing_period() bool {
	if nn.state == Active {
		return true
	} else {
		return false
	}
}

func (nn *Neuron) in_quiet_period() bool {
	if nn.state == Quiet {
		return true
	} else {
		return false
	}
}

func (nn *Neuron) merge_probability(trans_p float64) (p float64) {
	p = trans_p + nn.Cell.excit_p
	return
}

func (nn *Neuron) try_enough_energy() bool {
	if nn.Cell.pool < 0.1 {
		return false
	} else {
		return true
	}
}

func (nn *Neuron) try_excite() bool {
	r := rand.New(rand.NewSource(16))
	p := r.Float64()
	if p < nn.Cell.excit_p {
		return true
	} else {
		return false
	}
}

func (nn *Neuron) change_state_from(state int64, is_excited bool) {
	det_excit_step := step - nn.Cell.last_excit_timestamp
	det_trans_step := step - nn.Axon.last_trans_timestamp

	if is_excited == true {
		nn.Cell.last_excit_timestamp = step
		nn.state = Blocked
	} else {

		if state == Blocked {
			if det_excit_step < 5 {
				nn.state = Blocked
			} else {
				nn.state = Quiet
			}
		}

		if state == Active {
			if det_trans_step < 3 {
				nn.state = Active
			} else {
				//TODO:??
			}
		}

		if state == Quiet {
			nn.state = Quiet
		}

	}
}

func (nn *Neuron) change_state_to(state int64) {
	nn.state = state
}

func (this *Neuron) pass_potential(next *Neuron) bool {
	next.recover_energy()
	if next.in_blocking_period() {
		this.Axon.Decrease()
		if next.Cell.last_excit_timestamp >= 5 {
			next.change_state_to(Quiet)
		}
	}

	if next.in_activing_period() {
		temp_p := next.merge_probability(this.Cell.excit_p)
		if next.try_enough_energy() {
			if next.try_excite() { // if could be excited,
				next.Cell.Decrease() // equals try_avergy_pre_neurons
				this.Axon.Increase()

				next.Cell.last_excit_timestamp = step
				next.change_state_to(Blocked) // let next be into blocking_state

				return true
			} else { // want excite, but no engery
				next.Cell.excit_p = temp_p
				this.Axon.last_trans_timestamp = step
			}
		} else { // not enough energy
			next.Cell.last_excit_timestamp = step
			next.change_state_to(Blocked)
		}
	}

	if next.in_quiet_period() { // in scilent state
		temp_p := next.Cell.Base_p
		if next.try_enough_energy() {
			next.Cell.Decrease()
			this.Axon.Increase()
			next.Cell.last_excit_timestamp = step
			next.change_state_to(Blocked) // just to mark a timestamp
			return true
		} else {
			next.Cell.excit_p = temp_p
			this.Axon.last_trans_timestamp = step
			next.change_state_to(Active) // just to mark a timestamp
			// TODO: should there be a decrease of this.trans, maybe a distinguishing of growing up and mature is needed
		}
	}

	return false
}
