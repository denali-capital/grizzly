type Exchange interface {
    // functionality we want each exchange to implement
    // * getting data
    // * executing orders
    // * checking status of orders
}

// each exchange will have their own module that implements Exchange interface above
// can have module to compute statistics in background and access them
// Neural network module that handles predict / fit functionality

func main() {
    // need to create array of exchange objects
    // pass this array into function that computes statistics in background

    // run algo
}