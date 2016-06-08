package main

import (
	"ALE"
	"fmt"
	"neuron"
)

func main() {

	ale := ALE.ALE{}
	// controller, state := ale.Init()

	test := neuron.NeuralNetwork{}
	// test.Generate_nodes(1000)
	// test.Add_edge(0, 1)
	// fmt.Println("neuron pointers", test.Neurons[0].Post_neurons[0].Pre_neurons[0].Emmission_p)
	// test.Fast_generate_random_graph(1000, 0.3, 99)

	// _ = "breakpoint"

	test.Init(&ale)
	test.Read_state(&ale)
	// test.Generate_inputs(5, 10)
	// test.Generate_outputs(5, 10)

	// test.Inputs[0].Excited = true
	// test.Pick_excited_inputs_to_running_queue()

	// temp := test.Running_queue.Dequeue()

	fmt.Println("inputs:", test.Inputs)
	fmt.Println("outputs:", test.Outputs)
	// fmt.Println("temp:", temp)

}
