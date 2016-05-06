package neuron

import (
    "fmt"
    "math/rand"
    // "time"
    "math"
)

type NeuralNetwork struct {
	Neurons      []*Neuron
    
}

func (nk *NeuralNetwork) Generate_nodes(num int) {
    // initialize 'num' numbers of neurons in the network
    // nk.neurons = make([]*Neuron, num)
    for i := 0; i < num; i++ {
        p := &Neuron{Emmission_p:1, Transition_p:0}
        p.Init()
        nk.Neurons = append(nk.Neurons, p)
        fmt.Println(nk.Neurons)
    }
}

func (nk *NeuralNetwork) Add_edge(pre_neuron int, post_neuron int) {
    pre := nk.Neurons[pre_neuron]
    post := nk.Neurons[post_neuron]

    pre.Post_neurons = append(pre.Post_neurons, post)
    post.Pre_neurons = append(post.Pre_neurons, pre)

}

func (nk *NeuralNetwork) Remove_edge(pre_neuron int, post_neuron int) {
    // v := 0
    // w := 0
    // TODO: is this necessary
}

func (nk *NeuralNetwork) Fast_generate_random_graph(n int, p float64, seed int64) {
    r := rand.New(rand.NewSource(seed))
    // r := rand.New(rand.NewSource(time.Now().UnixNano())

    v := 0
    w := -1
    lp := math.Log(1.0 - p)

    for v < n {
        lr := math.Log(1.0 - r.Float64())
        w = w + 1 + int(lr/lp)
        if v == w {
            w = w + 1
        }
        for v < n && n <= w {
            w = w-n
            v = w + 1
            if v == w {
                w = w + 1
            }
        }
        if v < n {
            nk.Add_edge(v, w)
        }
    }

}

func (nk *NeuralNetwork) generate_inputs(num int, seed int64) (inputs []*Neuron){
    r := rand.New(rand.NewSource(seed))
    input_order := make([]int, 0)
    num_of_nodes := len(nk.Neurons)
    for num_of_nodes < num {
        input := r.Intn(num_of_nodes)
        exist := false
        for _, v := range input_order {
            if v==input {
                exist = true
                break
            }
        }
        
        if !exist {
            input_order = append(input_order, input)
        }
    }

    inputs = make([] *Neuron, 0)

    for _, v := range input_order {
        inputs = append(inputs, nk.Neurons[v])
    }

    return
}

func (nk *NeuralNetwork) generate_outputs(num int, seed int64) (outputs []*Neuron){
    r := rand.New(rand.NewSource(seed))
    output_order := make([]int, 0)
    num_of_nodes := len(nk.Neurons)
    for num_of_nodes < num {
        output := r.Intn(num_of_nodes)
        exist := false
        for _, v := range output_order {
            if v==output {
                exist = true
                break
            }
        }
        
        if !exist {
            output_order = append(output_order, output)
        }
    }

    outputs = make([] *Neuron, 0)

    for _, v := range output_order {
        outputs = append(outputs, nk.Neurons[v])
    }

    return
}
