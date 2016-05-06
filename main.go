package main

import(
    "neuron"
    // "fmt"
)

func main() {
    test := neuron.NeuralNetwork{}
    test.Generate_nodes(10)
    // test.Add_edge(0, 1)
    // fmt.Println("neuron pointers", test.Neurons[0].Post_neurons[0].Pre_neurons[0].Emmission_p)
    test.Fast_generate_random_graph(10, 0.3, 99)
}
