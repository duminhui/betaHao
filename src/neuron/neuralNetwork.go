package neuron

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"

	"github.com/eapache/queue"
)

const BRANCH_OF_EACH_NEURON int = 3
const NUMBER_OF_NEURONS int = 10
const NUMBER_OF_INPUTS int = 3
const NUMBER_OF_OUTPUTS int = 2

type NeuralNetwork struct {
	inputs  []*Neuron
	outputs []*Neuron

	Neurons [NUMBER_OF_NEURONS]*Neuron
}

func (nk *NeuralNetwork) Generate_nodes() {
	for i := 0; i < NUMBER_OF_NEURONS; i++ {
		p := &Neuron{}
		p.index = i

		nk.Neurons[i] = p
		// fmt.Println(nk.Neurons[i])
	}
}

func (nk *NeuralNetwork) Generate_random_graph() {

	for i := 0; i < NUMBER_OF_NEURONS; i++ {
		for j := 0; j < BRANCH_OF_EACH_NEURON; j++ {

			for {
				next_i := rand.Intn(NUMBER_OF_NEURONS)
				if next_i != i {
					// fmt.Println(&nk.Neurons[next_i])
					// 添加一条边
					nk.Neurons[i].branches[j].next = nk.Neurons[next_i]
					// fmt.Println(nk.Neurons[i])
					break
				}
			}

		}
	}

}

func (nk *NeuralNetwork) Generate_inputs() {
	nk.inputs = nk.Neurons[0:NUMBER_OF_INPUTS]
}

func (nk *NeuralNetwork) Generate_outputs() {
	nk.outputs = nk.Neurons[NUMBER_OF_NEURONS-NUMBER_OF_OUTPUTS : NUMBER_OF_NEURONS]
}

func (nk *NeuralNetwork) Read_outputs(learn_mode bool, expected_out []bool) (result []bool) {
	result = make([]bool, len(expected_out))

	if learn_mode == true { // 根据期望输出与实际输出作对比，完成输出神经元的学习

		for i := 0; i < len(nk.outputs); i++ {
			if nk.outputs[i].register.state == Excited && expected_out[i] == false {
				result[i] = false
				nk.outputs[i].register.branch.Decrease()
			} else {
				if nk.outputs[i].register.state == Excited {
					result[i] = true
				} else {
					result[i] = false
				}
			}
		}

	} else {

		for i := 0; i < len(nk.outputs); i++ {
			if nk.outputs[i].register.state == Excited {
				result[i] = true
			} else {
				result[i] = false
			}
		}

	}
	return
}

func (nk *NeuralNetwork) Write_inputs(RGB []byte, running_queue *queue.Queue) {
	var idx int64 = 0
	var offset int
	var bit int
	bin_buf := bytes.NewBuffer(RGB)
	for i := 0; i < len(RGB); i++ {
		offset = 1
		for j := 0; j < 8; j++ {
			var x int
			binary.Read(bin_buf, binary.BigEndian, &x)
			bit = x & offset
			if bit == 1 {
				running_queue.Add(nk.inputs[idx])
			}
			offset <<= 1
			idx++
		}
	}
}

func Get_cifar_data(ff *os.File) (expected_out []byte, vision_data []byte) {

	expected_out = make([]byte, 1)
	ff.Read(expected_out)
	// fmt.Printf("%d bytes: %d\n", n1, b1)
	vision_data = make([]byte, 3072)
	ff.Read(vision_data)
	return expected_out, vision_data
}

func Conversion_from_byte_array_to_bool_array(byte_array []byte) (bool_array []bool) {
	bool_array = make([]bool, len(byte_array)*8)

	var idx int64 = 0
	var offset int
	var bit int

	// fmt.Println(len(byte_array))
	// fmt.Println(byte_array[0])
	for i := 0; i < len(byte_array); i++ {
		offset = 1
		for j := 0; j < 8; j++ {
			var x int
			x = int(byte_array[i])

			bit = x & offset //101&100=100;010&100=0
			// fmt.Println(x, offset, bit)
			if bit > 0 {
				bool_array[idx] = true
			} else if bit == 0 {
				bool_array[idx] = false
			}

			offset <<= 1
			idx++
		}
	}
	return bool_array
}
func (nk *NeuralNetwork) Boot_up(step int) {
	// ff, _ := os.Open("./cifar-10-batches-bin/data_batch_3.bin")

	running_queue := queue.New()
	var neuron_deduplicator map[*Neuron]struct{}

	null_neuron := &Neuron{}
	running_queue.Add(null_neuron)

	var nn *Neuron
	// var expected_out []byte
	var vision_data []byte
	for i := 0; i < step; i++ {
		nn = running_queue.Peek().(*Neuron)
		running_queue.Remove()

		if nn == null_neuron {

			neuron_deduplicator = make(map[*Neuron]struct{})

			// expected_out, vision_data = Get_cifar_data(ff)

			// 暂时去掉，为了测试bool数组的读取情况
			// action := nk.Read_outputs(false, expected_out)
			nk.Write_inputs(vision_data, running_queue)

			running_queue.Add(nn)
			global_step = global_step + 1

		} else if nn.register.state == Excited {
			for i := 0; i < len(nn.branches); i++ {

				result := nn.branches[i].touch()
				if result {
					if _, ok := neuron_deduplicator[nn]; !ok {
						running_queue.Add(nn) //?
					}
				}

			}

		}
	}
	// ff.Close()
}

func Test() {

	// 新建一个神经网络NeuralNetwork，包含inputs、outputs、Neurons的指针组成的数组
	// inputs、outputs、Neuron三者的数据类型皆为：[]*Neuron
	// 由于Neuron结构体内存在索引（index int）变量，生成Neuron节点时需给每一个Neuron的index赋值
	// inputs和outputs可以通过自带的函数直接生成，而不需要重新初始化和命名新的空间

	// 初始化NeuralNetwork之后，需要用指针表示Neuron之间的关系
	var nk NeuralNetwork
	nk.Generate_nodes()
	nk.Generate_inputs()
	nk.Generate_outputs()
	// fmt.Println(nk.Neurons[0].branches[0].next)
	// nk.Neurons[0].branches[0].next = nk.Neurons[1]
	nk.Generate_random_graph()
	fmt.Println(nk.Neurons[0], nk.Neurons[1])
	// fmt.Println(nk)
	ff, _ := os.Open("./cifar-10-batches-bin/data_batch_3.bin")
	expected_out, vision_data := Get_cifar_data(ff)
	// fmt.Println(expected_out)
	// fmt.Println(vision_data)
	eo := Conversion_from_byte_array_to_bool_array(expected_out)
	vd := Conversion_from_byte_array_to_bool_array(vision_data)
	fmt.Println(eo)
	fmt.Println(vd)

	expected_out, vision_data = Get_cifar_data(ff)
	eo = Conversion_from_byte_array_to_bool_array(expected_out)
	vd = Conversion_from_byte_array_to_bool_array(vision_data)
	fmt.Println(eo)
	fmt.Println(vd)
}
