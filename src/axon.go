package cell

const MAX_TIME_SPAN uint = 5

type Axon struct {
	pre_excited  bool
	time_squence []bool
}

type AxonList struct {
	axon       Axon
	time_delay uint //every axon has its own delay when the generating time squence deploying on it.
}
