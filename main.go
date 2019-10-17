package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func loadData() ([][]float64, []string, error) {
	f, err := os.Open("iris.csv")
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ','
	r.LazyQuotes = true
	_, err = r.Read()
	if err != nil {
		return nil, nil, err
	}
	rows, err := r.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	X := [][]float64{}
	Y := []string{}
	for _, cols := range rows {
		x := make([]float64, 4)
		y := cols[4]
		for j, s := range cols[:4] {
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, nil, err
			}
			x[j] = v
		}
		X = append(X, x)
		Y = append(Y, y)
	}
	return X, Y, nil
}

func distance(lhs, rhs []float64) float64 {
	val := 0.0
	for i, _ := range lhs {
		val += math.Pow(lhs[i]-rhs[i], 2)
	}
	return math.Sqrt(val)
}

func rotate(data [][]float64) [][]float64 {
	result := make([][]float64, len(data[0]))
	for i := 0; i < len(data[0]); i++ {
		result[i] = make([]float64, len(data))
	}
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[0]); j++ {
			result[j][i] = data[i][j]
		}
	}
	return result
}

func minIdx(arr []float64) int {
	minv := arr[0]
	mini := 0
	for i, v := range arr[1:] {
		if v < minv {
			minv = v
			mini = i + 1
		}
	}
	return mini
}

type krange struct {
	min float64
	max float64
}

func minMax(XX [][]float64) []krange {
	result := []krange{}
	for _, arr := range rotate(XX) {
		r := krange{
			min: arr[0],
			max: arr[0],
		}
		for _, v := range arr[1:] {
			if r.min > v {
				r.min = v
			}
			if r.max < v {
				r.max = v
			}
		}
		result = append(result, r)
	}
	return result
}

func sameAll(a, b []int) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func fit(XX [][]float64, k int) []int {
	box := minMax(XX)
	representatives := make([][]float64, k)

	for i := 0; i < k; i++ {
		representatives[i] = make([]float64, len(box))
		for j := 0; j < len(box); j++ {
			vec := box[j].max - box[j].min
			off := box[j].min
			representatives[i][j] = rand.Float64()*vec + off
		}
	}

	idx := []int{}
	for _, arr := range XX {
		var vec []float64
		for _, r := range representatives {
			vec = append(vec, distance(arr, r))
		}
		idx = append(idx, minIdx(vec))
	}

	for {
		newRepresentatives := [][]float64{}
		for i, _ := range representatives {
			var group [][]float64
			for j, x := range XX {
				if idx[j] == i {
					group = append(group, x)
				}
			}
			if len(group) == 0 {
				continue
			}

			smallvec := []float64{}
			for _, arr := range rotate(group) {
				sum := 0.0
				for _, v := range arr {
					sum += v
				}
				smallvec = append(smallvec, sum/float64(len(arr)))
			}
			newRepresentatives = append(newRepresentatives, smallvec)
		}
		representatives = newRepresentatives

		newLabels := []int{}
		for _, d := range XX {
			var dvec []float64
			for _, r := range representatives {
				dvec = append(dvec, distance(d, r))
			}
			newLabels = append(newLabels, minIdx(dvec))
		}
		if sameAll(idx, newLabels) {
			break
		}
		idx = newLabels
	}
	return idx
}

func main() {
	X, Y, err := loadData()
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())

	indexes := fit(X, 3)

	labels := []string{
		"Iris-setosa",
		"Iris-versicolor",
		"Iris-virginica",
	}
	m := map[int]string{}
	for _, l := range labels {
		mm := map[int]int{}
		for j, y := range Y {
			if y != l {
				continue
			}
			mm[indexes[j]]++
		}
		maxv := 0
		maxi := 0
		for k, v := range mm {
			if maxv < v {
				maxv = v
				maxi = k
			}
		}
		m[maxi] = l
	}
	fmt.Println(m)
	crustered := make([]string, len(Y))
	for i, v := range indexes {
		crustered[i] = m[v]
	}
	correct := 0
	for i, _ := range crustered {
		if crustered[i] == Y[i] {
			correct += 1
		}
	}

	fmt.Printf("%f%%\n", float64(correct)/float64(len(crustered))*100)
}
