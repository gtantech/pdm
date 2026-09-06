// example taken from [Engineer4Free: Determine Total Float & Free Float (AKA "Slack") of activities in a network diagram](https://www.youtube.com/watch?v=NDa-Fq5jeuM)
package main

import (
	"fmt"
	"time"

	"github.com/gtantech/pdm"
	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/dependency/table"
	"github.com/gtantech/pdm/enums"
	"github.com/gtantech/pdm/relationship"
)

type Attributes struct {
	activity.Data
	Name string
}

func main() {
	project := pdm.New[Attributes]()

	A := activity.New(Attributes{Data: activity.NewData(3 * time.Hour), Name: "A"})
	B := activity.New(Attributes{Data: activity.NewData(4 * time.Hour), Name: "B"})
	C := activity.New(Attributes{Data: activity.NewData(2 * time.Hour), Name: "C"})
	D := activity.New(Attributes{Data: activity.NewData(5 * time.Hour), Name: "D"})
	E := activity.New(Attributes{Data: activity.NewData(1 * time.Hour), Name: "E"})
	F := activity.New(Attributes{Data: activity.NewData(2 * time.Hour), Name: "F"})
	G := activity.New(Attributes{Data: activity.NewData(4 * time.Hour), Name: "G"})
	H := activity.New(Attributes{Data: activity.NewData(3 * time.Hour), Name: "H"})

	//create project dependency table

	tbl := table.New[Attributes]()
	tbl.UpdatePredecessors(B, []table.PredecessorDependency[Attributes]{table.NewPredecessorDependency(A, relationship.New(enums.FS))})
	tbl.UpdatePredecessors(C, []table.PredecessorDependency[Attributes]{table.NewPredecessorDependency(A, relationship.New(enums.FS))})
	tbl.UpdatePredecessors(D, []table.PredecessorDependency[Attributes]{table.NewPredecessorDependency(B, relationship.New(enums.FS))})
	tbl.UpdatePredecessors(E, []table.PredecessorDependency[Attributes]{table.NewPredecessorDependency(C, relationship.New(enums.FS))})
	tbl.UpdatePredecessors(F, []table.PredecessorDependency[Attributes]{table.NewPredecessorDependency(C, relationship.New(enums.FS))})
	tbl.UpdatePredecessors(G, []table.PredecessorDependency[Attributes]{table.NewPredecessorDependency(D, relationship.New(enums.FS)), table.NewPredecessorDependency(E, relationship.New(enums.FS))})
	tbl.UpdatePredecessors(H, []table.PredecessorDependency[Attributes]{table.NewPredecessorDependency(F, relationship.New(enums.FS)), table.NewPredecessorDependency(G, relationship.New(enums.FS))})

	// add dependencies to pdm
	project.AddDependenciesFromTable(tbl)

	// print critical
	fmt.Printf("The total float and free float of each activity:\n")
	for activity := range project.Activities() {
		fmt.Printf("Activity %v:\n", activity.Data().Name)
		fmt.Printf("Total Float:%-5v \tFree Float:%-5v\n\n",
			project.TotalFloat(activity), project.FreeFloat(activity))
	}

}
