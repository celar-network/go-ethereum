package vm

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

type createdCodeRecorder struct {
	NoopOpCodeHooks
	code []byte
	err  error
}

func (r *createdCodeRecorder) CreatedCodeHook(_ *EVM, _ common.Address, code []byte) error {
	r.code = code
	return r.err
}

// The OPTIONALITY is the design property worth pinning, because it is the
// whole reason this is a separate interface rather than a third method on
// OpCodeHooks: existing implementations must keep compiling untouched.
//
// If someone later "tidies" the extension into OpCodeHooks, this test fails
// and names the consequence — downstream hooks types that assert against
// OpCodeHooks would have to be forked too.
func TestCreatedCodeHookIsOptional(t *testing.T) {
	var base OpCodeHooks = NoopOpCodeHooks{}
	if _, ok := base.(CreatedCodeHooks); ok {
		t.Fatal("NoopOpCodeHooks satisfies CreatedCodeHooks: the extension is " +
			"no longer optional, and every existing OpCodeHooks implementation " +
			"now has to provide it")
	}

	var extended OpCodeHooks = &createdCodeRecorder{}
	if _, ok := extended.(CreatedCodeHooks); !ok {
		t.Fatal("an implementation that provides CreatedCodeHook is not detected " +
			"by the assertion the EVM uses, so the hook would never be consulted")
	}
}
