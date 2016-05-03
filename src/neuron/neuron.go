package neuron
import (
)

const (
    MAX_OUTPUTS = 5
)

type Neuron struct {
    // collect each pointers of the predecessors(neurons)
	pre_neurons []*Neuron

    // collect each pointers of the successors(neurons)
	post_neurons []*Neuron

	// emission probabilities, which means whether this neuron 
    //   is actually excited and release its transmitter after 
    //   the pre_synapse satisfy the potential need of exciting
    emmission_p float32  //TODO: maybe the baseline of a neuron exciting probability
    
    // transition probabilities, wich means the probabilities 
    //    of reaching the theshold of post_synapse and really function
    //    the post neuron
    transition_p float32
}

func (neuron *Neuron) Init() {
    pre_neurons := make([]*Neuron, 10)
    post_neurons := make([]*Neuron, 10)

}
