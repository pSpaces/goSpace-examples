// Construction of a distance-to-PoI field

package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	. "github.com/pspaces/gospace"
	container "github.com/pspaces/gospace/container"
)

// Size of the "map", coordinates are in [0..MAXX,0..MAXY]
const MAXX = 100
const MAXY = 100

// Number of rows and columns to discretize the map into areas of interest
// Cell = area of interest
const XAREAS = 2
const YAREAS = 2

func main() {

	host, port, ndevices := args()

	// creating a new policy
	policy := NewComposablePolicy()

	// rule "pi" to be used to regulate queries
	policyName := "pi"
	// defining the template of the queries to be controlled
	spc := new(Space)
	var d float64
	var x float64
	var y float64
	var i int
	var j int
	var who int
	// the template of the queries
	template := CreateTemplate(minD, "device", &who, "in", &x, &y, &i, &j, "distanceToPoI", &d)
	templateFields := template.Fields()
	var ltf []interface{}
	ltf = make([]interface{}, len(templateFields)+1)
	copy(ltf[:2], []interface{}{templateFields[0], NewLabels(NewLabel(policyName))})
	copy(ltf[2:], templateFields[1:])
	ltp := CreateTemplate(ltf...)
	// define the action to be controlled
	a := NewAction(spc.QueryAgg, ltp)
	// define the transformations of template, tuple and result
	templateTrans := NewTransformation(TemplateIdentity)
	tupleTrans := NewTransformation(TupleIdentity)
	resultTrans := NewTransformation(TupleIdentity)
	transformations := NewTransformations(&templateTrans, &tupleTrans, &resultTrans)
	// create policy rule with action and transformation
	rule := NewAggregationRule(*a, *transformations)
	// add the aggregation policy with name "pi"
	policy.Add(NewAggregationPolicy(NewLabel(policyName), rule))

	field := NewSpace("tcp://" + host + ":" + port + "/field")
	// with policies
	// field := NewSpace("tcp://" + host + ":" + port + "/field, policy)

	// launche all devices
	for i := 0; i < ndevices; i++ {
		go device(&field, i)
	}

	// wait for all devices to be done
	for i := 0; i < ndevices; i++ {
		field.Get("done")
	}

	// Print the final result as a CSV (for the plot)
	fmt.Println()
	fmt.Println("############")
	fmt.Println("FINAL VALUES")
	fmt.Println("############")
	for i := 0; i < ndevices; i++ {
		tl, _ := field.GetAll("device", i, "in", &x, &y, &i, &j, "distanceToPoI", &d)
		for _, t := range tl {
			d = t.GetFieldAt(8).(float64)
			x = t.GetFieldAt(3).(float64)
			y = t.GetFieldAt(4).(float64)
			fmt.Printf("%f , %f , %f \n", x, y, d)
		}
	}
}

func device(field *Space, me int) {
	var d float64
	var d2 float64
	var x float64
	var y float64
	var x2 float64
	var y2 float64
	var i int
	var j int
	var who int

	// select a random position
	x = MAXX * rand.Float64()
	y = MAXY * rand.Float64()
	// compute the area of interest
	i, j = area(x, y)
	// update the distance to a PoI
	d = distanceToPoI(x, y)

	fmt.Printf("Device %d in (%f,%f) area (%d,%d) distance %f\n", me, x, y, i, j, d)
	field.Put("device", me, "in", x, y, i, j, "distanceToPoI", d)
	// with policy
	// field.Put(NewLabel("pi"), "device", me, "in", x, y, i, j, "distanceToPoI", d)

	// keep aggregating for some rounds
	for rounds := 0; rounds < 10; rounds++ {
		time.Sleep(1 * time.Second)
		// probe the area and adjacent ones + the diagonals
		for ii := i - 1; ii <= i+1; ii++ {
			for jj := j - 1; jj <= j+1; jj++ {
				t, e := field.QueryAgg(minD, "device", &who, "in", &x2, &y2, ii, jj, "distanceToPoI", &d2)
				if e == nil && reflect.TypeOf(t) == reflect.TypeOf(Tuple{}) && t.Length() == 9 {
					d2 = t.GetFieldAt(8).(float64)
					x2 = t.GetFieldAt(3).(float64)
					y2 = t.GetFieldAt(4).(float64)
					d2 = d2 + distanceTo(x, y, x2, y2)
					if d2 < d {
						d = d2
					}
				}
			}
		}
		// update the distance to a PoI
		fmt.Printf("Device %d in (%f,%f) area (%d,%d) distance %f\n", me, x, y, i, j, d)
		field.Get("device", me, "in", x, y, i, j, "distanceToPoI", &d2)
		field.Put("device", me, "in", x, y, i, j, "distanceToPoI", d)
		// with policy
		//field.Put(NewLabel("pi"), "device", me, "in", x, y, i, j, "distanceToPoI", d)
	}

	field.Put("done")

}

// area computes the area of interest (cell) for given coordinates
func area(x float64, y float64) (int, int) {
	return int(math.Min(float64(x/(MAXX/XAREAS)), XAREAS-1)),
		int(math.Min(float64(y/(MAXY/YAREAS)), YAREAS-1))
}

// distancetoPoI computes the distance to the nearest PoO
func distanceToPoI(x float64, y float64) float64 {
	// For now, the only PoI is at (0.0,0.0)
	// A device can see it only if in the same area
	i1, j1 := area(x, y)
	i2, j2 := area(0, 0)
	if i1 == i2 && j1 == j2 {
		return math.Sqrt((x * x) + (y * y))
	} else {
		return math.MaxFloat64
	}
}

// distanceTo computes the Euclidean distance between two points on the map
func distanceTo(x1 float64, y1 float64, x2 float64, y2 float64) float64 {
	return math.Sqrt(((x1 - x2) * (x1 - x2)) + ((y1 - y2) * (y1 - y2)))
}

// minD finds the tuple with the tuple with the smallest distance field
func minD(ts ...Intertuple) container.Intertuple {

	//fmt.Printf("Aggregating %d tuples\n", len(ts))

	if len(ts) == 0 {
		tt := make([]interface{}, 1)
		t := CreateTuple(tt)
		return &t
	}

	var d float64
	d = 257 //math.MaxFloat64
	//tt := make([]interface{}, 1)
	t := ts[0]

	for z := 0; z < len(ts); z++ {
		if ts[z].GetFieldAt(8).(float64) < d {
			d = ts[z].GetFieldAt(8).(float64)
			t = ts[z]
		}
	}

	return t

}

func args() (host string, port string, ndevices int) {

	// default values
	host = "localhost"
	port = "31145"
	ndevices = 10

	flag.Parse()
	argn := flag.NArg()

	if argn > 3 {
		fmt.Printf("Usage: [ndevices] [address] [port] \n")
		return
	}

	if argn >= 2 {
		host = flag.Arg(1)
	}

	if argn >= 3 {
		port = strings.Join([]string{":", flag.Arg(2)}, "")
	}

	if argn >= 1 {
		ndevices, _ = strconv.Atoi(flag.Arg(0))
	}

	return host, port, ndevices
}

// TemplateIdentity is the template identity function.
func TemplateIdentity(i interface{}) (tp container.Template) {
	tpf := i.([]interface{})

	tp = CreateTemplate(tpf...)

	return tp
}

// TupleIdentity is the tuple identity function.
func TupleIdentity(i interface{}) (it container.Intertuple) {
	tf := i.([]interface{})

	t := CreateTuple(tf...)
	it = &t

	return it
}
