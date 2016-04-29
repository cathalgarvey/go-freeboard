package freeboard

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

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

// NewDashboard clears the contents of the freeboard and initialises a new dashboard
func (fb *FBWrapper) NewDashboard() {
	fb.FreeboardObject.Call("newDashboard")
}

// Serialize returns a serialised object of the current board.
func (fb *FBWrapper) Serialize() *js.Object {
	return fb.FreeboardObject.Call("serialize")
}

// LoadDashboard accepts a serialised dashboard and a callback for when loading completes.
func (fb *FBWrapper) LoadDashboard(serialised, callback *js.Object) {
	fb.FreeboardObject.Call("loadDashboard", serialised, callback)
}

// SetEditing programmatically controls the editing state of the board
func (fb *FBWrapper) SetEditing(editing, animate *js.Object) {
	fb.FreeboardObject.Call("setEditing", editing, animate)
}

// IsEditing returns boolean depending on whether the dashboard is in
// the view-only or edit state.
func (fb *FBWrapper) IsEditing() *js.Object {
	return fb.FreeboardObject.Call("isEditing")
}

// LoadDatasourcePlugin accepts a datasource plugin and loads it.
// This can be passed either a *js.Object for a JS plugin, or a
// map defining a Go plugin; but use LoadGoDatasourcePlugin for that.
func (fb *FBWrapper) LoadDatasourcePlugin(ds interface{}) {
	fb.FreeboardObject.Call("loadDatasourcePlugin", ds)
}

// LoadGoDatasourcePlugin accepts a datasource plugin
// written in Go and loads it.
func (fb *FBWrapper) LoadGoDatasourcePlugin(ds DsPluginDefinition) {
	fb.LoadDatasourcePlugin(ds.ToFBInterface())
}

// LoadWidgetPlugin accepts a widget plugin and loads it.
// This can be passed either a *js.Object for a JS plugin, or a
// map defining a Go plugin; but use LoadGoWidgetPlugin for that.
func (fb *FBWrapper) LoadWidgetPlugin(wt interface{}) {
	fb.FreeboardObject.Call("loadWidgetPlugin", wt)
}

// LoadGoWidgetPlugin accepts a widget plugin written in Go
// and loads it.
func (fb *FBWrapper) LoadGoWidgetPlugin(wt WtPluginDefinition) {
	fb.LoadWidgetPlugin(wt.ToFBInterface())
}

// ShowLoadingIndicator shows or hides the loading indicator.
func (fb *FBWrapper) ShowLoadingIndicator(show *js.Object) {
	fb.FreeboardObject.Call("showLoadingIndicator", show)
}

// ShowDialog shows a styled dialog box with custom content.
//     * contentElement (DOM or jquery element) - The DOM or jquery element to display within the content of the dialog box.
//     * title (string) - The title of the dialog box displayed on the top left.
//     * okButtonTitle (string) - The string to display in the button that will be used as the OK button. A null or undefined value will result in no button being displayed.
//     * cancelButtonTitle (string) - The string to display in the button that will be used as the Cancel button. A null or undefined value will result in no button being displayed.
//     * okCallback (function) - A function that will be called if the user presses the OK button.
func (fb *FBWrapper) ShowDialog(contentElement dom.HTMLElement, title, okButtonTitle, cancelButtonTitle string, okCallback interface{}) {
	fb.FreeboardObject.Call("showDialog", contentElement, title, okButtonTitle, cancelButtonTitle, okCallback)
}

// GetDatasourceSettings returns the current settings for a datasource
// or null if no datasource with the given name is found.
func (fb *FBWrapper) GetDatasourceSettings(name string) *js.Object {
	return fb.FreeboardObject.Call("getDatasourceSettings", name)
}

// SetDatasourceSettings updates settings on a datasource.
func (fb *FBWrapper) SetDatasourceSettings(name string, settings *js.Object) {
	fb.FreeboardObject.Call("setDatasourceSettings", name, settings)
}

// On attaches a callback to a global freeboard event.
// At present, only "dashboard_loaded" and "initialized" are
// fired as events by freeboard.
func (fb *FBWrapper) On(eventName string, callback *js.Object) {
	fb.FreeboardObject.Call("on", eventName, callback)
}
