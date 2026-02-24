package data_processor

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type mockReadCloser struct {
	mock.Mock
}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *mockReadCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestParser(t *testing.T) {
	type testCase struct {
		inp []byte
		out Input
		err error
	}

	testCases := []testCase{
		{
			inp: []byte("test\n+\n1\n2"),
			out: Input{
				Id:   "test",
				Op:   "+",
				Val1: 1,
				Val2: 2,
			},
			err: nil,
		},
		{
			inp: []byte("test\n+\nlalala\n2"),
			out: Input{
				Id:   "",
				Op:   "",
				Val1: 0,
				Val2: 0,
			},
			err: strconv.ErrSyntax,
		},
		{
			inp: []byte("test\n+\n1\nlalala"),
			out: Input{
				Id:   "",
				Op:   "",
				Val1: 0,
				Val2: 0,
			},
			err: strconv.ErrSyntax,
		},
		{
			inp: []byte(""),
			out: Input{
				Id:   "",
				Op:   "",
				Val1: 0,
				Val2: 0,
			},
			err: errInvalidFormat,
		},
	}

	for _, tc := range testCases {
		output, err := parser(tc.inp)
		if output != tc.out {
			t.Errorf("Expected %v; got %v", tc.out, output)
		} else if !errors.Is(err, tc.err) {
			t.Errorf("Expected error to be %v; got %v", tc.err, err)
		}
	}
}

func TestDataProcessor(t *testing.T) {
	type testCase struct {
		paramIn  chan []byte
		paramOut chan Result
		expVal   Result
	}

	testCases := []testCase{
		{
			paramIn:  make(chan []byte, 100),
			paramOut: make(chan Result, 100),
			expVal: Result{
				Id:    "test-1",
				Value: 1,
			},
		},
		{
			paramIn:  make(chan []byte, 100),
			paramOut: make(chan Result, 100),
			expVal: Result{
				Id:    "test-2",
				Value: 2,
			},
		},
		{
			paramIn:  make(chan []byte, 100),
			paramOut: make(chan Result, 100),
			expVal: Result{
				Id:    "test-3",
				Value: 3,
			},
		},
		{
			paramIn:  make(chan []byte, 100),
			paramOut: make(chan Result, 100),
			expVal: Result{
				Id:    "test-4",
				Value: 4,
			},
		},
		{
			paramIn:  make(chan []byte, 100),
			paramOut: make(chan Result, 100),
			expVal: Result{
				Id:    "test-5/Error - Unknown Operand",
				Value: 0,
			},
		},
		{
			paramIn:  make(chan []byte, 100),
			paramOut: make(chan Result, 100),
			expVal: Result{
				Id:    "Unknown/Error - Invalid Format",
				Value: 0,
			},
		},
		{
			paramIn:  make(chan []byte, 100),
			paramOut: make(chan Result, 100),
			expVal: Result{
				Id:    "test-7/Error - Division by zero",
				Value: 0,
			},
		},
	}

	testCases[0].paramIn <- []byte("test-1\n+\n1\n0")
	testCases[1].paramIn <- []byte("test-2\n-\n3\n1")
	testCases[2].paramIn <- []byte("test-3\n*\n3\n1")
	testCases[3].paramIn <- []byte("test-4\n/\n8\n2")
	testCases[4].paramIn <- []byte("test-5\nrandom\n100\n200")
	testCases[5].paramIn <- []byte("two\nlines")
	testCases[6].paramIn <- []byte("test-7\n/\n0\n0")
	close(testCases[0].paramIn)
	close(testCases[1].paramIn)
	close(testCases[2].paramIn)
	close(testCases[3].paramIn)
	close(testCases[4].paramIn)
	close(testCases[5].paramIn)
	close(testCases[6].paramIn)

	for _, tc := range testCases {
		DataProcessor(tc.paramIn, tc.paramOut)
		select {
		case val := <-tc.paramOut:
			if val != tc.expVal {
				t.Errorf("Expected %v; got %v", tc.expVal, val)
			}
		}
	}
}

func TestWriteData(t *testing.T) {
	type testCase struct {
		paramIn chan Result
		paramW  *bytes.Buffer
		expVal  []byte
	}

	wOne := bytes.NewBuffer(make([]byte, 0, 100))
	wTwo := bytes.NewBuffer(make([]byte, 0, 100))
	testCases := []testCase{
		{
			paramIn: make(chan Result, 1),
			paramW:  wOne,
			expVal:  []byte("testCase:4\n"),
		},
		{
			paramIn: make(chan Result, 1),
			paramW:  wTwo,
			expVal:  []byte("testCase/Error - Unknown Operand:0\n"),
		},
	}

	testCases[0].paramIn <- Result{Id: "testCase", Value: 4}
	testCases[1].paramIn <- Result{Id: "testCase/Error - Unknown Operand", Value: 0}

	close(testCases[0].paramIn)
	close(testCases[1].paramIn)

	for _, tc := range testCases {
		WriteData(tc.paramIn, tc.paramW)
		val := make([]byte, len(tc.expVal))
		_, err := tc.paramW.Read(val)
		if err != nil && !errors.Is(err, io.EOF) {
			t.Errorf("Unexpected error: %v", err)
		}

		stringVal, stringExpVal := string(val), string(tc.expVal)
		if stringVal != stringExpVal {
			t.Errorf("Expected %s; got %s", stringExpVal, stringVal)
		}
	}
}

func TestNewController(t *testing.T) {
	type testCase struct {
		paramW        *httptest.ResponseRecorder
		paramR        *http.Request
		expVal        []byte
		expStatusCode int
	}

	stubChan := make(chan []byte, 100)
	testHandler := NewController(stubChan)

	mockBuffer := &mockReadCloser{}
	mockBuffer.On("Read", mock.AnythingOfType("[]uint8")).Return(0, fmt.Errorf("Error Reading"))
	mockBuffer.On("Close").Return(fmt.Errorf("Error Closing"))

	testCases := []testCase{
		{
			paramW:        httptest.NewRecorder(),
			paramR:        httptest.NewRequest("POST", "/", nil),
			expVal:        []byte("Bad Input"),
			expStatusCode: http.StatusBadRequest,
		},
		{
			paramW:        httptest.NewRecorder(),
			paramR:        httptest.NewRequest("POST", "/", bytes.NewBufferString("test-1\n-\n1\n1")),
			expVal:        []byte("OK: "),
			expStatusCode: http.StatusAccepted,
		},
		{
			paramW:        httptest.NewRecorder(),
			paramR:        httptest.NewRequest("POST", "/", mockBuffer),
			expVal:        []byte("Bad Input"),
			expStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		testHandler.ServeHTTP(tc.paramW, tc.paramR)
		val := tc.paramW.Result()
		defer val.Body.Close()
		valData, err := io.ReadAll(val.Body)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !bytes.Contains(valData, tc.expVal) {
			t.Errorf("Expected '%s' to contain '%s'", string(valData), string(tc.expVal))
		}

		if val.StatusCode != tc.expStatusCode {
			t.Errorf("Expected Status Code %v; got %v", tc.expStatusCode, val.StatusCode)
		}
	}
	// Simulating full inp channel
	newStubChan := make(chan []byte, 1)
	newTestHandler := NewController(newStubChan)
	newStubChan <- []byte("some value")

	lastTestCase := testCase{
		paramW:        httptest.NewRecorder(),
		paramR:        httptest.NewRequest("POST", "/", bytes.NewBufferString("test-3\n+\n2\n1")),
		expVal:        []byte("Too Busy: "),
		expStatusCode: http.StatusServiceUnavailable,
	}

	newTestHandler.ServeHTTP(lastTestCase.paramW, lastTestCase.paramR)
	val := lastTestCase.paramW.Result()
	defer val.Body.Close()
	valData, err := io.ReadAll(val.Body)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !bytes.Contains(valData, lastTestCase.expVal) {
		t.Errorf("Expected '%s' to contain '%s'", string(valData), string(lastTestCase.expVal))
	}

	if val.StatusCode != lastTestCase.expStatusCode {
		t.Errorf("Expected Status Code %v; got %v", lastTestCase.expStatusCode, val.StatusCode)
	}
}
