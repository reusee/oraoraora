package main

/*
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>
#cgo pkg-config: gtk+-3.0 webkit2gtk-3.0
*/
import "C"
import (
	"fmt"
	"runtime"
)

var p = fmt.Printf

func init() {
	runtime.GOMAXPROCS(32)
}

func main() {
	var argc C.int
	C.gtk_init(&argc, nil)

	win := C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
	connect(win, "destroy", func() {
		C.gtk_main_quit()
	})

	grid := C.gtk_grid_new()
	C.gtk_container_add(asContainer(win), grid)

	pages := C.gtk_notebook_new()
	C.gtk_grid_attach(asGrid(grid), pages, 0, 0, 1, 1)
	C.gtk_widget_set_vexpand(pages, C.gtk_true())
	C.gtk_widget_set_hexpand(pages, C.gtk_true())
	C.gtk_notebook_set_tab_pos(asNotebook(pages), C.GTK_POS_LEFT)

	var newView func() *View
	newView = func() *View {
		view := NewView()

		label := C.gtk_label_new(nil)
		C.gtk_misc_set_alignment(asMisc(label), 0, 0.5)

		C.gtk_notebook_append_page(asNotebook(pages), view.Widget, label)

		connect(view.View, "ready-to-show", func() {
			C.gtk_widget_show_all(view.Widget)
		})
		connect(view.View, "create", func() *C.GtkWidget {
			return newView().Widget
		})
		connect(view.View, "notify::title", func() {
			var title string
			for _, r := range fromGStr(C.webkit_web_view_get_title(view.View)) {
				if len(title) > 32 {
					break
				}
				title += string(r)
			}
			C.gtk_label_set_text(asLabel(label), toGStr(title))
		})

		return view
	}

	view := newView()
	C.webkit_web_view_load_uri(view.View, toGStr("http://www.bilibili.tv"))
	connect(view.View, "load-changed", func(_, ev interface{}) {
		p("load changed %d\n", ev.(int))
	})

	C.gtk_widget_show_all(win)
	C.gtk_main()
}
