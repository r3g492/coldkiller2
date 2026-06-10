//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit

#import <AppKit/AppKit.h>

void hideSystemUI() {
	[NSApp setPresentationOptions:
		NSApplicationPresentationHideDock |
		NSApplicationPresentationHideMenuBar];
}

void restoreSystemUI() {
	[NSApp setPresentationOptions:NSApplicationPresentationDefault];
}
*/
import "C"

const windowModeConfigurable = true

func hideSystemUI()    { C.hideSystemUI() }
func restoreSystemUI() { C.restoreSystemUI() }
