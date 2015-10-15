package cell
import {
	"fmt"
}

const{
	MAX_INPUTS := 20
	MAX_OUTPUTS := 5
    ENERGY := -70
}

type Dentrite struct {
	var preAxon = map[int]Axon
	weights int
    time_delay int
    excited_count int
    last_excited_time timestamp
}

type Axon stuct {
    excited bool
    //time_squence []bool
}

type Neuron struct {
    ENERGY 
	excitedTheshold int;
	intputs [MAX_INPUTS]Dentrite
	var outputs [MAX_OUTPUTS]Axon
}

func createInput(neuron *Neuron) {
}

func cancelInput(neuron *Neuron) {

}

func generateOutput(neuron *Neuron) {

}

func isExcited(neuron *Neuron) {

}

func doExcited() {

}
