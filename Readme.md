# Freeboard bindings for GopherJS
by Cathal Garvey, ©2016, released under the GNU AGPL version 3 or later.

[![Documentation on Godoc](https://godoc.org/github.com/cathalgarvey/go-freeboard?status.svg)](https://godoc.org/github.com/cathalgarvey/go-freeboard)

This is a set of bindings to the Freeboard dashboard framework, written in Go for GopherJS. It also includes a simple framework for plugins, comprising an interface for plugin objects and a struct that simplifies the process of defining a plugin for inclusion in Freeboard.

## How to write a plugin using Golang
See the `testplugin` folder for a silly datasource example about cats. This works, for the most part.

I have not yet tested making a widget plugin using this framework.

As always, caveat emptor.
