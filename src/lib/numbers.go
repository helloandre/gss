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

package lib

import (
    "text/scanner"
    "strconv"
    "bytes"
    // "fmt"
)

// collect the next characters in s that are numbers
// assumes at least TokenText is a number
func Number(s *scanner.Scanner) string {
    var buffer bytes.Buffer
    
    for s.Scan() != scanner.EOF {
        buffer.WriteString(s.TokenText())
        
        if _, ok := strconv.Atoi(string(s.Peek())); ok != nil {
            return buffer.String()
        }
    }
    
    // reached the end of file, but that's ok potentially?
    return buffer.String()
}

func DetectNumberCollectable(str string, next rune) bool {
    n := string(next)
    // fail fast, fail often
    if len(str) > 1 || len(n) > 1 { return false }
    
    // fucking error returns....
    _, ok_str := strconv.Atoi(str)
    _, ok_next := strconv.Atoi(n)
    
    return ok_str == nil && ok_next == nil
}