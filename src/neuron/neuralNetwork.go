package neuron

import (
	"fmt"
	"math/rand"
	// "log"
	// "time"
	// "ALE"
	"virtualEnvironment"
	"github.com/oleiade/lane"
	"math"
	// "sync"
)

var step int64 = 0

type Input struct {
	inputs           []*Neuron // 运行态
	mapping_relation map[int64]*Neuron
}

type Output struct {
	outputs []*Neuron
	// game_operator   []*Neuron // 运行态
	// mask_influences  []*Neuron // 运行态
	mapping_relation map[*Neuron]int64
}

type NeuralNetwork struct {
	Neurons       []*Neuron
	Running_queue *lane.Queue
	env           Environmenter

	input             Input
	output            Output
	input_number      []int64
	output_number     []int64
	num_of_controller int64
	num_of_state      int64
}

type Environmenter interface {
	Init() (num_of_controller int64, num_of_state int64)
	Read_state() (screen_list []int64, is_terminated int64, is_scored int64)
	Write_action(outputs []int64)
	Final()
}

func (nk *NeuralNetwork) Generate_nodes(num int) {
	// initialize 'num' numbers of neurons in the network
	// nk.neurons = make([]*Neuron, num)
	// logger.Println(num, " of nodes.")
	for i := 0; i < num; i++ {
		p := &Neuron{}
		p.Init()

		nk.Neurons = append(nk.Neurons, p)
	}
}

func (nk *NeuralNetwork) Add_edge(pre_neuron int, post_neuron int) {
	pre := nk.Neurons[pre_neuron]
	post := nk.Neurons[post_neuron]

	pre.Post_neurons = append(pre.Post_neurons, post)
	post.Pre_neurons = append(post.Pre_neurons, pre)

}

func (nk *NeuralNetwork) Remove_edge(pre_neuron int, post_neuron int) {
	// v := 0
	// w := 0
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
	fmt.Println(count, " of edges.")
	fmt.Println("ER graph generated")

}

func (nk *NeuralNetwork) Generate_inputs(num int64, seed int64) {
	r := rand.New(rand.NewSource(seed))
	nk.input_number = make([]int64, num)
	num_of_nodes := int64(len(nk.Neurons))
	for int64(len(nk.input_number)) < num {
		input := r.Int63n(num_of_nodes)
		exist := false
		for _, v := range nk.input_number {
			if v == input {
				exist = true
				break
			}
		}

		if !exist {
			nk.input_number = append(nk.input_number, input)
		}
	}
	fmt.Printf("len of input_order: %v \n", len(nk.input_number))
	// fmt.Printf("len of inputs: %v \n", num)

	nk.input.mapping_relation = make(map[int64]*Neuron, num)
	for i, v := range nk.input_number {
		nk.input.mapping_relation[int64(i)] = nk.Neurons[v]
		// fmt.Printf("inpurt_order: %v, %T\n", i, i)
		// nk.Inputs[int64(i)] = nk.Neurons[v]
	}

	return
}

func (nk *NeuralNetwork) Generate_outputs(num int64, seed int64) {
	r := rand.New(rand.NewSource(seed))
	nk.output_number = make([]int64, num)
	num_of_nodes := int64(len(nk.Neurons))
	for int64(len(nk.output_number)) < num {
		output := r.Int63n(num_of_nodes)
		exist := false
		for _, v := range nk.output_number {
			if v == output {
				exist = true
				break
			}
		}

		if !exist {
			nk.output_number = append(nk.output_number, output)
		}
	}

	nk.output.mapping_relation = make(map[*Neuron]int64, num)

	for i, v := range nk.output_number {
		// nk.Outputs = append(nk.Outputs, nk.Neurons[i])
		nk.output.mapping_relation[nk.Neurons[v]] = int64(i)
	}

	return
}

func (nk *NeuralNetwork) Init() {
	// ale := ALE.ALE{}
	ale := virtualEnvironment.VirtualENV{}
	nk.env = &ale
	nk.Generate_nodes(10000)

	nk.Fast_generate_random_graph(1000, 0.3, 99)

	num_of_outputs, num_of_inputs := nk.env.Init()

	nk.Generate_inputs(num_of_inputs, 10) // num_of_inputs, seed
	nk.Generate_outputs(num_of_outputs, 100)

}

func (nk *NeuralNetwork) Pick_excited_inputs_to_running_queue() {
	// nk.Running_queue = lane.NewQueue()
	for _, value := range nk.input.mapping_relation {
		if value.Excited == true {
			nk.Running_queue.Enqueue(value)
		}
		// fmt.Println("value:", value.Excited)
	}

}

func (nk *NeuralNetwork) check_outputs() {
	inpt := make([]int64, len(nk.output_number))

	fmt.Println("Before output:", inpt)
	fmt.Println("	outputs length: ", len(nk.output.outputs))

	for _, item := range nk.output.outputs {
		inpt = append(inpt, nk.output.mapping_relation[item])
	}

	fmt.Println("After output:", inpt)

	nk.env.Write_action(inpt)

}

func (nk *NeuralNetwork) put_into_queue(nn *Neuron) {
	nk.Running_queue.Enqueue(nn)
}

func (nk *NeuralNetwork) put_inputs_into_queue(inputs []int64) {
	fmt.Println("inputs list length: ", nk.Running_queue.Size())
	for _, v := range inputs {
		if v > 0 {
			nk.Running_queue.Enqueue(nk.input.mapping_relation[v])
		}
	}
	fmt.Println("inputs list length: ", nk.Running_queue.Size())
}

func (nk *NeuralNetwork) check_inputs() {
	screen_inputs, is_terminated, is_scored := nk.env.Read_state()
	merge_inputs := append(screen_inputs, is_terminated, is_scored)
	nk.put_inputs_into_queue(merge_inputs)

}

func (nk *NeuralNetwork) check_if_outputs(neu *Neuron) {
	if v, ok := nk.output.mapping_relation[neu]; ok {
		nk.output_number = append(nk.output_number, v)
		fmt.Println("outputs nodes added.")
	}
}

func (nk *NeuralNetwork) finish_exciting_transmitting(neu interface{}) {
	if nn, ok := neu.(*Neuron); ok {
		// fmt.Printf("neuron.trans before: ", &nn.trans.p)
		for _, next := range nn.Post_neurons {
			//fmt.Println("next state: ", next.state, next.cell.base_p, next.cell.excit_p, next.cell.pool, next.cell.last_excit_timestamp)
			suc := nn.pass_potential(next)
			// fmt.Println("next state: ", next.state, next.cell.base_p, next.cell.excit_p, next.cell.pool, next.cell.last_excit_timestamp)
			// fmt.Println("neuron.trans after: ", &nn.trans.p)
			if suc == true {
				fmt.Println("transed success")
				nk.put_into_queue(next)
				nk.check_if_outputs(next)
			}
		}
	}

}

func (nk *NeuralNetwork) Boot_up(step int) {
	// putting nil Neuron pointer at each start of step
	// when dequeue a nil pointer, the system will judge inputs and outputs
	start_frame := "sep"
	nk.Running_queue = lane.NewQueue()
	nk.Running_queue.Enqueue(start_frame)
	for step > 0 {
		if nk.Running_queue.Empty() {
			fmt.Println("Dequeue is empty. Unexpectedly exit.")
			break
		}
		neu := nk.Running_queue.Dequeue()
		if neu == start_frame {
			fmt.Println("step: ", step)
			fmt.Println("before list length: ", nk.Running_queue.Size())
			nk.check_outputs()
			fmt.Println("list length: ", nk.Running_queue.Size())
			nk.check_inputs()
			fmt.Println("after inputs list length: ", nk.Running_queue.Size())
			nk.Running_queue.Enqueue(start_frame)

			nk.input_number = make([]int64, 0, 10)
			nk.output_number = make([]int64, 0, 10)

			step--
		} else {
			// fmt.Println("in finish trans")
			nk.finish_exciting_transmitting(neu)
		}
	}
}
