package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestBearerToken(t *testing.T) {

  testCases := []struct{
    input string
    want string
  } {
    { input: "Bearer test", want: "test", },
    { input: "test", want: "", },
  }

  for _, tt := range testCases {
    header := http.Header{}
    header.Add("Authorization", tt.input)

    output, _ := GetBearerToken(&header)

    if output != tt.want {
      t.Errorf("Processing failed wanted: %v actual: %v", tt.want, output)
      t.Fail()
    }
    
  }
}

func TestJWTToken(t *testing.T) {

  testCases := []struct{
    input string
    want string
  } {
    { input: "TECHNO-USER", want: "TECHNO-USER", },
  }

  for _, tt := range testCases {
    input, _ := uuid.Parse(tt.input)

    jwt, err := MakeJWT(input, "TOP-SECRET", 24 * time.Hour)
    if err != nil {
      t.Logf("Unable to create JWT: %v", err)
      t.Fail()
    }

    jwtUuid, err := ValidateJWT(jwt, "TOP-SECRET")
    if err != nil {
      t.Logf("Unable to validate JWT: %v", err)
      t.Fail()
    }

    expected, _ := uuid.Parse(tt.want)
    if jwtUuid != expected {
      t.Logf("uuids are not the same, expected: %v got: %v", expected, jwtUuid)
      t.Fail()
    }
  }
}
