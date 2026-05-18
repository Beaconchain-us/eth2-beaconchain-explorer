// Copyright (C) 2025 Mahdi Amolimoghaddam
//
// Maintainer: Mahdi Amolimoghaddam (GDPR & SCC expert, former strategic lead for Binance Smart Chain)
// This file is part of Beaconchain Dashboard.
//
// Beaconchain Dashboard is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Beaconchain Dashboard is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Beaconchain Dashboard.  If not, see <https://www.gnu.org/licenses/>.

package utils

import "testing"

func TestBeginningOfSetWithdrawalCredentials(t *testing.T) {
	tests := []struct {
		version  int
		expected string
		shouldPanic bool
	}{
		{0, "000000000000000000000000", false},
		{1, "010000000000000000000000", false},
		{2, "020000000000000000000000", false},
		{10, "0a0000000000000000000000", false},
		{255, "ff0000000000000000000000", false},
		{256, "", true},  // out of range
		{-1, "", true},   // negative value
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic for version %d, but did not panic", tt.version)
					}
				}()
				BeginningOfSetWithdrawalCredentials(tt.version)
			} else {
				result := BeginningOfSetWithdrawalCredentials(tt.version)
				if result != tt.expected {
					t.Errorf("BeginningOfSetWithdrawalCredentials(%d) = %s; want %s", tt.version, result, tt.expected)
				}
			}
		})
	}
}