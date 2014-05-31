package main

/*
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>
*/
import "C"
import "unsafe"

type View struct {
	Widget *C.GtkWidget
	View   *C.WebKitWebView
}

func NewView() *View {
	widget := C.webkit_web_view_new()
	view := &View{
		Widget: widget,
		View:   (*C.WebKitWebView)(unsafe.Pointer(widget)),
	}
	return view
}
