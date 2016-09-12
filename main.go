package main

import "neuron"

func main() {

	/*
		logfile, err := os.OpenFile("runtime.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
		if err != nil {
			fmt.Println("%s\r\n", err.Error())
			os.Exit(-1)
		}

		defer logfile.Close()

		logger := log.New(logfile, "\n", log.Ldate|log.Ltime|log.Llongfile)

	*/

	// controller, state := ale.Init()

	test := neuron.NeuralNetwork{}

	step := 100
	test.Init()
	test.Boot_up(step)
	test.Write_to("edges.txt")

}
