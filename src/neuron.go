package main
import (
)

const (
    MAX_OUTPUTS = 5
)

type Axon struct {
	map[int]Neuron
    delta_p int
    // time_delay int
    // excited_count int
}

type Neuron struct {
	excited_p int
	outputs [MAX_OUTPUTS]Axon
}
