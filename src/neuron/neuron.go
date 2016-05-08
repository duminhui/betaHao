package neuron
import (
    // "fmt"
)

const (
    MAX_OUTPUTS = 5
)

type Neuron struct {
    // collect each pointers of the predecessors(neurons)
	Pre_neurons []*Neuron

    // collect each pointers of the successors(neurons)
	Post_neurons []*Neuron

    Excited bool

	// emission probabilities, which means whether this neuron 
    //   is actually excited and release its transmitter after 
    //   the pre_synapse satisfy the potential need of exciting
    Emmission_p float32  //TODO: maybe the baseline of a neuron exciting probability
    
    // transition probabilities, wich means the probabilities 
    //    of reaching the theshold of post_synapse and really function
    //    the post neuron
    Transition_p float32
}

func (nn *Neuron) Init() {
    // nn.pre_neurons = make([]*Neuron, 10)
    // nn.post_neurons = make([]*Neuron, 10)
}

func (nn *Neuron) GetEmmissionP() float32 {
    return nn.Emmission_p
}
