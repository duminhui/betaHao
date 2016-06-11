package neuron

import (
	"fmt"
	"math/rand"
	// "time"
	"github.com/oleiade/lane"
	"math"
	// "sync"
)

type NeuralNetwork struct {
	Neurons       []*Neuron
	Running_queue *lane.Queue
	start_frame   string

	Inputs            map[int64]*Neuron
	Outputs           map[*Neuron]int64
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
		p := &Neuron{Emmission_p: 1, Transition_p: 0}
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
			// fmt.Println("v: %v, w: %v ", v, w)
			// fmt.Println("W:", w)
			// _ = "breakpoint"
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
	fmt.Printf("len of input_order: %v", len(input_order))

	// inputs = make([] *Neuron, 0)
	/*
	   for i, _ := range input_order {
	       // nk.Inputs = append(nk.Inputs, nk.Neurons[v])
	       fmt.Printf("inpurt_order: %v, %T\n", i, i)
	       // nk.Inputs[int64(i)] = nk.Neurons[v]
	   }
	*/

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

	nk.Outputs = make(map[*Neuron]int64, num)

	for i := 0; i < len(output_order); i++ {
		// nk.Outputs = append(nk.Outputs, nk.Neurons[i])
		nk.Outputs[nk.Neurons[i]] = int64(i)
	}

	return
}

func (nk *NeuralNetwork) Init(env Environmenter) {
	// instance := NeuralNetwork{}
	nk.Generate_nodes(1000)

	nk.Fast_generate_random_graph(1000, 0.3, 99)

	num_of_outputs, num_of_inputs := env.Init()

	nk.Generate_inputs(num_of_inputs, 10) // num_of_inputs, seed
	nk.Generate_outputs(num_of_outputs, 100)

}

func (nk *NeuralNetwork) Pick_excited_inputs_to_running_queue() {
	nk.Running_queue = lane.NewQueue()
	for _, value := range nk.Inputs {
		if value.Excited == true {
			nk.Running_queue.Enqueue(value)
		}
		fmt.Println("value:", value.Excited)
	}

}

func (nk *NeuralNetwork) check_outputs() {
	excited_outputs := make([]*Neuron, 0)
	for key, _ := range nk.Outputs {
		if key.Excited == true {
			excited_outputs = append(excited_outputs, key)
		}
		// fmt.Println("value:", value.Excited)
	}

}

func (nk *NeuralNetwork) explain_outputs() {

}

func (nk *NeuralNetwork) check_inputs() {

}

func (nk *NeuralNetwork) finish_exciting_transmitting() {

}

func (nk *NeuralNetwork) Boot_up(step int) {
	// putting nil Neuron pointer at each start of step
	// when dequeue a nil pointer, the system will judge inputs and outputs
	nk.start_frame = "start"
	nk.Running_queue.Enqueue(nk.start_frame)
	for ; step > 0; step-- {
		neu := nk.Running_queue.Dequeue()
		if neu == nk.start_frame {
			nk.Running_queue.Enqueue(nk.start_frame)
			nk.check_outputs()
			nk.explain_outputs()
			nk.check_inputs()
		} else {
			nk.finish_exciting_transmitting()
		}
	}
}

func (nk *NeuralNetwork) Read_state(env Environmenter) {
	screen_inputs, is_terminated, is_scored := env.Read_state()
	_ = "breakpoint"
	fmt.Println(screen_inputs)
	fmt.Println(is_terminated)
	fmt.Println(is_scored)
}
