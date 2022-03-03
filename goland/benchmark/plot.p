# Generate with `gnuplot plot.p`

set title "[ns/op] for Dela to execute a smart contract that performs N times a crypto operation"

set terminal pdf
set output 'plot.pdf'

set logscale y
set logscale x

set key left top

set ylabel "[ns/op]"
set xlabel "Number of operations"
#set xrange [-1:110]

set style data linespoints
plot "result.dat" using 1:2 title 'Native' linewidth 2 , \
     "result.dat" using 1:3 title 'Unikernel' linewidth 2, \
     "result.dat" using 1:4 title 'EVM' linewidth 2, \
     "result.dat" using 1:5 title 'WASM' linewidth 2, \
     "result.dat" using 1:6 title 'GraalVM' linewidth 2

#