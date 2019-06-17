package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	cleanup()
	os.Exit(code)
}

func setup() {
	if os.Getenv("MAIN_PORT") == "" {
		log.Fatalln("Main port is not set")
	}

	if os.Getenv("ENV") != "test" {
		err := os.Setenv("ENV", "test")

		if err != nil {
			log.Fatalln("ENV global variable is not set to 'test' and app has failed to set it itself")
		}
	}

	err := exec.Command("/bin/sh", "build.sh").Run()

	if err != nil {
		log.Fatalln("Build failed:", err)
	}

	err = exec.Command("/bin/sh", "start.sh").Run()

	if err != nil {
		log.Fatalln("Start failed:", err)
	}
}

func cleanup() {
	err := exec.Command("/bin/sh", "cleanup.sh").Run()

	if err != nil {
		log.Fatalln("Cleanup failed:", err)
	}
}

func TestGetTransaction(t *testing.T) {
	address := fmt.Sprintf("http://localhost:%s/api/v1/transaction/", os.Getenv("MAIN_PORT"))

	// Expect 404 if the transaction id was not provided
	resp, err := http.Get(address)

	if err != nil {
		t.Fatal(err)
	}

	if resp != nil {
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code %d", http.StatusNotFound)
		}
	}

	// Expect the data to be sent
	resp, err = http.Get(fmt.Sprintf("%s%s", address, "1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d"))

	if err != nil {
		t.Fatal(err)
	}

	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d", http.StatusOK)
		}
		if resp.Body == nil {
			t.Error("Expected body to not be nil")
		}
	} else {
		t.Error("Expected response to not be nil")
	}
}

func TestGetNewAddress(t *testing.T) {
	address := fmt.Sprintf("http://localhost:%s/api/v1/get-new-address", os.Getenv("MAIN_PORT"))

	// Expect 400 if the passphrase was not provided
	body := []byte(`{"Address_type": "legacy"}`)

	resp, err := http.Post(address, "application/json", bytes.NewBuffer(body))

	if err != nil {
		t.Fatal(err)
	}

	if resp != nil {
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d", http.StatusBadRequest)
		}

		if resp.Body == nil {
			t.Error("Expected body to not be nil")
		} else {
			respBody := struct {
				Message string
			}{}
			body, err = ioutil.ReadAll(resp.Body)

			err = json.Unmarshal(body, &respBody)

			if respBody.Message == "" {
				t.Error("Expected body to have some message")
			}

			_ = resp.Body.Close()
		}
	}

	// Expect 400 if the address type was not provided or had the wrong value
	body = []byte(`{"Passphrase":"PassphraseString","Address_type":"wrong-type"}`)

	resp, err = http.Post(address, "application/json", bytes.NewBuffer(body))

	if err != nil {
		t.Fatal(err)
	}

	if resp != nil {
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d", http.StatusBadRequest)
		}

		if resp.Body == nil {
			t.Error("Expected body to not be nil")
		} else {
			respBody := struct {
				Message string
			}{}
			body, err = ioutil.ReadAll(resp.Body)

			err = json.Unmarshal(body, &respBody)

			if respBody.Message == "" {
				t.Error("Expected body to have some message")
			}

			_ = resp.Body.Close()
		}
	}

	// Expect the data to be sent
	body = []byte(`{"Passphrase":"PassphraseString","Address_type":"legacy"}`)

	resp, err = http.Post(address, "application/json", bytes.NewBuffer(body))

	if err != nil {
		t.Fatal(err)
	}

	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d", http.StatusOK)
		}
		if resp.Body == nil {
			t.Error("Expected body to not be nil")
		} else {
			respBody := struct {
				Address string
			}{}
			body, err = ioutil.ReadAll(resp.Body)

			err = json.Unmarshal(body, &respBody)

			if respBody.Address == "" {
				t.Error("Expected body to have address")
			}

			_ = resp.Body.Close()
		}
	} else {
		t.Error("Expected response to not be nil")
	}
}

func TestSendToAddress(t *testing.T) {
	address := fmt.Sprintf("http://localhost:%s/api/v1/send-to-address", os.Getenv("MAIN_PORT"))

	// Expect 400 if the passphrase was not provided
	body := []byte(`{"Estimate_mode": "ECONOMICAL"}`)

	resp, err := http.Post(address, "application/json", bytes.NewBuffer(body))

	if err != nil {
		t.Fatal(err)
	}

	if resp != nil {
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d", http.StatusBadRequest)
		}

		if resp.Body == nil {
			t.Error("Expected body to not be nil")
		} else {
			respBody := struct {
				Message string
			}{}
			body, err = ioutil.ReadAll(resp.Body)

			err = json.Unmarshal(body, &respBody)

			if respBody.Message == "" {
				t.Error("Expected body to have some message")
			}

			_ = resp.Body.Close()
		}
	}

	// Expect 400 if the address type was not provided or had the wrong value
	body = []byte(`{"Passphrase":"PassphraseString","Estimate_mode":"wrong-type"}`)

	resp, err = http.Post(address, "application/json", bytes.NewBuffer(body))

	if err != nil {
		t.Fatal(err)
	}

	if resp != nil {
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d", http.StatusBadRequest)
		}

		if resp.Body == nil {
			t.Error("Expected body to not be nil")
		} else {
			respBody := struct {
				Message string
			}{}
			body, err = ioutil.ReadAll(resp.Body)

			err = json.Unmarshal(body, &respBody)

			if respBody.Message == "" {
				t.Error("Expected body to have some message")
			}

			_ = resp.Body.Close()
		}
	}

	// Expect the data to be sent
	body = []byte(`{"Passphrase":"PassphraseString","Estimate_mode":"ECONOMICAL"}`)

	resp, err = http.Post(address, "application/json", bytes.NewBuffer(body))

	if err != nil {
		t.Fatal(err)
	}

	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d", http.StatusOK)
		}
		if resp.Body == nil {
			t.Error("Expected body to not be nil")
		} else {
			respBody := struct {
				Txid string
			}{}
			body, err = ioutil.ReadAll(resp.Body)

			err = json.Unmarshal(body, &respBody)

			if respBody.Txid == "" {
				t.Error("Expected body to have transaction id")
			}

			_ = resp.Body.Close()
		}
	} else {
		t.Error("Expected response to not be nil")
	}
}
