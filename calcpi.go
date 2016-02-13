/*
 * Compare to my Scala version using the actor model
 * https://prezi.com/4asyld78xkm6/actor-systems-with-akka/
 *
 * PI/4 = 1/1 - 1/3 + 1/5 - 1/7 + 1/9 - ...
 */
package main
import "fmt"
import "time"

func calcPiForTerms(start int, numElements int) chan float64 {

    out := make(chan float64)

    go func() {
        accum := 0.0

        for i := start; i < start + numElements; i++ {
            accum += 4.0 * float64(1 - (i % 2) * 2) / float64(2 * i + 1)
        }

        out <- accum
    }()

    return out
}

func calculatePi(numWorkers int, elementsPerWorker int) chan float64 {
    out := make(chan float64)

    go func() {
        var channels []chan float64

        for i := 0; i < numWorkers; i++ {
            c := calcPiForTerms(i * elementsPerWorker, elementsPerWorker)
            channels = append(channels, c)
        }
        out <- sum(channels)
    }()

    return out
}

func main() {

    totalElements := 10000000000
    numWorkers := 8
    elementsPerWorker := totalElements / numWorkers

    fmt.Printf("Calculating PI with %d workers and %d elements per worker\n", numWorkers, elementsPerWorker);

    start := time.Now()

    piChan := calculatePi(numWorkers, elementsPerWorker)

    fmt.Printf("Pi is %.10f\n", <-piChan) // Synchronization here. No need for a latch.
    fmt.Printf("Elapsed time: %s\n", time.Since(start))
}

////////////////// Generic utilities

func sum(channels []chan float64) float64 {
    sum := 0.0
    for _,c := range channels {
        sum += <-c
    }
    return sum
}
