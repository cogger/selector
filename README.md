# selector 

[![GoDoc](https://godoc.org/github.com/cogger/selector?status.png)](http://godoc.org/github.com/cogger/selector)  
[![Build Status](https://travis-ci.org/cogger/selector.svg?branch=master)](https://travis-ci.org/cogger/selector)  
[![Coverage Status](https://coveralls.io/repos/cogger/selector/badge.svg?branch=master)](https://coveralls.io/r/cogger/selector?branch=master)  
[![License](http://img.shields.io/:license-apache-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

selector adds generic data to contexts

## Usage
~~~ go
// main.go
package main

import (
	"github.com/cogger/selector"
	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1/wait"
	"gopkg.in/cogger/cogger.v1/cogs"
)

func main() {
	ctx := context.Background()
	sel := selector.New()

	sel = sel.Case(func(ctx context.Context)bool{
		//Case 1
		//test for something
		return false
	},cogs.Simple(ctx, func()error{
		//do something
		return nil
	}))

	sel = sel.Case(func(ctx context.Context)bool{
		//Case 2
		//test for something
		return false
	},cogs.Simple(ctx, func()error{
		//do something
		return nil
	}))

	sel = sel.Case(func(ctx context.Context)bool{
		//Case 3
		//test for something
		return true
	},cogs.Simple(ctx, func()error{
		//do something
		return nil
	}))

	sel = sel.Default(cogs.Simple(ctx, func()error{
		return nil
	}))

	errs := wait.Resolve(ctx, sel)
	
}

~~~

