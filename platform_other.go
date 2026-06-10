//go:build !darwin && !linux

package main

const windowModeConfigurable = true

func hideSystemUI()    {}
func restoreSystemUI() {}
