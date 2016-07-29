package neuron

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	// "log"
	// "time"
	// "ALE"
	"github.com/oleiade/lane"
	"math"
	"virtualEnvironment"
	// "sync"
)

var step int64

type Input struct {
	// inputs           []*Neuron // 运行态
	mapping_relation map[int64]*Neuron
}

type Output struct {
	outputs          []*Neuron // 运行态
	mapping_relation map[*Neuron]int64
}

type NeuralNetwork struct {
	Neurons       []*Neuron
	running_queue *lane.Queue
	env           Environmenter

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

	pre.post_neurons = append(pre.post_neurons, post)
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
		// fmt.Printf("inpurt_order: %v, %T\n", i, i)
		// nk.Inputs[int64(i)] = nk.Neurons[v]
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
		// nk.Outputs = append(nk.Outputs, nk.Neurons[i])
		nk.output.mapping_relation[nk.Neurons[v]] = int64(i)
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

}

func (nk *NeuralNetwork) Pick_excited_inputs_to_running_queue() {
	// nk.Running_queue = lane.NewQueue()
	for _, value := range nk.input.mapping_relation {
		if value.Excited == true {
			nk.running_queue.Enqueue(value)
		}
		// fmt.Println("value:", value.Excited)
	}

}

func (nk *NeuralNetwork) check_outputs() {
	inpt := make([]int64, 0)

	// fmt.Println("Before output:", inpt)
	// fmt.Println("	outputs length: ", len(nk.output.outputs))

	for _, item := range nk.output.outputs {
		inpt = append(inpt, nk.output.mapping_relation[item])
	}

	fmt.Println(" output list:", inpt)

	nk.env.Write_action(inpt)

}

func (nk *NeuralNetwork) Write_to(filename string) {
	/*
		for _, v := range nk.Neurons {
			v.reversal_tag = false

		}

		begin := nk.Neurons[0]
	*/

	reslt, _ := json.Marshal(nk)

	out, _ := os.Create(filename)
	defer out.Close()

	out.Write(reslt)
}

func (nk *NeuralNetwork) put_into_queue(nn *Neuron) {
	nk.running_queue.Enqueue(nn)
}

func (nk *NeuralNetwork) put_inputs_into_queue(inputs []int64) {
	// fmt.Println("inputs list length: ", nk.running_queue.Size())
    fmt.Println(" input mapping_relation: ", len(nk.input.mapping_relation))
	for i, v := range inputs {
		if v > 0 {
			nk.running_queue.Enqueue(nk.input.mapping_relation[int64(i)])
		}
	}
	// fmt.Println("inputs list length: ", nk.running_queue.Size())
}

func (nk *NeuralNetwork) check_inputs() {
	screen_inputs, is_terminated, is_scored := nk.env.Read_state()
	merge_inputs := append(screen_inputs, is_terminated, is_scored)
    fmt.Println(" input list: ", merge_inputs)
	nk.put_inputs_into_queue(merge_inputs)

}

func (nk *NeuralNetwork) check_if_outputs(neu *Neuron) {
	if _, ok := nk.output.mapping_relation[neu]; ok {
		nk.output.outputs = append(nk.output.outputs, neu)
		fmt.Println("outputs nodes added.")
	}
}

func (nk *NeuralNetwork) finish_exciting_transmitting(neu interface{}) {
	if nn, ok := neu.(*Neuron); ok {
		// fmt.Printf("neuron.trans before: ", &nn.trans.p)
		for _, next := range nn.post_neurons {
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
	nk.running_queue = lane.NewQueue()
	nk.running_queue.Enqueue(start_frame)
	for step > 0 {
		if nk.running_queue.Empty() {
			fmt.Println("Dequeue is empty. Unexpectedly exit.")
			break
		}
		neu := nk.running_queue.Dequeue()
		if neu == start_frame {
			fmt.Println("In step: ", step)
			// fmt.Println("before list length: ", nk.running_queue.Size())
			nk.check_outputs()
			// fmt.Println("list length: ", nk.running_queue.Size())
			nk.check_inputs()
			// fmt.Println("after inputs list length: ", nk.running_queue.Size())
			nk.running_queue.Enqueue(start_frame)

			// nk.input.inputs = make([]int64, 0, 10)
			nk.output.outputs = make([]*Neuron, 0)

			step--
		} else {
			// fmt.Println("in finish trans")
			nk.finish_exciting_transmitting(neu)
		}
	}
}
