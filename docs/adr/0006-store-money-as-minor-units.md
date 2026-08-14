# Store Money as Minor Units

Money amounts will be stored as integer minor units, such as satang for THB, rather than floating-point numbers. Split rounding will be deterministic so manual and equal splits can be recalculated without introducing invisible precision errors.
