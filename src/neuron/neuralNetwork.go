package neuron

import (
	"fmt"
	//"math"
	//"math/rand"
	"os"
)

type NeuralNetwork struct {
	inputs  map[*Neuron]int64
	outputs map[*Neuron]int64

	Neurons []*Neuron
}

func (nk *NeuralNetwork) Generate_nodes(num int) {
	for i := 0; i < num; i++ {\
	}
}