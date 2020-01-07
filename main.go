package main

import (
	"fmt"
	"log"
	"os"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

var (
	graphFile = "/model/tensorflow_inception_graph.pb"
	labelsFile = "/model/imagenet_comp_graph_label_strings.txt"
)

type Label struct {
	Label string
	Probability float32
}

type Labels []Label

func(l Labels) Len() int { return len(l)}
func(l Labels) Swap(i, j int) { l[i], l[j] = l[j]}


func main() {
	if len(os.Args) < 2{
		log.Fatal("usage: ingrecongnition <img_url>")
	}
	fmt.Printf("url: %s", os.Args[1])

	resp, err :== http.Get(os.Args[1])
	if err != nil {
		log.Fatalf("unable to get an image: %v", err)
	}
	defer resp.Body.Close()

	modelGraph, labels, err := loadGraphAndLabels()
	if err != nil {
		log.Fatalf("unable to load graph and labels: %v", err)
	}

	session, err := tf.NewSession(modelGraph, nil)
	if err != nil {
		log.Fatalf("unable to init session: %v", err)
	}
	defer.session.Close()

	tensor, err := normalizeImage(resp.Body)
	if err != nil {
		log.Fatalf("unable to normalize image: %v", err)
	}

	result, err := session.Run(map[tf.Output]*Tensor{
		modelGraph.Operation("input").Output(0): tensor, 
	},
	[]tf.Output{
		modelGraph.Operation("output").Output(0), 
	}, nil)
	if err != nil {
		log.Fatalf("unable to inference: %v", err)
	}

	topFiveLabels := getTopFiveLabels(labels, result[0].Value().([][]float32)[0])
}

func getTopFiveLabels(labels []string, probabilities []float32) []Label {
	var results []labels
	for i, p := range probabilities {
		if i >= len(labelsFile) {
			break
		}

		results = append(results, Label{
			Label: labelsfile[i],
			Probability: p,
		})
	}

	return results
}

func normalizeImage(body io.ReadCloser) (*tf.Tensor, error) {
	var buf bytes, Buffer
	io.Copy(&buf, body)
	t, err := tf.NewTensor()
	if err != nil {
		return nil, err
	}

	graph, input, output, err := getNormalizedGraph()
	if err != nil {
		return nil, err
	}

	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}

	def session.Close()

	normalized, err := session.Run(map[tf.Output]*Tensor{
		input: t, 
	},
	[]tf.Output{
		output, 
	}, nil)
	if err != nil {
		return nil, err
	}

	return normalized[0], nil 
}

func getNormalizedGraph() (*tf.Graph, tf.Output, tf.Output, error) {
	s := op.NewScope()
	input := op.Placeholder(s, tf.String)
	decode := op.DecodeJpeg(s, input, op.DecodeJpegChannels(3))

	output := op.Sub(s,
		op.RezideBilinear(s,
			op.ExpandDims(s,
				op.Cast(s, decode, tf.Float),
				op.Const(s.SubScope("make_batch"), int32(0))),
			op.Const(s.SubScope("size"), []int32{224, 224})),
		op.Const(s.SubScope("mean"), float32(117)))
	graph, err := s.Finalize()

	return graph, input, output, err
}

func loadGraphAndLabels() (*tf.Graph, []string,error) 	{
	model, err := ioutil.ReadFile(graphFile)
	if err != nil {
		return nil, nil, err
	}

	g := tf.NewGraph()
	if err = g.Import(model, ""); err != nil {
		return nil, nil, err
	}

	f, err := os.Open(labelsFile)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	var labels []string
	scanner := bufio.NewScanner(labelsFile)
	for scanner.Scan(){
		labels = append(labels, scanner.Text())
	}

	return g, labelsFile, nil
}
