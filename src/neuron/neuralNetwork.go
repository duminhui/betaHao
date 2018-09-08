package neuron

import (
	"fmt"
	// "os"
	"math/rand"

	"github.com/eapache/queue"
)

const BRANCH_OF_EACH_NEURON int = 3

type NeuralNetwork struct {
	inputs  map[*Neuron]int
	outputs map[*Neuron]int

	Neurons []*Neuron
}

func (nk *NeuralNetwork) Generate_nodes(num int) {
	for i := 0; i < num; i++ {
		p := &Neuron{}
		p.index = i
		nk.Neurons = append(nk.Neurons, p)
	}
}

func (nk *NeuralNetwork) Add_edge(br *Branch, nn *Neuron) {
	br.next = nn
}

func (nk *NeuralNetwork) Generate_random_graph(n int) {

	for i := 0; i < len(nk.Neurons); i++ {
		for j := 0; j < BRANCH_OF_EACH_NEURON; j++ {

			for {
				next_i := rand.Intn(len(nk.Neurons))
				if next_i != i {
					nk.Add_edge(nk.Neurons[i].branches[j], nk.Neurons[next_i])
					fmt.Println(nk.Neurons[i])
					break
				}
			}

		}
	}

}

func (nk *NeuralNetwork) Generate_inputs(num int) {
	nk.inputs = map[*Neuron]int{}

	for len(nk.inputs) < num {
		i := rand.Intn(len(nk.Neurons))
		nk.inputs[nk.Neurons[i]] = i
	}
}

func (nk *NeuralNetwork) Generate_outputs(num int) {
	nk.outputs = map[*Neuron]int{}

	for len(nk.outputs) < num {
		i := rand.Intn(len(nk.Neurons))
		nk.inputs[nk.Neurons[i]] = i
	}
}

func (nk *NeuralNetwork) Read_outputs(learn_mode bool, expected_out []bool) (result []bool) {
	result = make([]bool, len(expected_out))
	if learn_mode == true {
		for np, idx := range nk.outputs {
			if np.register.state == Excited && expected_out[idx] == false {
				result[idx] = false
				np.register.branch.Decrease()
			} else {
				if np.register.state == Excited {
					result[idx] = true
				} else {
					result[idx] = false
				}
			}
		}
	} else {
		for np, idx := range nk.outputs {
			if np.register.state == Excited {
				result[idx] = true
			} else {
				result[idx] = false
			}
		}
	}
	return
}

func (nk *NeuralNetwork) Write_inputs() {

}

func (nk *NeuralNetwork) Boot_up(step int) {
	running_queue := queue.New()
	var neuron_deduplicator map[*Neuron]struct{}

	null_neuron := &Neuron{}
	running_queue.Add(null_neuron)

	var nn *Neuron

	for i := 0; i < step; i++ {
		nn = running_queue.Peek().(*Neuron)
		running_queue.Remove()

		if nn == null_neuron {

			neuron_deduplicator = make(map[*Neuron]struct{})
			action := nk.Read_outputs(false, expected_out)
			nk.Write_inputs()

			running_queue.Add(nn)
			global_step = global_step + 1

		} else if nn.register.state == Excited {
			for i := 0; i < len(nn.branches); i++ {

				result := nn.branches[i].touch()
				if result {
					if _, ok := neuron_deduplicator[nn]; !ok {
						running_queue.Add(nn)
					}
				}

			}

		}
	}

}

func Test1() {
	np := &Neuron{}
	fmt.Println(np)
}

func Test() {
	var sets map[*Neuron]int
	sets = make(map[*Neuron]int)
	cell1 := &Neuron{index: 1}
	sets[cell1] = 1

	cell1.index = 2

	sets[cell1] = 2

	cell3 := &Neuron{index: 3}
	sets[cell3] = 3

	fmt.Println(cell1.index, sets[cell1])

	// delete(sets, cell1)
	for k, v := range sets {
		println(k, v)
	}

	delete(sets, cell3)

	for k, v := range sets {
		println(k, v)
	}

	q := queue.New()
	q.Add(cell1)
	q.Add(cell3)

	var t *Neuron
	t = q.Peek().(*Neuron)
	fmt.Println(*t)
	q.Remove()
	t = q.Peek().(*Neuron)
	fmt.Println(*t)
	q.Remove()

	fmt.Println(q.Length())
}
