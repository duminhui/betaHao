package main

import(
    "neuron"
    "fmt"
)

func main() {
    test := neuron.NeuralNetwork{}
    test.Generate_nodes(100)
    // test.Add_edge(0, 1)
    // fmt.Println("neuron pointers", test.Neurons[0].Post_neurons[0].Pre_neurons[0].Emmission_p)
    test.Fast_generate_random_graph(100, 0.3, 99)
    inputs := test.Generate_inputs(5, 10)
    outputs := test.Generate_outputs(5, 10)

    fmt.Println("inputs:", inputs)
    fmt.Println("outputs:", outputs)

}
