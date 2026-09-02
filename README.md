# Precedence Diagram Method for Go
This is a Go package that provides an implementation to help plan a project via [PDM](https://en.wikipedia.org/wiki/Precedence_diagram_method).

## Table of Contents
- [Install](#install)
- [Example](#example)
- [License](#license)

## Install
Install via `go get`. Note that Go 1.23 or newer is required.

```sh
# After: go mod init ...
go get -u github.com/gtantech/pdm/v2
```

## Example
package main

import (
	"fmt"
	"time"

	"github.com/gtantech/pdm/v2"
	"github.com/gtantech/pdm/v2/activity"
	"github.com/gtantech/pdm/v2/dependency"
	"github.com/gtantech/pdm/v2/enums"
)

type Attributes struct {
	activity.Data
	Name string
}

func main() {
	project := pdm.New[Attributes]()
	// add activities (optional)
	A := project.AddActivity(activity.New(Attributes{Data: activity.NewData(5 * time.Hour), Name: "A"}))
	B := project.AddActivity(activity.New(Attributes{Data: activity.NewData(4 * time.Hour), Name: "B"}))
	C := project.AddActivity(activity.New(Attributes{Data: activity.NewData(5 * time.Hour), Name: "C"}))
	D := project.AddActivity(activity.New(Attributes{Data: activity.NewData(6 * time.Hour), Name: "D"}))
	E := project.AddActivity(activity.New(Attributes{Data: activity.NewData(3 * time.Hour), Name: "E"}))
	F := project.AddActivity(activity.New(Attributes{Data: activity.NewData(4 * time.Hour), Name: "F"}))

	// add dependencies to pdm
	//                                                    // A has no dependencies
	project.AddDependency(A, B, dependency.New(enums.FS)) // B depends on A
	project.AddDependency(A, C, dependency.New(enums.FS)) // C depends on A
	project.AddDependency(B, D, dependency.New(enums.FS)) //       .
	project.AddDependency(C, E, dependency.New(enums.FS)) //       .
	project.AddDependency(D, F, dependency.New(enums.FS)) // F depends on D and E
	project.AddDependency(E, F, dependency.New(enums.FS))

	// update the early/late start/finish of each activity
	project.UpdateActivityTimestamps()

	// print the early/late start/finish of each activity
	for _, activity := range []activity.Activity[Attributes]{A, B, C, D, E, F} {
		fmt.Printf("Activity %v:\n", activity.Data().Name)
		fmt.Printf("Early Start:%-5v \tEarly Finish:%-5v\n",
			activity.Timestamps().Early().Start(), activity.Timestamps().Early().Finish())
		fmt.Printf("Late Start:%-5v \tLate Finish:%-5v\n\n",
			activity.Timestamps().Late().Start(), activity.Timestamps().Late().Finish())
	}

}

```

## License

Licensed under [MIT License](./LICENSE)

## Thanks!

Thanks for reading and happy coding! Add a star to the project if you find it useful!