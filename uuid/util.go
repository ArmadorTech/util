/* Copyright (C) 2017 Armador Technologies
 * Author: Jose Luis Tallon <jltallon@armador.xyz>
 * * License: APL-2.0
*/

package uuid




// And returns result of binary AND of two UUIDs.
func And(u1 UUID, u2 UUID) UUID {
	u := UUID{}
	for i := 0; i < 16; i++ {
		u[i] = u1[i] & u2[i]
	}
	return u
}

// Or returns result of binary OR of two UUIDs.
func Or(u1 UUID, u2 UUID) UUID {
	u := UUID{}
	for i := 0; i < 16; i++ {
		u[i] = u1[i] | u2[i]
	}
	return u
}

// Equal returns true if u1 and u2 equals, otherwise returns false.
func Equal(u1 UUID, u2 UUID) bool {
	return equalBytes(u1[:],u2[:],16)
}


func equalBytes(lhs []byte, rhs []byte, n int) bool {
	var a byte = 0
	for i:=0; i<n; i++ {
		a |= lhs[i] ^ rhs[i]
	}
	return (0==a)
}