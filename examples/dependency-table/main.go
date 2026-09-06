// example taken from [Engineer4Free: Determine Late Start (LS) and Late Finish (LF) of acitivies in PDM network diagram](https://www.youtube.com/watch?v=bobDz_Bh3tI)
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

	A := activity.New(Attributes{Data: activity.NewData(5 * time.Hour), Name: "A"})
	B := activity.New(Attributes{Data: activity.NewData(4 * time.Hour), Name: "B"})
	C := activity.New(Attributes{Data: activity.NewData(5 * time.Hour), Name: "C"})
	D := activity.New(Attributes{Data: activity.NewData(6 * time.Hour), Name: "D"})
	E := activity.New(Attributes{Data: activity.NewData(3 * time.Hour), Name: "E"})
	F := activity.New(Attributes{Data: activity.NewData(4 * time.Hour), Name: "F"})

	//create project dependency table

	tbl := table.New[Attributes]()
	tbl.UpdatePredecessors(B, []table.PredecessorDependency[Attributes]{table.NewPredecessorDependency(A, relationship.New(enums.FS))})
	tbl.UpdatePredecessors(C, []table.PredecessorDependency[Attributes]{table.NewPredecessorDependency(A, relationship.New(enums.FS))})
	tbl.UpdatePredecessors(D, []table.PredecessorDependency[Attributes]{table.NewPredecessorDependency(B, relationship.New(enums.FS))})
	tbl.UpdatePredecessors(E, []table.PredecessorDependency[Attributes]{table.NewPredecessorDependency(C, relationship.New(enums.FS))})
	tbl.UpdatePredecessors(F, []table.PredecessorDependency[Attributes]{table.NewPredecessorDependency(D, relationship.New(enums.FS)), table.NewPredecessorDependency(E, relationship.New(enums.FS))})

	// add dependencies to pdm
	project.AddDependenciesFromTable(tbl)

	// print the early/late start/finish of each activity
	for _, activity := range []activity.Activity[Attributes]{A, B, C, D, E, F} {
		fmt.Printf("Activity %v:\n", activity.Data().Name)
		fmt.Printf("Early Start:%-5v \tEarly Finish:%-5v\n",
			activity.Early().Start(), activity.Early().Finish())
		fmt.Printf("Late Start:%-5v \tLate Finish:%-5v\n\n",
			activity.Late().Start(), activity.Late().Finish())
	}

}
