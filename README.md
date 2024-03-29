# foptgen

foptgen is a binary tool that generates functional options for Go structs.

## Install

```bash
go install github.com/200sc/foptgen/main/foptgen
```

## Usage 

```bash
$ foptgen --dir=internal/components/titlebar --struct=Constructor --overwrite=false
Output target internal/components/titlebar/opts_gen.go already exists; overwrite? (y/N)
y
```

Input struct:

```go
type Constructor struct {
	Color        color.Color
	Height       float64
	Layers       []int
	Fire         bool
	StringArrays [][][]string
}
```

Output (to file `internal\components\titlebar\opts_gen.go`):

```go
// Code generated by foptgen; DO NOT EDIT.

package titlebar

import "image/color"

type Option func(Constructor) Constructor

func WithColor(v color.Color) Option {
	return func(s Constructor) Constructor {
		s.Color = v
		return s
	}
}

func WithHeight(v float64) Option {
	return func(s Constructor) Constructor {
		s.Height = v
		return s
	}
}

func WithLayers(v []int) Option {
	return func(s Constructor) Constructor {
		s.Layers = v
		return s
	}
}

func WithFire(v bool) Option {
	return func(s Constructor) Constructor {
		s.Fire = v
		return s
	}
}

func WithStringArrays(v [][][]string) Option {
	return func(s Constructor) Constructor {
		s.StringArrays = v
		return s
	}
}
```
