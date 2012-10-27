/**
 * Copyright (c) 2012 Andre Bluehs <hello@andrebluehs.net>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy 
 * of this software and associated documentation files (the "Software"), to deal 
 * in the Software without restriction, including without limitation the rights 
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies 
 * of the Software, and to permit persons to whom the Software is furnished to do so, 
 * subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all 
 * copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR 
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS 
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR 
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER 
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN 
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package tokenize_css

import (
    "bufio"
    "bytes"
    "fmt"
    "text/scanner"
    "../../lib"
)

var (
    inBraces = false
    inColon = false
    inUrl = false
    preserveSpacing = false
    isSpace = false
    s scanner.Scanner
    write_this bytes.Buffer
)

/**
 * First pass to build tokens, state machine wooo!
 */
func Parse(r *bufio.Reader) []lib.CSSToken {
    
    s.Init(r)
    // ignore newlines and some whitespaces
    s.Whitespace = 1<<'\t' | 1<<'\r' | 1<<'\n'
    // ignore ints and floats because it gets confused about octals
    s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments
    
    tokens := []lib.CSSToken{}
    
    for s.Scan() != scanner.EOF {
        str := s.TokenText()
        // next := s.Peek()
        
        if s.ErrorCount > 0 {
            // scanner emits it's own error message here
            panic("error parsing file")
        }
        
        // some comments are ok, some we output, some css just doesn't support so why are they in there...?
        if lib.IsComment(str){
            if lib.IsInlineComment(str) {
                panic(fmt.Sprintf("inline comment detected at %d", s.Pos().Line))
            }
            if !lib.ShouldPreserveComment(str) {
                continue
            }
        }
        
        // nuke spaces before: space, semicolon, open brace
        // or if we explicitly want to
        if str == " " && !preserveSpacing {
            continue
        }   
        
        // only turn on space preserving if we're inside colons or outside braces
        // with exceptions to both of those (see above)
        if (inColon && str != " ") || (!inBraces && str != " ") {
            preserveSpacing = true
        }
        
        // using a rather large switch/case here because we don't lex this
        switch str {
        case "{":
            inBraces = true
            preserveSpacing = false
        case "}":
            inBraces = false
            preserveSpacing = false
        case ":":
            if inBraces {
                inColon = true
            }
        case ";":
            inColon = false
            preserveSpacing = false
        case "(":
            inUrl = true
        case ")":
            inUrl = false
        }    
        
        // append that guy to our slice, preserving some state
        // we don't care about depth, so ignore that guy
        tokens = append(tokens, lib.CSSToken{str, s.Pos().Line, inBraces, inColon, inUrl})
    }
    
    if inBraces {
        panic("mismatched braces")
    }
    
    return tokens
}