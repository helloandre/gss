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

package main

import (
    "bytes"
    "flag"
    "bufio"
    "os"
    "./compress/css"
    "./tokenize/css"
    "strings"
)

const (
    unknown = "SOME_RIDICULOUS_STRING"
)

var (
    output = flag.String("o", unknown, "output file")
    file_type = flag.String("type", unknown, "css")
    input_file string
    r *bufio.Reader
    w *bufio.Writer
)

func main() {
    // get all the args for the file
    flag.Parse()
    args := os.Args
    
    // set up reader
    if last_arg := args[len(args)-1]; last_arg != *output && last_arg != *file_type {
        fi, err := os.Open(last_arg)
        if err != nil { panic(err) }
        defer fi.Close()
        r = bufio.NewReader(fi)
        
        input_file = last_arg
    } else {
        r = bufio.NewReader(os.Stdin)
        input_file = unknown
    }
    
    // set up writer
    if *output != unknown {
        fo, err := os.Create(*output)
        if err != nil { panic(err) }
        defer fo.Close()
        w = bufio.NewWriter(fo)
    } else {
        w = bufio.NewWriter(os.Stdout)
    }
    
    if *file_type == unknown && input_file == unknown {
        panic("cannot determine file type")
    } else if *file_type == unknown {
        splitter := strings.Split(input_file, ".")
        *file_type = splitter[len(splitter)-1]
    }
    
    var output bytes.Buffer
    if *file_type == "css" {
        tokens := tokenize_css.Parse(r)
        output = compress_css.Compress(tokens)
    } else {
        panic("filetype not supported")
    }
    
    if _, err := w.Write(output.Bytes()); err != nil { panic(err) }
    if err := w.Flush(); err != nil { panic(err) }
}

