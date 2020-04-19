package controller

import (
	"github.com/apache/iotdb-operator/pkg/controller/iotdb"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, iotdb.Add)
}
