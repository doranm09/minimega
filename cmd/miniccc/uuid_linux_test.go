// Copyright 2015-2023 National Technology & Engineering Solutions of Sandia, LLC (NTESS).
// Under the terms of Contract DE-NA0003525 with NTESS, the U.S. Government retains certain
// rights in this software.

//go:build linux
// +build linux

package main

import (
	"errors"
	"strings"
	"testing"
)

func TestReadLinuxUUIDPrefersDMI(t *testing.T) {
	uuid, err := readLinuxUUID(func(path string) ([]byte, error) {
		if path == linuxUUIDPaths[0] {
			return []byte("A5BA6920-5BCF-4022-B8CF-015425F7B05C\n"), nil
		}
		return nil, errors.New("unexpected fallback")
	})
	if err != nil {
		t.Fatal(err)
	}
	if uuid != "a5ba6920-5bcf-4022-b8cf-015425f7b05c" {
		t.Fatalf("unexpected UUID: %q", uuid)
	}
}

func TestReadLinuxUUIDFallsBackToQEMUFWCfg(t *testing.T) {
	uuid, err := readLinuxUUID(func(path string) ([]byte, error) {
		if path == linuxUUIDPaths[0] {
			return nil, errors.New("DMI unavailable")
		}
		return []byte("54FCF635-B72E-4D4D-9A7A-21830E61935B\x00"), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if uuid != "54fcf635-b72e-4d4d-9a7a-21830e61935b" {
		t.Fatalf("unexpected UUID: %q", uuid)
	}
}

func TestReadLinuxUUIDReportsAllSources(t *testing.T) {
	_, err := readLinuxUUID(func(path string) ([]byte, error) {
		return nil, errors.New("not present")
	})
	if err == nil {
		t.Fatal("expected an error")
	}
	for _, path := range linuxUUIDPaths {
		if !strings.Contains(err.Error(), path) {
			t.Fatalf("missing path %q in error: %v", path, err)
		}
	}
}
