package neuron

import (
	"math/rand"
)

const (
	MAX_BRANCHES = 3
)

const (
	Excited     int64 = 0
	Unexcited   int64 = 1
	UnInhibited int64 = 2
	Inhibited   int64 = 3
)

type Branch struct {
	next        *Neuron
	probability float64
}

func (branch *Branch) Increase() {
	branch.probability += 0.1
	if branch.probability > 1 {
		branch.probability = 1
	}
}

func (branch *Branch) Decrease() {
	branch.probability -= 0.1
	if branch.probability < -1 {
		branch.probability = -1
	}
}

type Register_of_previous_touch struct {
	branch      *Branch
	probability float64
	step        int64
	state       int64
}

var global_step int64 = 0

type Neuron struct {
	index     int64
	is_input  bool
	is_output bool

	register Register_of_previous_touch
	branches []*Branch
}

func (nn *Neuron) Increase() {
	for i := 0; i < len(nn.branches); i++ {
		nn.branches[i].Increase()
	}
}

func (nn *Neuron) Decrease() {
	for i := 0; i < len(nn.branches); i++ {
		nn.branches[i].Decrease()
	}
}

func (br *Branch) binarization(p float64) (result_state int64) {
	// p = br.probability + br.next.register.probability

	if p >= 0 {

		if p > 1 {
			p = 1
		}
		if p < rand.Float64() {
			result_state = Unexcited
		} else {
			result_state = Excited
		}

	} else {

		if p < -1 {
			p = -1
		}
		if p > -rand.Float64() {
			result_state = UnInhibited
		} else {
			result_state = Inhibited
		}
	}
}

func (br *Branch) impulse(nn *Neuron) {
	result_state := br.next.register.state
	if result_state == Excited {
		br.Increase()
		br.next.Increase()
	}
	if result_state == Unexcited {
		br.Decrease()
		if br.probability < 0 {
			br.next.register.state = UnInhibited
		}
		nn.register.branch.Decrease()
	}
	if result_state == UnInhibited {
		br.Increase()
		if br.probability > 0 {
			br.next.register.state = Unexcited
		}
		nn.register.branch.Decrease()
	}
	if result_state == Inhibited {
		br.Decrease()
		br.next.Decrease()
	}
}

func (br *Branch) register_to_next_neuron(nn *Neuron) {
	br.next.register.branch = br
	br.next.register.probability = br.probability
	br.next.register.step = global_step
}

func (br *Branch) touch() {
	delta_step := global_step - br.next.register.step
	var result_state int64
	if br.probability >= 0 {
		if delta_step == 0 {
			//if br.next.register.state == Excited {
			// do nothing
			//}
			if br.next.register.state == Unexcited {
				br.next.register.state == Excited
			}
			if br.next.register.state == UnInhibited {
				probability := br.probability + br.next.register.probability
				br.binarization(probability)
			}
			if br.next.register.state == Inhibited {
				probability := br.probability - 1
			}
		}
		if delta_step > 0 && delta_step <= 3 {

		}
	}
	if br.probability < 0 {

	}
}

func (nn *Neuron) run() {
	for i := 0; i < len(nn.branches); i++ {
		nn.branches[i].next
	}
}
