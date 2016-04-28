package freeboard

import "github.com/gopherjs/gopherjs/js"

// FreeBoard is the Freeboard JS object
var FreeBoard *js.Object

func init() {
	FreeBoard = js.Global.Get("freeboard")
}

// Initialize is called with:
//  * allowEdit (whether to permit the board to be edited)
//  * finished (callback when loading is finished)
func Initialize(allowEdit bool, finished func()) {
	FreeBoard.Call("initialize", allowEdit, finished)
}

// Serialize returns a serialised object of the current board.
func Serialize() *js.Object {
	return FreeBoard.Call("serialize")
}

// LoadDashboard accepts a serialised dashboard and a callback for when loading completes.
func LoadDashboard(serialised *js.Object, callback func()) {
	FreeBoard.Call("loadDashboard", serialised, callback)
}

// LoadDatasourcePlugin accepts a datasource plugin and loads it.
func LoadDatasourcePlugin(ds DsPluginDefinition) {
	FreeBoard.Call("loadDatasourcePlugin", ds.ToFBInterface())
}

// LoadWidgetPlugin accepts a widget plugin and loads it.
func LoadWidgetPlugin(wt WtPluginDefinition) {
	FreeBoard.Call("loadWidgetPlugin", wt.ToFBInterface())
}
