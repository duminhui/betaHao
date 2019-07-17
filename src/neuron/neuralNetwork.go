package neuron

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"

	"github.com/eapache/queue"
)

const BRANCH_OF_EACH_NEURON int = 3

type NeuralNetwork struct {
	inputs  []*Neuron
	outputs []*Neuron

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

	for len(nk.inputs) < num {
		i := rand.Intn(len(nk.Neurons))
		tmp_neuron := nk.Neurons[i]

		check_value := 0
		for j := 0; j < len(nk.inputs); j++ {
			if nk.inputs[j] == tmp_neuron {
				check_value = 1
			}
		}

		if check_value == 0 {
			nk.inputs = append(nk.inputs, tmp_neuron)
		}

	}
}

func (nk *NeuralNetwork) Generate_outputs(num int) {

	for len(nk.inputs) < num {
		i := rand.Intn(len(nk.Neurons))
		tmp_neuron := nk.Neurons[i]

		check_value := 0
		for j := 0; j < len(nk.inputs); j++ {
			if nk.outputs[j] == tmp_neuron {
				check_value = 1
			}
		}

		if check_value == 0 {
			nk.outputs = append(nk.outputs, tmp_neuron)
		}

	}
}

func (nk *NeuralNetwork) Read_outputs(learn_mode bool, expected_out []bool) (result []bool) {
	result = make([]bool, len(expected_out))

	if learn_mode == true { // 根据期望输出与实际输出作对比，完成输出神经元的学习

		for i := 0; i < len(nk.outputs); i++ {
			if nk.outputs[i].register.state == Excited && expected_out[i] == false {
				result[i] = false
				nk.outputs[i].register.branch.Decrease()
			} else {
				if nk.outputs[i].register.state == Excited {
					result[i] = true
				} else {
					result[i] = false
				}
			}
		}

	} else {

		for i := 0; i < len(nk.outputs); i++ {
			if nk.outputs[i].register.state == Excited {
				result[i] = true
			} else {
				result[i] = false
			}
		}

	}
	return
}

func (nk *NeuralNetwork) Write_inputs(RGB []byte, running_queue *queue.Queue) {
	var idx int64 = 0
	var offset int
	var bit int
	bin_buf := bytes.NewBuffer(RGB)
	for i := 0; i < len(RGB); i++ {
		offset = 1
		for j := 0; j < 8; j++ {
			var x int
			binary.Read(bin_buf, binary.BigEndian, &x)
			bit = x & offset
			if bit == 1 {
				running_queue.Add(nk.inputs[idx])
			}
			offset <<= 1
			idx++
		}
	}
}

func Get_cifar_data(ff *os.File) (expected_out []byte, vision_data []byte) {

	expected_out = make([]byte, 1)
	ff.Read(expected_out)
	// fmt.Printf("%d bytes: %d\n", n1, b1)
	vision_data = make([]byte, 3072)
	ff.Read(vision_data)
	return expected_out, vision_data
}

func (nk *NeuralNetwork) Boot_up(step int) {
	ff, _ := os.Open("./cifar-10-batches-bin/data_batch_3.bin")

	running_queue := queue.New()
	var neuron_deduplicator map[*Neuron]struct{}

	null_neuron := &Neuron{}
	running_queue.Add(null_neuron)

	var nn *Neuron
	var expected_out []byte
	var vision_data []byte
	for i := 0; i < step; i++ {
		nn = running_queue.Peek().(*Neuron)
		running_queue.Remove()

		if nn == null_neuron {

			neuron_deduplicator = make(map[*Neuron]struct{})

			expected_out, vision_data = Get_cifar_data(ff)

			action := nk.Read_outputs(false, expected_out)
			nk.Write_inputs(vision_data, running_queue)

			running_queue.Add(nn)
			global_step = global_step + 1

		} else if nn.register.state == Excited {
			for i := 0; i < len(nn.branches); i++ {

				result := nn.branches[i].touch()
				if result {
					if _, ok := neuron_deduplicator[nn]; !ok {
						running_queue.Add(nn) //?
					}
				}

			}

		}
	}

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
