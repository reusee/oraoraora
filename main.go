package main

/*
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>
#cgo pkg-config: gtk+-3.0 webkit2gtk-3.0
*/
import "C"
import "fmt"

var p = fmt.Printf

func main() {
	// init gtk
	var argc C.int
	C.gtk_init(&argc, nil)

	// root window
	win := C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
	connect(win, "destroy", func() { // quit gtk when window is closed
		C.gtk_main_quit()
	})

	// root grid
	grid := C.gtk_grid_new()
	C.gtk_container_add(asContainer(win), grid)

	// webview container
	pages := C.gtk_notebook_new()
	C.gtk_grid_attach(asGrid(grid), pages, 0, 0, 1, 1)
	C.gtk_widget_set_vexpand(pages, C.gtk_true())
	C.gtk_widget_set_hexpand(pages, C.gtk_true())
	C.gtk_notebook_set_tab_pos(asNotebook(pages), C.GTK_POS_LEFT)

	// new view constructor
	var newView func() *View
	newView = func() *View {
		view := NewView()

		label := C.gtk_label_new(nil)
		C.gtk_label_set_use_markup(asLabel(label), C.gtk_true())
		C.gtk_misc_set_alignment(asMisc(label), 0, 0.5)

		C.gtk_notebook_append_page(asNotebook(pages), view.Widget, label)

		// new view is requested
		connect(view.View, "create", func() *C.GtkWidget {
			return newView().Widget
		})

		// set tab label to page title
		connect(view.View, "notify::title", func() {
			var title string
			// trim long title
			for _, r := range fromGStr(C.webkit_web_view_get_title(view.View)) {
				if len(title) > 32 {
					break
				}
				title += string(r)
			}
			C.gtk_label_set_markup(asLabel(label), toGStr(fmt.Sprintf(`<span font="10">%s</span>`, title)))
		})

		return view
	}

	// first view
	view := newView()
	C.webkit_web_view_load_uri(view.View, toGStr("http://www.bilibili.tv"))

	// show window and run
	C.gtk_widget_show_all(win)
	C.gtk_main()
}
