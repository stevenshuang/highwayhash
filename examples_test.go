// Copyright (c) 2018 Minio Inc. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package highwayhash

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// ExampleNew shows how to use HighwayHash-256 to compute fingerprints of files.
func ExampleNew() {
	key, err := hex.DecodeString("000102030405060708090A0B0C0D0E0FF0E0D0C0B0A090807060504030201000") // use your own key here
	if err != nil {
		fmt.Printf("Cannot decode hex key: %v", err) // add error handling
		return
	}

	file, err := os.Open("./LICENSE") // specify your file here
	if err != nil {
		fmt.Printf("Failed to open the file: %v", err) // add error handling
		return
	}
	defer file.Close()

	hash, err := New(key)
	if err != nil {
		fmt.Printf("Failed to create HighwayHash instance: %v", err) // add error handling
		return
	}

	if _, err = io.Copy(hash, file); err != nil {
		fmt.Printf("Failed to read from file: %v", err) // add error handling
		return
	}

	checksum := hash.Sum(nil)
	fmt.Println(hex.EncodeToString(checksum))

	// Output: f4217352f920fb8d287c4948eb0d843d8c4d43b61a7c0d658f28036bd1d270b7
}

// ExampleNew64 shows how to use HighwayHash-64 to implement a content-addressable storage.
func ExampleNew64() {
	key, err := hex.DecodeString("000102030405060708090A0B0C0D0E0FF0E0D0C0B0A090807060504030201000") // use your own key here
	if err != nil {
		fmt.Printf("Cannot decode hex key: %v", err) // add error handling
		return
	}

	AddressOf := func(key []byte, file string) (uint64, error) { // function to compute address based on content
		fsocket, err := os.Open(file)
		if err != nil {
			return 0, err
		}
		defer fsocket.Close()

		hash, err := New64(key)
		if err != nil {
			return 0, err
		}

		_, err = io.Copy(hash, fsocket)
		return hash.Sum64(), err
	}

	dir, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Printf("Failed to read current directory: %v", err) // add error handling
		return
	}

	lookupMap := make(map[uint64]string, len(dir))
	for _, file := range dir {
		if file.IsDir() {
			continue // skip sub-directroies in our example
		}
		address, err := AddressOf(key, file.Name())
		if err != nil {
			fmt.Printf("Failed to read file %s: %v", file.Name(), err) // add error handling
			return
		}
		lookupMap[address] = file.Name()
	}
	// Output:
}
