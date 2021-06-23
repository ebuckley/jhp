package jhp

import (
	"fmt"
	"testing"
)

func TestOutEcho(t *testing.T) {
	ex := toEcho(`
<h1>Hello world</h1>
<p>This should be normal output</p>
<?jhp
const things = [1,2,3,4];
?>
<p>It should run each jhp block independently</p>

<?jhp things.forEach(function (tx) { ?>	
Everything in between the jhp should turn in to an echo command <?jhp echo(tx); ?>
<?jhp }); ?>
<p>Hmm... this is actually going to be a bit tricky isn't it...</p>`)
	// if we rewrite a jhp to be a page, this is waht it would look like...
	expectedOut := `echo("\n\u003ch1\u003eHello world\u003c/h1\u003e\n\u003cp\u003eThis should be normal output\u003c/p\u003e\n");

const things = [1,2,3,4];
echo("\n\u003cp\u003eIt should run each jhp block independently\u003c/p\u003e\n\n");
 things.forEach(function (tx) { echo("\t\nEverything in between the jhp should turn in to an echo command ");
 echo(tx); echo("\n");
 }); echo("\n\u003cp\u003eHmm... this is actually going to be a bit tricky isn't it...\u003c/p\u003e");`
	if ex != expectedOut {
		t.Log("expected:")
		fmt.Print(expectedOut)
		fmt.Print("\n")
		t.Log("----- got -----")
		fmt.Print(ex)
		fmt.Print("\n")
		t.Log("---------------")
		t.Fail()
	}
}

func TestName(t *testing.T) {
	ex := toEcho(`<a href="<?jhp echo(params.offset); ?>"></a>`)
	if ex != `` {
		fmt.Println(ex)
		t.Fail()
	}
}

