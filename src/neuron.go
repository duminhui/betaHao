package main
import (
)

const (
    MAX_OUTPUTS = 5
)

type Axon struct {
	next map[*int]Neuron
    delta_p float32
    // time_delay int
    // excited_count int
}

type Neuron struct {
	excited_p float32
	outputs [MAX_OUTPUTS]Axon
}
