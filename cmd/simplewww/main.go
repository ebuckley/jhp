package main

import (
	"github.com/dop251/goja"
	"jhp"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		vm := goja.New()
		err := jhp.Register(vm, jhp.RequestReader{
			DB:     nil,
			Params: nil,
		}, w)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_, _ = vm.RunString(" echo(\"<h1>HELLO jhp</h1>\" + Date.now()) ");
	})

	http.HandleFunc("/sample.jhp", func(w http.ResponseWriter, r *http.Request) {
		vm := goja.New()
		req := jhp.RequestReader{
			DB:     nil,
			Params: nil,
		}
		err := jhp.Render(vm, req, w, `
<h1>Hello world</h1>
<p>This should be normal output</p>
<?jhp
const things = [1,2,3,4];
?>
<p>It should run each jhp block independently</p>

<?jhp things.forEach(function (tx) { ?>	
Everything in between the jhp should turn in to an echo command <?jhp echo(tx); ?>
<?jhp }); ?>
<p>Hmm... this is actually going to be a bit tricky isn't it...</p>
`)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
	log.Println("Running on :9293")
	err := http.ListenAndServe(":9293", nil)
	if err != nil {
		panic(err)
	}
}

