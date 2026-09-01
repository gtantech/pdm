package main

import (
	"fmt"
	"time"

	"github.com/gtantech/pdm"
	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/activity/data"
	"github.com/gtantech/pdm/dependency"
	"github.com/gtantech/pdm/enums"
)

type Data struct {
	data.ActivityData
	Name string
}

func main() {
	A := activity.New(Data{ActivityData: data.New(3 * time.Hour), Name: "A"})
	B := activity.New(Data{ActivityData: data.New(4 * time.Hour), Name: "B"})
	C := activity.New(Data{ActivityData: data.New(2 * time.Hour), Name: "C"})
	D := activity.New(Data{ActivityData: data.New(5 * time.Hour), Name: "D"})
	E := activity.New(Data{ActivityData: data.New(1 * time.Hour), Name: "E"})
	F := activity.New(Data{ActivityData: data.New(2 * time.Hour), Name: "F"})
	G := activity.New(Data{ActivityData: data.New(4 * time.Hour), Name: "G"})
	H := activity.New(Data{ActivityData: data.New(3 * time.Hour), Name: "H"})

	tableOfDependencies := make(map[activity.Activity[Data]][]activity.Activity[Data])
	tableOfDependencies[A] = []activity.Activity[Data]{}     // A has no dependencies
	tableOfDependencies[B] = []activity.Activity[Data]{A}    // B depends on A
	tableOfDependencies[C] = []activity.Activity[Data]{A}    // C depends on A
	tableOfDependencies[D] = []activity.Activity[Data]{B}    //		.
	tableOfDependencies[E] = []activity.Activity[Data]{C}    //		.
	tableOfDependencies[F] = []activity.Activity[Data]{C}    //		.
	tableOfDependencies[G] = []activity.Activity[Data]{D, E} // G depends on D and E
	tableOfDependencies[H] = []activity.Activity[Data]{F, G} // H depends on F and G

	project := pdm.New[Data]()

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
