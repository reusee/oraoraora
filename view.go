package main

/*
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>
*/
import "C"
import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"unsafe"
)

var cookieFilePath string

func init() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	cookieFilePath = filepath.Join(user.HomeDir, ".cookies")
}

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

	// ready to show
	connect(view.View, "ready-to-show", func() {
		C.gtk_widget_show_all(view.Widget)
	})

	// page load state changed
	connect(view.View, "load-changed", func(_, ev interface{}) {
		p("load changed %d\n", ev.(int))
	})

	// context
	context := C.webkit_web_view_get_context(view.View)
	C.webkit_web_context_set_spell_checking_enabled(context, C.gtk_false())
	C.webkit_web_context_set_tls_errors_policy(context, C.WEBKIT_TLS_ERRORS_POLICY_IGNORE)
	C.webkit_web_context_set_disk_cache_directory(context, toGStr(os.TempDir()))

	// settings
	settings := C.webkit_web_view_get_settings(view.View)
	C.webkit_settings_set_enable_java(settings, C.gtk_false())
	C.webkit_settings_set_enable_tabs_to_links(settings, C.gtk_false())
	C.webkit_settings_set_enable_dns_prefetching(settings, C.gtk_true())
	C.webkit_settings_set_javascript_can_access_clipboard(settings, C.gtk_true())
	C.webkit_settings_set_enable_site_specific_quirks(settings, C.gtk_true())
	C.webkit_settings_set_enable_smooth_scrolling(settings, C.gtk_true())

	// handle cookie
	cookieManager := C.webkit_web_context_get_cookie_manager(context)
	C.webkit_cookie_manager_set_persistent_storage(cookieManager, toGStr(cookieFilePath),
		C.WEBKIT_COOKIE_PERSISTENT_STORAGE_TEXT)

	return view
}
