# markov.go

An extremely simple Markov generator based on [Jordan Scales's Ruby-based one](https://github.com/jdan/markov.rb).

## Usage

```sh
$ go build -o ./markov
$ ./markov --help
Usage of ./markov:
  -c int
    	Number of results to generate (default 10)
  -n int
    	Length of n-gram (default 2)
  -w	Split text into words
$ ./markov -n 3 < names.txt
new Farayah
new	Tamalandre
exists	Donn
new	Eristinique
new	Eliusephona
exists	Lavar
exists	Arlena
new	Beja
new	Johnnyse
new	Jedad
$ ./markov -n 1 -w -c 1 < war-n-peace.txt
new	Why can’t make me time for his saber which some raisins, fine day, gentlemen.
```