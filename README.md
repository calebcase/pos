# Proof of Space

POS aims to provide a [Proof of Space/Proof of Capacity][pos] system with the
following properties:

* Prover must have (fast) random access to the space
* Nominal exchange of data between the challenger and prover
* Challenger can validate the proof without the space
* Quality of the proof should be proportional to the amount of space claimed

## Algorithm

- P: Make a claim of *b* bytes of space to C.
- C: Select a random seed *s* and send to P.
- P: Use the seed *s* to initialize a [PRNG][prng] and generate an output file
  *o* of *b* bytes.
- P: Tell C that it is done generating the output file.
- C: Start a timer *t*.
- C: Select a random mask *m* of *n* bytes and send to P.
- C: Select a desired number of hash inputs *i* and send to P.
- P: Select *n* bytes of output file *o* and [XOR][xor] with mask *m*
  to produce a new seed *s2*.
- P: Use seed *s2* to initialize a PRNG. Generate *i* blocks of *n* bytes.
  Modulo each block by *b* bytes to produce a list of indices *x*.
- P: Compute a hash *h* by reading the *x* th byte from *o* in series and send
  to C.
- C: Stop the timer *t*. If *t* is larger than the allowed time *At*, then reject
  the proof.
- C: Following a similar procedure as P, compute the hash *h*. This can be
  done without the required space by computing the output stream *o* twice. On
  the first pass only the *n* bytes are retained to be used to compute *s2*. On
  the second pass retain the values of the indices needed for the hash *h*.

The allowed time *At* must be low enough that P cannot perform the two-pass
method within the interval timed by *t*.

---

[pos]: https://en.wikipedia.org/wiki/Proof-of-space
[prng]: https://en.wikipedia.org/wiki/Pseudorandom_number_generator
[xor]: https://en.wikipedia.org/wiki/Exclusive_or
