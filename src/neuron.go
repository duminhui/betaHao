package cell
import {
	"fmt"
}

const{
	MAX_INPUTS := 20
	MAX_OUTPUTS := 5
}

type Dentrite struct {
	var preNeurons = map[int]Neuron
	weights int
}

type BranchOfAxon struct {
	//Anox can form Axon Collaterals(branches). Equally, we define these as multipal outputs.
	var postNeurons = map[int]Neuron //TODO:write add function to roll query
	Property bool //need to consider more, perhaps we can put it into dentrites weights
}

type Neuron struct {
	excitedTheshold int;
	intputs [MAX_INPUTS]Dentrite
	var outputs [MAX_OUTPUTS]BranchOfAxon
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
