// Copyright 2019 goodxp(goodxp@gmail.com). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package wqSGF is a SGF parser for Go game(weiqi).
//  - load and save .sgf files (from/to game tree structure for further coding)
//  - encode/decode SGF strings
//  - helper functions for convertion of SGF property value types
package wqSGF

import (
	"actTree"
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
)

// load a .sfg file and parse it into a tree structure
// assumed charset 'UTF-8'
func Load(file string) *actTree.Tree {
	sgf, _ := ioutil.ReadFile(file)
	return Parse(string(sgf))
}

// save kifu(game record with a tree structure) into a .sgf file
// used charset 'UTF-8'
func Save(kifu *actTree.Tree, file string) {
	ioutil.WriteFile(file, ([]byte)(ToSGF(kifu)), 0644)
}

func ToSGF(kifu *actTree.Tree) string {
	var buf bytes.Buffer
	buf.WriteString("(")
	enterNode := func(n *actTree.Node) bool {
		if n.HasSib() {
			buf.WriteString("(")
		}
		buf.WriteString(n.Value.(*Node).ToSGF())
		return false
	}
	leaveNode := func(n *actTree.Node) bool {
		if n.HasSib() {
			buf.WriteString(")")
		}
		return false
	}
	actTree.WalkThrough(kifu.Root, enterNode, leaveNode)
	buf.WriteString(")")
	return buf.String()
}

func Parse(sgf string) *actTree.Tree {
	kifu := actTree.New()
	var stack []*actTree.Node
	var node *actTree.Node

	onTree := func(begin bool) bool {
		if begin { //push: node => stack
			stack = append(stack, node)
		} else if len(stack) > 0 { //pop: node <= stack
			n := len(stack) - 1
			node = stack[n]
			stack = stack[:n]
		}
		return false
	}
	onNode := func(sgfNode string) bool {
		n := new(Node)
		n.FromSGF(sgfNode)
		node = kifu.Add(n, node)
		return false
	}
	Scan(sgf, onTree, onNode)
	return kifu
}

func Scan(sgf string,
	onTree func(begin bool) bool,
	onNode func(sgfNode string) bool) {

	for _, sn := range matchSGF(reTree, sgf) {
		switch sn {
		case "(":
			if onTree != nil && onTree(true) {
				return
			}
		case ")":
			if onTree != nil && onTree(false) {
				return
			}
		default:
			if onNode != nil && onNode(sn) {
				return
			}
		}
	}
}

const (
	reTree  = `\(|\)|(;(\s*[A-Z]+(\s*((\[\])|(\[(.|\s)*?([^\\]\]))))+)*)`
	reNode  = `[A-Z]+(\s*((\[\])|(\[(.|\s)*?([^\\]\]))))+`
	reIdent = `[A-Z]+`
	reVals  = `(\[\])|(\[(.|\s)*?([^\\]\]))`
)

func matchSGF(re string, sgf string) []string {
	r, err := regexp.Compile(re)
	if err != nil {
		fmt.Errorf("compile regexp %q failed. error = %v", re, err)
		return nil
	}
	return r.FindAllString(sgf, -1)
}

type Node struct {
	Props []Prop
}

func (n *Node) FromSGF(sn string) {
	for _, sp := range matchSGF(reNode, sn) {
		var p Prop
		p.FromSGF(sp)
		n.Props = append(n.Props, p)
	}
}

func (n *Node) ToSGF() string {
	var buf bytes.Buffer
	buf.WriteString(";")
	for _, p := range n.Props {
		buf.WriteString(p.ToSGF())
	}
	return buf.String()
}

//Vals are with their SGF format: E.g. "[value]", "[value:value]"
//Vals == nil means the Prop format is "Id"
//Vals[0] == "" means the Prop format is "Id[]"
type Prop struct {
	Id   string
	Vals []string
}

func (p *Prop) FromSGF(sp string) {
	si := matchSGF(reIdent, sp)
	p.Id = si[0]
	for _, sv := range matchSGF(reVals, sp) {
		p.Vals = append(p.Vals, sv)
	}
}

func (p *Prop) ToSGF() string {
	var buf bytes.Buffer
	buf.WriteString(p.Id)
	for _, v := range p.Vals {
		buf.WriteString(v)
	}
	return buf.String()
}
