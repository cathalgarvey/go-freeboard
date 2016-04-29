package freeboard

import "github.com/gopherjs/gopherjs/js"

// FBWrapper is a wrapper for the Freeboard object.
// Use the "FB" variable instantiated at initialisation.
type FBWrapper struct {
	FreeboardObject *js.Object
}

// FB is the Freeboard JS object as wrapped from the global context.
var FB *FBWrapper

func init() {
	fbobj := js.Global.Get("freeboard")
	FB = &FBWrapper{fbobj}
}

// Initialize is called with:
//  * allowEdit (whether to permit the board to be edited)
//  * finished (callback when loading is finished)
func (fb *FBWrapper) Initialize(allowEdit bool, finished func()) {
	fb.FreeboardObject.Call("initialize", allowEdit, finished)
}

// Serialize returns a serialised object of the current board.
func (fb *FBWrapper) Serialize() *js.Object {
	return fb.FreeboardObject.Call("serialize")
}

// LoadDashboard accepts a serialised dashboard and a callback for when loading completes.
func (fb *FBWrapper) LoadDashboard(serialised *js.Object, callback func()) {
	fb.FreeboardObject.Call("loadDashboard", serialised, callback)
}

// LoadDatasourcePlugin accepts a datasource plugin and loads it.
func (fb *FBWrapper) LoadDatasourcePlugin(ds DsPluginDefinition) {
	fb.FreeboardObject.Call("loadDatasourcePlugin", ds.ToFBInterface())
}

// LoadWidgetPlugin accepts a widget plugin and loads it.
func (fb *FBWrapper) LoadWidgetPlugin(wt WtPluginDefinition) {
	fb.FreeboardObject.Call("loadWidgetPlugin", wt.ToFBInterface())
}
