package main

import(
    "neuron"
    "fmt"
    "ALE"
)

func main() {

    // ale := ALE{}
    // _,_ := ale.Init()

    test := neuron.NeuralNetwork{}
    test.Generate_nodes(100)
    // test.Add_edge(0, 1)
    // fmt.Println("neuron pointers", test.Neurons[0].Post_neurons[0].Pre_neurons[0].Emmission_p)
    test.Fast_generate_random_graph(10000, 0.3, 99)

    test.Init(ALE)
    // test.Generate_inputs(5, 10)
    // test.Generate_outputs(5, 10)

    test.Inputs[0].Excited = true
    // test.Pick_excited_inputs_to_running_queue()

    // temp := test.Running_queue.Dequeue()

    fmt.Println("inputs:", test.Inputs)
    fmt.Println("outputs:", test.Outputs)
    // fmt.Println("temp:", temp)

}
