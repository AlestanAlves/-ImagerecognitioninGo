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
		log.Fatalf("unable to load graph and labels: %w", err)
	}
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
