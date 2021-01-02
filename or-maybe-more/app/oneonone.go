// oneonone.go
// This file contains the routines for our one-on-one protocol.

package app

type OneOnOneSession struct {
	parties map[HashString]*Client
}
