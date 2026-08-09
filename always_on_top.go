//go:build linux && !wayland

package main

/*
#cgo LDFLAGS: -lX11
#include <stdlib.h>
#include <unistd.h>
#include <X11/Xatom.h>
#include <X11/Xlib.h>

static Window findWindowForPID(Display *display, Window parent, unsigned long pid) {
	Atom pidAtom = XInternAtom(display, "_NET_WM_PID", True);
	if (pidAtom != None) {
		Atom actualType;
		int actualFormat;
		unsigned long itemCount;
		unsigned long bytesAfter;
		unsigned char *property = NULL;

		if (XGetWindowProperty(display, parent, pidAtom, 0, 1, False,
			XA_CARDINAL, &actualType, &actualFormat, &itemCount, &bytesAfter,
			&property) == Success) {
			if (property != NULL && itemCount == 1 &&
				*(unsigned long *)property == pid) {
				XFree(property);
				return parent;
			}
			if (property != NULL) {
				XFree(property);
			}
		}
	}

	Window root;
	Window ignoredParent;
	Window *children = NULL;
	unsigned int childCount = 0;
	if (!XQueryTree(display, parent, &root, &ignoredParent, &children, &childCount)) {
		return None;
	}

	Window found = None;
	for (unsigned int i = 0; i < childCount && found == None; i++) {
		found = findWindowForPID(display, children[i], pid);
	}
	if (children != NULL) {
		XFree(children);
	}
	return found;
}

static int setWindowAlwaysOnTop(int enabled) {
	Display *display = XOpenDisplay(NULL);
	if (display == NULL) {
		return 1;
	}

	Window root = DefaultRootWindow(display);
	Window window = findWindowForPID(display, root, (unsigned long)getpid());
	if (window == None) {
		XCloseDisplay(display);
		return 2;
	}

	Atom stateAtom = XInternAtom(display, "_NET_WM_STATE", False);
	Atom aboveAtom = XInternAtom(display, "_NET_WM_STATE_ABOVE", False);
	XEvent event = {0};
	event.xclient.type = ClientMessage;
	event.xclient.window = window;
	event.xclient.message_type = stateAtom;
	event.xclient.format = 32;
	event.xclient.data.l[0] = enabled ? 1 : 0;
	event.xclient.data.l[1] = aboveAtom;
	event.xclient.data.l[3] = 1;

	int sent = XSendEvent(display, root, False,
		SubstructureRedirectMask | SubstructureNotifyMask, &event);
	XFlush(display);
	XCloseDisplay(display);
	return sent == 0 ? 3 : 0;
}
*/
import "C"

import "fmt"

// setAlwaysOnTop is an OS-specific implementation for Linux applications
// running through X11 or XWayland.
func setAlwaysOnTop(enabled bool) error {
	requested := C.int(0)
	if enabled {
		requested = 1
	}

	switch result := C.setWindowAlwaysOnTop(requested); result {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("cannot connect to the X11 display")
	case 2:
		return fmt.Errorf("cannot find the Hanji window")
	default:
		return fmt.Errorf("KWin rejected the window-state request")
	}
}
