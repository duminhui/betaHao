package neuron

import (
	"fmt"
	"math/rand"
    "encoding/json"
    // "gopkg.in/fatih/set.v0"
)

const (
	MAX_OUTPUTS = 5
)

const (
	Quiet int64 = iota // be quiet when initialized
	// Active // cancel Active state, because it's the same as independent excite in multiple steps
    Hebb
	Blocked
)

type Cell struct {
	Base_p               float64
	float_p              float64
	pool                 float64
	last_excit_timestamp int64
}

func (pl *Cell) Decrease() {
	pl.pool -= 0.2
}

type Transmission struct {
	post_neurons []*Neuron
	p      []float64
}

func (trans Transmission) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)

    b = append(b, []byte(`{"Next":[`)...)
    for i:=0; i < len(trans.post_neurons); i++ {
        b = append(b, []byte(fmt.Sprintf("%d", trans.post_neurons[i].Key))...)
        if i !=len(trans.post_neurons)-1 {
            b = append(b, []byte(`,`)...)
        }
    }

    b = append(b, []byte(`],`)...)

    b = append(b, []byte(`"P":`)...)
    t, _ := json.Marshal(trans.p)
    b = append(b, t...)

    b = append(b, []byte(`}`)...)
    // fmt.Println(string(b))
	return b, nil
}

type Axon struct {
	// P                    float64
	last_trans_timestamp int64
	Trans                Transmission
}

func (ts *Axon) Decrease(i int) {
	ts.Trans.p[i] -= 0.01
	if ts.Trans.p[i] < -1 {
		ts.Trans.p[i] = -1
	}
}

func (ts *Axon) Increase(i int) {
    ts.Trans.p[i] += 0.01
    if (ts.Trans.p[i] > 1) {
        ts.Trans.p[i] = 1
    }
}

type Branch struct {
    neu *Neuron
    idx int
}

func (nn *Neuron) initial_branch() {
    nn.excited_neurons = make([]*Branch, 0)
}

func (nn *Neuron) add_to_branch(neu *Neuron,  idx int) {
    nn.excited_neurons = append(nn.excited_neurons, &Branch{neu, idx})
}

type Neuron struct {
	reversal_tag bool
	Key          int64
	// collect each pointers of the predecessors(neurons)
	pre_neurons []*Neuron
	// collect each pointers of the successors(neurons)
    excited_neurons []*Branch // run-time tag, for excited pre_neurons

	state int64
	Cell  Cell
	Axon  Axon

	Is_input bool
    Is_output bool
}

func (nn *Neuron) Init() {
	nn.Cell.Base_p = 0.5
	nn.Cell.pool = 1
}

/*
func (nn *Neuron) recover_energy() {
	det_excit_step := step - nn.Cell.last_excit_timestamp
	nn.Cell.Recover(det_excit_step)
}
*/

func (nn *Neuron) merge(trans_p float64) {
	nn.Cell.float_p = nn.Cell.float_p + trans_p
	return
}

func (nn *Neuron) try_enough_energy() bool {
	if nn.Cell.pool < 0.1 {
		return false
	} else {
		return true
	}
}

func (nn *Neuron) binarization() bool {
	r := rand.New(rand.NewSource(16))
	p := r.Float64()
	if p < nn.Cell.float_p {
		return true
	} else {
		return false
	}
}

func (nn *Neuron) change_state() {
    if(nn.state == Quiet) {
        if(len(nn.excited_neurons) >= 2) {
            nn.state = Hebb
            nn.Cell.last_excit_timestamp = step
        }
    }

    if(nn.state == Hebb) {
        if(nn.Cell.last_excit_timestamp < step) {
            nn.state = Blocked
            nn.Cell.last_excit_timestamp = step
        }
    }
    if(nn.state == Blocked) {
        if(nn.Cell.last_excit_timestamp < step - 3) {
            nn.state = Quiet
        }
    }
}

func (nn *Neuron) is_excited() (is_excited bool) {
    if(nn.state == Hebb) {
        for i := 0; i < len(nn.excited_neurons); i++ {
			p := nn.excited_neurons[i].neu
			idx := nn.excited_neurons[i].idx
            p.Axon.Increase(idx)
        }
        is_excited = true
    }

    if(nn.state == Quiet) {
        is_excited = nn.binarization()
    }

    if(nn.state == Blocked) {
        for i :=0; i < len(nn.excited_neurons); i++ {
			p := nn.excited_neurons[i].neu
			idx := nn.excited_neurons[i].idx
            p.Axon.Decrease(idx)
        }
        is_excited = true
    }
    return
}

func (nn *Neuron) pass_potential(next *Neuron, i int) {
	p := nn.Axon.Trans.p[i]
	next.merge(p)
}
