package main

/*
#include <stdlib.h>
#include <glib-object.h>
#include <gtk/gtk.h>

static inline GType gvalue_get_type(GValue *v) {
	return G_VALUE_TYPE(v);
}

static inline GType gtype_get_fundamental(GType t) {
	return G_TYPE_FUNDAMENTAL(t);
}

extern void closureMarshal(GClosure*, GValue*, guint, GValue*, gpointer, gpointer);

GClosure* new_closure(void *data) {
	GClosure *closure = g_closure_new_simple(sizeof(GClosure), NULL);
	g_closure_set_meta_marshal(closure, data, (GClosureMarshal)(closureMarshal));
	return closure;
}

*/
import "C"

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

// signal connect

var refHolder []interface{}
var refHolderLock sync.Mutex

var connect = Connect

func Connect(obj interface{}, signal string, cb interface{}) uint64 {
	cbp := &cb
	refHolderLock.Lock()
	refHolder = append(refHolder, cbp) //FIXME deref
	refHolderLock.Unlock()
	closure := C.new_closure(unsafe.Pointer(cbp))
	cSignal := (*C.gchar)(unsafe.Pointer(C.CString(signal)))
	defer C.free(unsafe.Pointer(cSignal))
	id := C.g_signal_connect_closure(C.gpointer(unsafe.Pointer(reflect.ValueOf(obj).Pointer())),
		cSignal, closure, C.gboolean(0))
	return uint64(id)
}

func fromGValue(v *C.GValue) (ret interface{}) {
	valueType := C.gvalue_get_type(v)
	fundamentalType := C.gtype_get_fundamental(valueType)
	switch fundamentalType {
	case C.G_TYPE_OBJECT:
		ret = unsafe.Pointer(C.g_value_get_object(v))
	case C.G_TYPE_STRING:
		ret = fromGStr(C.g_value_get_string(v))
	case C.G_TYPE_UINT:
		ret = int(C.g_value_get_uint(v))
	case C.G_TYPE_BOXED:
		ret = unsafe.Pointer(C.g_value_get_boxed(v))
	case C.G_TYPE_BOOLEAN:
		ret = C.g_value_get_boolean(v) == C.gboolean(1)
	case C.G_TYPE_ENUM:
		ret = int(C.g_value_get_enum(v))
	default:
		panic(fmt.Sprintf("from type %s %T", fromGStr(C.g_type_name(fundamentalType)), v))
	}
	return
}

// string

func fromGStr(s *C.gchar) string {
	return C.GoString((*C.char)(unsafe.Pointer(s)))
}

var _gstrs = make(map[string]*C.gchar)

func toGStr(s string) *C.gchar {
	if gstr, ok := _gstrs[s]; ok {
		return gstr
	}
	gstr := (*C.gchar)(unsafe.Pointer(C.CString(s)))
	_gstrs[s] = gstr
	return gstr
}

// asXXX

func asContainer(o interface{}) *C.GtkContainer {
	return (*C.GtkContainer)(unsafe.Pointer(reflect.ValueOf(o).Pointer()))
}

func asGrid(o interface{}) *C.GtkGrid {
	return (*C.GtkGrid)(unsafe.Pointer(reflect.ValueOf(o).Pointer()))
}

func asNotebook(o interface{}) *C.GtkNotebook {
	return (*C.GtkNotebook)(unsafe.Pointer(reflect.ValueOf(o).Pointer()))
}

func asLabel(o interface{}) *C.GtkLabel {
	return (*C.GtkLabel)(unsafe.Pointer(reflect.ValueOf(o).Pointer()))
}

func asMisc(o interface{}) *C.GtkMisc {
	return (*C.GtkMisc)(unsafe.Pointer(reflect.ValueOf(o).Pointer()))
}
