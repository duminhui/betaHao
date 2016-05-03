package neuron

import (
    "fmt"
)

type NeuralNetwork struct {
	neurons      []*Neuron
    
}

func (neuralnetwork *NeuralNetwork) generate_nodes(num int) {
    // initialize 'num' numbers of neurons in the network
    neurons := make([]*Neuron, num)
    for i := 0; i <= num; i++ {
        p := &Neuron{emmission_p:1, transition_p:0}
        neurons = append(neurons, p)
        fmt.Printf("aaa %s \n", *p.pre_neurons[i])
    }
    fmt.Printf("bbb %s\n", neurons[0])
    fmt.Printf("bbb %s\n", neurons[1])

    // return &neurons
}
