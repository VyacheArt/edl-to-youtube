package main

import "github.com/VyacheArt/edl-to-youtube/converter"

//func main() {
//	edlContent, err := os.ReadFile("/Users/vyacheart/Downloads/Timeline 1.edl")
//	if err != nil {
//		panic(err)
//	}
//
//	edlList, err := edl.Parse(string(edlContent))
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Printf("%+v", edlList)
//}

func main() {
	a := converter.NewApplication()
	if err := a.Run(); err != nil {
		panic(err)
	}
}
