// Package liftover provides liftover functionality via running a pyliftover script.
package liftover

// Converter defines a set of methods for types implementing Converter.
type Converter interface {
	Convert38to19(cwd string, inputFile string, outputFile string) error
	Convert19to38(cwd string, inputFile string, outputFile string) error
}
