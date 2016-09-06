package neuron

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	// "log"
	// "time"
	// "ALE"
	// "github.com/oleiade/lane"
	"math"
	"virtualEnvironment"

	"gopkg.in/fatih/set.v0"
	// "sync"
)

var step int64

type Input struct {
	// inputs           []*Neuron // 运行态
	mapping_relation map[int64]*Neuron
}

type Output struct {
	outputs          []int64 // 运行态
	mapping_relation map[*Neuron]int64
}

func (output *Output) clear() {
	output.outputs = make([]int64, 0)
}

type NeuralNetwork struct {
	Neurons     []*Neuron
	old_set     set.Interface // running-time binary
	current_set set.Interface // running-time binary
	env         Environmenter

	input             Input
	output            Output
	Input_number      []int64
	Output_number     []int64
	Num_of_controller int64
	Num_of_state      int64
}

type Environmenter interface {
	Init() (num_of_controller int64, num_of_state int64)
	Read_state() (screen_list []int64, is_terminated int64, is_scored int64)
	Write_action(outputs []int64)
	Final()
}

func (nk *NeuralNetwork) Generate_nodes(num int) {
	// initialize 'num' numbers of neurons in the network
	for i := 0; i < num; i++ {
		p := &Neuron{}
		p.Key = int64(i)
		p.Init()

		nk.Neurons = append(nk.Neurons, p)
	}
}

func (nk *NeuralNetwork) Add_edge(pre_neuron int, post_neuron int) {
	pre := nk.Neurons[pre_neuron]
	post := nk.Neurons[post_neuron]

	pre.Axon.Trans.post_neurons = append(pre.Axon.Trans.post_neurons, post)
	pre.Axon.Trans.p = append(pre.Axon.Trans.p, 0)
	post.pre_neurons = append(post.pre_neurons, pre)

}

func (nk *NeuralNetwork) Remove_edge(pre_neuron int, post_neuron int) {
	// TODO: is this necessary
}

func (nk *NeuralNetwork) Fast_generate_random_graph(n int, p float64, seed int64) {
	count := 0
	r := rand.New(rand.NewSource(seed))
	// r := rand.New(rand.NewSource(time.Now().UnixNano())
	v := 0
	w := -1
	lp := math.Log(1.0 - p)

	for v < n {
		lr := math.Log(1.0 - r.Float64())
		w = w + 1 + int(lr/lp)
		if v == w {
			w = w + 1
		}
		for v < n && n <= w {
			w = w - n
			v = v + 1
			if v == w {
				w = w + 1
			}
		}
		if v < n {
			nk.Add_edge(v, w)
			count++
		}
	}
	fmt.Println("ER graph generated", count, " of edges.")

}

func (nk *NeuralNetwork) Generate_inputs(num int64, seed int64) {
	r := rand.New(rand.NewSource(seed))
	nk.Input_number = make([]int64, 0, num)
	num_of_nodes := int64(len(nk.Neurons))
	for int64(len(nk.Input_number)) < num {
		input := r.Int63n(num_of_nodes)
		exist := false
		for _, v := range nk.Input_number {
			if v == input {
				exist = true
				break
			}
		}

		if !exist {
			nk.Input_number = append(nk.Input_number, input)
		}
	}
	fmt.Printf("len of input_order: %v \n", len(nk.Input_number))

	nk.input.mapping_relation = make(map[int64]*Neuron, num)
	for i, v := range nk.Input_number {
		nk.input.mapping_relation[int64(i)] = nk.Neurons[v]
		nk.Neurons[v].Is_input = true
	}

	return
}

func (nk *NeuralNetwork) Generate_outputs(num int64, seed int64) {
	r := rand.New(rand.NewSource(seed))
	nk.Output_number = make([]int64, 0, num)
	num_of_nodes := int64(len(nk.Neurons))
	var output int64

	for int64(len(nk.Output_number)) < num {
		output = r.Int63n(num_of_nodes)
		exist := false
		for _, v := range nk.Output_number {
			if v == output {
				exist = true
				break
			}
		}

		if !exist {
			nk.Output_number = append(nk.Output_number, output)
		}
	}

	nk.output.mapping_relation = make(map[*Neuron]int64, num)

	for i, v := range nk.Output_number {
		nk.output.mapping_relation[nk.Neurons[v]] = int64(i)
		nk.Neurons[v].Is_output = true
	}

	fmt.Println("nk.Output_number: ", nk.Output_number)
	fmt.Println("generate output mapping: ", nk.output.mapping_relation)

	return
}

func (nk *NeuralNetwork) Init() {
	// ale := ALE.ALE{}
	ale := virtualEnvironment.VirtualENV{}
	nk.env = &ale
	nk.Generate_nodes(20)

	nk.Fast_generate_random_graph(20, 0.3, 99)

	nk.Num_of_controller, nk.Num_of_state = nk.env.Init()
	fmt.Println("environment state: ", nk.Num_of_state)
	fmt.Println("environment controller: ", nk.Num_of_controller)

	nk.Generate_inputs(nk.Num_of_state, 10) // num_of_inputs, seed
	nk.Generate_outputs(nk.Num_of_controller, 100)
	nk.current_set = set.New()
	nk.old_set = set.New()
}

func (nk *NeuralNetwork) check_outputs() {
	nk.env.Write_action(nk.output.outputs)
}

func (nk *NeuralNetwork) Write_to(filename string) {
	reslt, _ := json.Marshal(nk)

	out, _ := os.Create(filename)
	defer out.Close()

	out.Write(reslt)
}

func (nk *NeuralNetwork) check_inputs() {
	screen_inputs, is_terminated, is_scored := nk.env.Read_state()
	merge_inputs := append(screen_inputs, is_terminated, is_scored)
	fmt.Println(" input list: ", merge_inputs)
	// TODO: merge_inputs or mapping_relation ?
	for i, v := range merge_inputs {
		if v > 0 {
			// fmt.Println("input: ", nk.input.mapping_relation[int64(i)])
			nk.current_set.Add(nk.input.mapping_relation[int64(i)])
		}
	}
}

func (nk *NeuralNetwork) check_if_outputs(neu *Neuron) {
	if idx, ok := nk.output.mapping_relation[neu]; ok {
		nk.output.outputs = append(nk.output.outputs, idx)
		fmt.Println("outputs nodes added.")
	}
}

func (nk *NeuralNetwork) finish_exciting_transmitting(neu interface{}) {
	if nn, ok := neu.(*Neuron); ok {
		for i, next := range nn.Axon.Trans.post_neurons {
			nn.pass_potential(next, i)
			// pass后需要将后端神经元去重，入set
			nk.current_set.Add(next)
		}
	}
}

func (nk *NeuralNetwork) Boot_up(step int) {
	// current_set := set.New()
	// old_set := set.New()

	for step > 0 {
		nk.check_inputs()

		for !nk.old_set.IsEmpty() {

			nn := nk.old_set.Pop().(*Neuron)
			nn.change_state()

			if nn.is_excited() {
				for i, next := range nn.Axon.Trans.post_neurons {
					nn.pass_potential(next, i)
					nk.current_set.Add(next)
				}

				if nn.Is_output {
					nk.output.outputs = append(nk.output.outputs, nk.output.mapping_relation[nn])
				}
			}
		}

		nk.check_outputs()

		nk.old_set.Clear()
		nk.output.clear()
		nk.old_set = nk.current_set
		nk.current_set.Clear()
		step--
	}
}
