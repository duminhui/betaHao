package neuron

import (
	"math/rand"
)

const (
	MAX_BRANCHES = 3
)

const (
	Excited     int = 0
	Unexcited   int = 1
	UnInhibited int = 2
	Inhibited   int = 3
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
	branch *Branch
	// probability float64
	step  int
	state int
}

var global_step int = 0

type Neuron struct {
	index int
	//is_input  bool
	//is_output bool

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

func (br *Branch) binarization(p float64) {
	// p = br.probability + br.next.register.probability

	if p >= 0 {

		if p > 1 {
			p = 1
		}
		if p < rand.Float64() {
			br.next.register.state = Unexcited
		} else {
			br.next.register.state = Excited
		}

	} else {

		if p < -1 {
			p = -1
		}
		if p > -rand.Float64() {
			br.next.register.state = UnInhibited
		} else {
			br.next.register.state = Inhibited
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

func (br *Branch) register_to_next_neuron(nn *Neuron) { //, state int64) {
	br.next.register.branch = br
	// br.next.register.probability = br.probability
	br.next.register.step = global_step
	// br.next.register.state = state
}

func (br *Branch) touch() (result bool) {
	delta_step := global_step - br.next.register.step

	var probability float64

	if br.probability >= 0 {

		if delta_step == 0 {
			if br.next.register.state == Excited {
				probability = 1
			}
			if br.next.register.state == Unexcited {
				probability = 1
			}
			if br.next.register.state == UnInhibited {
				probability = br.probability + br.next.register.branch.probability
			}
			if br.next.register.state == Inhibited {
				probability = br.probability - 1
			}

			br.binarization(probability)
			br.impulse(br.next)
			br.register_to_next_neuron(br.next)
		} else if delta_step > 0 && delta_step <= 3 {
			if br.next.register.state == Excited {
				br.Decrease()
			}
			if br.next.register.state == Unexcited || br.next.register.state == UnInhibited || br.next.register.state == Inhibited {
				probability = br.probability + br.next.register.branch.probability
				br.binarization(probability)
				br.impulse(br.next)
				br.register_to_next_neuron(br.next)
			}
		} else {
			probability = br.probability
			br.binarization(probability)
			br.impulse(br.next)
			br.register_to_next_neuron(br.next)
		}

	} else if br.probability < 0 {

		if delta_step == 0 {
			if br.next.register.state == Excited {
				probability = br.probability + 1
			}
			if br.next.register.state == Unexcited {
				probability = br.probability + br.next.register.branch.probability
			}
			if br.next.register.state == UnInhibited {
				probability = -1
			}
			if br.next.register.state == Inhibited {
				probability = -1
			}

			br.binarization(probability)
			br.impulse(br.next)
			br.register_to_next_neuron(br.next)

		} else if delta_step > 0 && delta_step <= 3 {
			if br.next.register.state == Excited || br.next.register.state == Unexcited || br.next.register.state == UnInhibited {
				probability = br.probability + br.next.register.branch.probability
				br.binarization(probability)
				br.impulse(br.next)
				br.register_to_next_neuron(br.next)
			}
			if br.next.register.state == Inhibited {
				br.Increase()
			}
		} else {
			probability = br.probability
			br.binarization(probability)
			br.impulse(br.next)
			br.register_to_next_neuron(br.next)
		}

	}

	if br.next.register.state == Excited {
		result = true
	} else {
		result = false
	}

	return result
}
