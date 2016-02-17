package main

import (
    "fmt"
    "./src"
)

/*type NeuralNetwork struct {
	HiddenLayer      [][]Neuron
	InputLayer       [][]Neuron
	OutputLayer      [][]Neuron
}*/

func main(){

    var neu Neuron
    neu.excited_p = 0.1

    fmt.Printf("The per is %s", neu.excited_p)
    
}
