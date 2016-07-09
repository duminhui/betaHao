package neuron

import (
	"fmt"
	"math/rand"
	// "time"
	"ALE"
	"github.com/oleiade/lane"
	"math"
	// "sync"
)

var step int64 = 0

type Input struct {
	inputs           []int64 // 运行态
	mapping_relation map[int64]*Neuron
}

type Output struct {
	game_operatror   []*Neuron // 运行态
	mask_influences  []*Neuron // 运行态
	mapping_relation map[*Neuron]int64
}

type NeuralNetwork struct {
	Neurons       []*Neuron
	Running_queue *lane.Queue
	env           Environmenter

	input             Input
	output            Output
	num_of_controller int64
	num_of_state      int64
}

type Environmenter interface {
	Init() (num_of_controller int64, num_of_state int64)
	Read_state() (screen_list []int64, is_terminated int64, is_scored int64)
	Write_action(excited_actions []int64, excited_sensors_controller []int64)
	Final()
}

func (nk *NeuralNetwork) Generate_nodes(num int) {
	// initialize 'num' numbers of neurons in the network
	// nk.neurons = make([]*Neuron, num)
	for i := 0; i < num; i++ {
		p := &Neuron{}
		p.cell.base_p = 1
		// :p.Init()
		nk.Neurons = append(nk.Neurons, p)
		// fmt.Println(nk.Neurons)
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
		}
	}
	fmt.Println("ER graph generated")

}

func (nk *NeuralNetwork) Generate_inputs(num int64, seed int64) {
	r := rand.New(rand.NewSource(seed))
	input_order := make([]int64, num)
	num_of_nodes := int64(len(nk.Neurons))
	for int64(len(input_order)) < num {
		input := r.Int63n(num_of_nodes)
		exist := false
		for _, v := range input_order {
			if v == input {
				exist = true
				break
			}
		}

		if !exist {
			input_order = append(input_order, input)
		}
	}
	fmt.Printf("len of input_order: %v \n", len(input_order))
	// fmt.Printf("len of inputs: %v \n", num)

	nk.input.mapping_relation = make(map[int64]*Neuron, num)
	for i, v := range input_order {
		nk.input.mapping_relation[int64(i)] = nk.Neurons[v]
		// fmt.Printf("inpurt_order: %v, %T\n", i, i)
		// nk.Inputs[int64(i)] = nk.Neurons[v]
	}

	_ = "breakpoint"

	return
}

func (nk *NeuralNetwork) Generate_outputs(num int64, seed int64) {
	r := rand.New(rand.NewSource(seed))
	output_order := make([]int64, num)
	num_of_nodes := int64(len(nk.Neurons))
	for int64(len(output_order)) < num {
		output := r.Int63n(num_of_nodes)
		exist := false
		for _, v := range output_order {
			if v == output {
				exist = true
				break
			}
		}

		if !exist {
			output_order = append(output_order, output)
		}
	}

	nk.output.mapping_relation = make(map[*Neuron]int64, num)

	for i, v := range output_order {
		// nk.Outputs = append(nk.Outputs, nk.Neurons[i])
		nk.output.mapping_relation[nk.Neurons[v]] = int64(i)
	}

	return
}

func (nk *NeuralNetwork) Init() {
	ale := ALE.ALE{}
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
	oper := make([]int64, len(nk.output.game_operatror))
	mask := make([]int64, len(nk.output.mask_influences))

	for _, item := range nk.output.game_operatror {
		oper = append(oper, nk.output.mapping_relation[item])
	}

	for _, item := range nk.output.mask_influences {
		oper = append(oper, nk.output.mapping_relation[item])
	}

	nk.env.Write_action(oper, mask)

}

func (nk *NeuralNetwork) put_into_queue(nn *Neuron) {
	nk.Running_queue.Enqueue(nn)
}

func (nk *NeuralNetwork) put_inputs_into_queue(inputs []int64) {
	for _, v := range inputs {
		if v > 0 {
			nk.Running_queue.Enqueue(nk.input.mapping_relation[v])
		}
	}
}

func (nk *NeuralNetwork) check_inputs() {
	screen_inputs, is_terminated, is_scored := nk.env.Read_state()
	merge_inputs := append(screen_inputs, is_terminated, is_scored)
	nk.put_inputs_into_queue(merge_inputs)

}

func (nk *NeuralNetwork) finish_exciting_transmitting(neu interface{}) {
	if nn, ok := neu.(Neuron); ok {
		for _, next := range nn.Post_neurons {
			fmt.Println("neuron.trans before: ", neuron.trans)
			suc := nn.pass_potential(next)
			fmt.Println("neuron.trans after: ", neuron.trans)
			if suc == true {
				nk.put_into_queue(next)
			}
		}
	}

}

func (nk *NeuralNetwork) Boot_up(step int) {
	// putting nil Neuron pointer at each start of step
	// when dequeue a nil pointer, the system will judge inputs and outputs
	start_frame := "sep"
	_ = "breakpoint"
	nk.Running_queue = lane.NewQueue()
	nk.Running_queue.Enqueue(start_frame)
	for step > 0 {
		_ = "breakpoint"
		if nk.Running_queue.Empty() {
			fmt.Println("Dequeue is empty. Unexpectedly exit.")
			break
		}
		neu := nk.Running_queue.Dequeue()
		if neu == start_frame {
			nk.check_outputs()
			nk.check_inputs()
			nk.Running_queue.Enqueue(start_frame)
			step--
		} else {
			nk.finish_exciting_transmitting(neu)
		}
	}
}
