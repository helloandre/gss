GSS
====

GSS is a css compiler written in [Go](http://golang.org). It's not incredibly efficient time-wise, but in my testing has produced a smaller output file that YUICompressor.

Why?
---

Because I wanted to play with Go and was interested in parsing CSS. GSS does two passes over the input file tokenizing then compressing. This makes it really inefficent, but it was fun to build.

Eventuall I would like to do something similar with Javascript in Go as well.

Usage
---

You can either use the `run` script included or compile it once and run it with just that executable. The usage below can be used for either approach.

`./run -o <output_file> <input_file>`

Testing
---

I've included Twitter's Bootstrap CSS as a benchmark. My testing shows GSS creates a smaller file that functions the same. To run this test:

`./run -o test/input/bootstrap.css test/output/bootstrap.css`

Please let me know or file an issue if you see something amiss.