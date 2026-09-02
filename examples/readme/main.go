package main

import (
	"fmt"
	"time"

	"github.com/gtantech/pdm"
	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/dependency"
	"github.com/gtantech/pdm/enums"
)

type Attributes struct {
	activity.Data
	Name string
}

func main() {
	A := activity.New(Attributes{Data: activity.NewData(3 * time.Hour), Name: "A"})
	B := activity.New(Attributes{Data: activity.NewData(4 * time.Hour), Name: "B"})
	C := activity.New(Attributes{Data: activity.NewData(2 * time.Hour), Name: "C"})
	D := activity.New(Attributes{Data: activity.NewData(5 * time.Hour), Name: "D"})
	E := activity.New(Attributes{Data: activity.NewData(1 * time.Hour), Name: "E"})
	F := activity.New(Attributes{Data: activity.NewData(2 * time.Hour), Name: "F"})
	G := activity.New(Attributes{Data: activity.NewData(4 * time.Hour), Name: "G"})
	H := activity.New(Attributes{Data: activity.NewData(3 * time.Hour), Name: "H"})

	tableOfDependencies := make(map[activity.Activity[Attributes]][]activity.Activity[Attributes])
	tableOfDependencies[A] = []activity.Activity[Attributes]{}     // A has no dependencies
	tableOfDependencies[B] = []activity.Activity[Attributes]{A}    // B depends on A
	tableOfDependencies[C] = []activity.Activity[Attributes]{A}    // C depends on A
	tableOfDependencies[D] = []activity.Activity[Attributes]{B}    //		.
	tableOfDependencies[E] = []activity.Activity[Attributes]{C}    //		.
	tableOfDependencies[F] = []activity.Activity[Attributes]{C}    //		.
	tableOfDependencies[G] = []activity.Activity[Attributes]{D, E} // G depends on D and E
	tableOfDependencies[H] = []activity.Activity[Attributes]{F, G} // H depends on F and G

	project := pdm.New[Attributes]()

	for key := range tableOfDependencies {
		successor := key
		for _, predecessor := range tableOfDependencies[successor] {
			project.AddDependency(predecessor, successor, dependency.New(enums.FS))
		}
	}

	project.UpdateActivityTimestamps()

	for key := range tableOfDependencies {
		fmt.Printf("Activity %v:\n", key.Data().Name)
		fmt.Printf("Early Start:%-5v \tEarly Finish:%-5v\n", key.Timestamps().Early().Start(), key.Timestamps().Early().Finish())
		fmt.Printf("Late Start:%-5v \tLate Finish:%-5v\n\n", key.Timestamps().Late().Start(), key.Timestamps().Late().Finish())
	}

}
