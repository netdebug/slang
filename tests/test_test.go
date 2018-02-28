package tests

import (
	"io/ioutil"
	"github.com/Bitspark/slang/tests/assertions"
	"testing"
	"github.com/Bitspark/slang/pkg/slapi"
)

func TestTestOperator__TrivialTests(t *testing.T) {
	a := assertions.New(t)
	succs, fails, err := slapi.TestOperator("test_data/voidOp_test.json", ioutil.Discard, true)
	a.Nil(err)
	a.Equal(1, succs)
	a.Equal(0, fails)
}

func TestTestOperator__SimpleFail(t *testing.T) {
	a := assertions.New(t)
	succs, fails, err := slapi.TestOperator("test_data/voidOp_corruptTest.json", ioutil.Discard, true)
	a.Nil(err)
	a.Equal(0, succs)
	a.Equal(1, fails)
}

func TestTestOperator__ComplexTest(t *testing.T) {
	a := assertions.New(t)
	succs, fails, err := slapi.TestOperator("test_data/nested_op/usingSubCustomOpDouble_test.json", ioutil.Discard, true)
	a.Nil(err)
	a.Equal(2, succs)
	a.Equal(0, fails)
}

func TestTestOperator__SuiteTests(t *testing.T) {
	a := assertions.New(t)

	succs, fails, err := slapi.TestOperator("test_data/suite/polynomial_test.yaml", ioutil.Discard, false)
	a.Nil(err)
	a.Equal(1, succs)
	a.Equal(0, fails)

	succs, fails, err = slapi.TestOperator("test_data/suite/main_test.yaml", ioutil.Discard, false)
	a.Nil(err)
	a.Equal(2, succs)
	a.Equal(0, fails)
}

func TestOperator_Pack(t *testing.T) {
	a := assertions.New(t)
	succs, fails, err := slapi.TestOperator("test_data/slib/pack_test.yaml", ioutil.Discard, false)
	a.Nil(err)
	a.Equal(1, succs)
	a.Equal(0, fails)
}

func TestTestOperator__SumReduce(t *testing.T) {
	a := assertions.New(t)
	succs, fails, err := slapi.TestOperator("test_data/sum/reduce_test.yaml", ioutil.Discard, true)
	a.Nil(err)
	a.Equal(4, succs)
	a.Equal(0, fails)
}

func TestTestOperator__MergeSort(t *testing.T) {
	a := assertions.New(t)
	succs, fails, err := slapi.TestOperator("test_data/slib/merge_sort_test.yaml", ioutil.Discard, true)
	a.Nil(err)
	a.Equal(5, succs)
	a.Equal(0, fails)
}

func TestTestOperator_Properties(t *testing.T) {
	a := assertions.New(t)
	succs, fails, err := slapi.TestOperator("test_data/properties/prop_op_test.yaml", ioutil.Discard, true)
	a.Nil(err)
	a.Equal(3, succs)
	a.Equal(0, fails)
}

func TestTestOperator_NestedProperties(t *testing.T) {
	a := assertions.New(t)
	succs, fails, err := slapi.TestOperator("test_data/properties/prop2_op_test.yaml", ioutil.Discard, true)
	a.Nil(err)
	a.Equal(3, succs)
	a.Equal(0, fails)
}

func TestTestOperator_NestedDelegates(t *testing.T) {
	a := assertions.New(t)
	succs, fails, err := slapi.TestOperator("test_data/delegates/wrapper_test.yaml", ioutil.Discard, true)
	a.Nil(err)
	a.Equal(3, succs)
	a.Equal(0, fails)
}
