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

package compress_css

import (
    "bytes"
    "fmt"
    "strings"
    "strconv"
    "../../lib"
)

var (
    nukeSpaces = []string{"!", "{", "}", ";", ">", "+", "(", ")", ","}
    removeAfterZero = []string{"px", "em", "%", "in", "cm", "mm", "pc", "pt", "ex"}
)

/**
 * Second pass to compress tokens
 */
func Compress(tokens []lib.CSSToken) bytes.Buffer {
    
    var write_this bytes.Buffer
    
    for i:=0;i<len(tokens);i++ {
        token := tokens[i]
        current := token.Str
        
        next := "GO_COMPRESSOR_FAKE_NEXT"
        if i != (len(tokens)-1) {
            next = tokens[i+1].Str
        }
        
        isZero := len(current) == 1 && string(current[0]) == "0"
        
        /**
         * Optimizations
         */
        // we don't need a semicolon before a close brace
        if (current == ";" && next == "}") || (current == ";" && next == ";") {
            continue
        }
        // nuke any leading spaces
        // !{};:>+(),
        if current == " " && lib.Contains(nukeSpaces, next) {
            continue
        } 
        // only remove space before colon in braces to avoid p :link{} => p:link{}
        if (token.InBraces && current == " " && next == ":"){
            continue
        }
        // nuke any trailing spaces
        if next == " " && lib.Contains(nukeSpaces, current) {
            i++
        }
        // 0.x => .x
        if isZero && next == "." {
            continue    
        }
        // 0px => 0
        if isZero && lib.Contains(removeAfterZero, next) {
            i++
        }
        
        /****
         * Commence big, ugly, string-maniuplatey optimizations
         */
        // convert colors from rgb(51,102,153) to #336699 to hopefully benefit from the next optimization
        if current == "rgb" && next == "(" {
            color, length := lib.GetUntil(tokens[i+2:], []string{")"})
            splited := strings.Split(color, ",")
            if len(splited) == 3 {
                first, _ := strconv.ParseInt(splited[0], 10, 8)
                second, _ := strconv.ParseInt(splited[1], 10, 8)
                third, _ := strconv.ParseInt(splited[2], 10, 8)
                
                current = "#" + fmt.Sprintf("%x", first) + fmt.Sprintf("%x", second) + fmt.Sprintf("%x", third)
            }
            i += length+2 // compensate for ( and )
        } 
        
        // dedupe css colors #00FF00 => #0F0
        // assumes there will  be no spaces between numbers
        // must not be in url because not all urls are stings, and can have a hash
        if token.InColon && !token.InUrl && (current == "#" || string(current[0]) == "#") {
            // we did somethign in the last step
            style_str, length := "", 0
            if len(current) == 7 {
                style_str, length = current[1:], 0
                current = "#"
            } else {
                style_str, length = lib.GetUntil(tokens[i+1:], []string{" ", ";", "}"})
            }
            
            if len(style_str) == 6 && string(style_str[0]) == string(style_str[1]) && string(style_str[2]) == string(style_str[3]) && string(style_str[4]) == string(style_str[5]) {
                current = current + string(style_str[0]) + string(style_str[2]) + string(style_str[4])
            } else {
                current = current + style_str
            }
            i += length
        }
        
        write_this.WriteString(current)
    }
    
    return write_this
}