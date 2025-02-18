package main

import "testing"

func BenchmarkRender(b *testing.B) {
	//for i := 0; i < 1; i++ {
	res, _ := render("{{#each tags}}{{#each tags}}{{#each tags}}{{this}}{{/each}}{{/each}}{{/each}}", map[string]any{
		"tags": numbers,
	})
	println(len(res))
	//}
}
