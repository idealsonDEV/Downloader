// main
package main

//"os"
//"fmt"

func main() {
	//url := os.Args[1]
	//down := NewDownloader("http://www92.uptobox.com/dl/KWe9ZM_wwKPc2j5DE_p3GGJgx0XAsky2zHz1UhxWci5c-IoaTn5mycKxj5bXprS5aMGiFgU9l4s9_YHXJO16ZrUYQNgSFH0hF1DgEuZEvhTgo_hGS1ZzprVg0tTISCrkXKjCw2Ulc6gmF4cZnmCqKg/AAX2MP3.zip", 0)
	//down.Do()
	//down.Join()
	down := Resume("AAX2MP3.zip", "http://www92.uptobox.com/dl/KWe9ZM_wwKPc2j5DE_p3GGJgx0XAsky2zHz1UhxWci5c-IoaTn5mycKxj5bXprS5aMGiFgU9l4s9_YHXJO16ZrUYQNgSFH0hF1DgEuZEvhTgo_hGS1ZzprVg0tTISCrkXKjCw2Ulc6gmF4cZnmCqKg/AAX2MP3.zip")
	//down.DoResume()
	down.Join()
}
