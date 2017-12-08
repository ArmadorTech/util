/** simple-uuid
 * Copyright (C) 2017 Armador Technologies
 * Author: Jose Luis Tallon <jltallon@armador.xyz>
 * License: APL-2.0
 *
 * Loosely based upon code Copyright (C) 2013-2015 
 * by Maxim Bublis <b@codemonkey.ru>, released under the MIT License
 */

// (Simple) UUID provides a simple and fast implementation of UUID primitives
package uuid


import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/md5"
	"hash"
)


type UUID [16]byte

// The nil UUID is special form of UUID that is specified to have all
// 128 bits set to zero.
var Nil = UUID{}




// UUID layout variants.
const (
	VariantNCS 		= iota
	VariantRFC4122
	VariantMicrosoft
	VariantFuture
)

// Used in string method conversion
const dash byte = '-'


// Predefined namespace UUIDs.
var (
	NamespaceDNS, _  = FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	NamespaceURL, _  = FromString("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	NamespaceOID, _  = FromString("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	NamespaceX500, _ = FromString("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
)



// Version returns algorithm version used to generate UUID.
func (u UUID) Version() uint {
	return uint(u[6] >> 4)
}

// Variant returns UUID layout variant.
func (u UUID) Variant() uint {
	switch {
		case (u[8] & 0x80) == 0x00:
			return VariantNCS
		case (u[8]&0xc0)|0x80 == 0x80:
			return VariantRFC4122
		case (u[8]&0xe0)|0xc0 == 0xc0:
			return VariantMicrosoft
	}
	return VariantFuture
}

// Bytes returns bytes slice representation of UUID.
func (u UUID) Bytes() []byte {
	return u[:]
}

// Returns canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u UUID) String() string {
	return toString(&u)
}

// URN returns the RFC 2141 URN form of uuid,
// urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx,  or "" if uuid is invalid.
func (u *UUID) URN() string {
	var buf [36 + 9]byte
	copy(buf[0:], urnPrefix)
	copy(buf[9:], toString(u))
	return string(buf[:])
}


// SetVersion sets version bits.
func (u *UUID) SetVersion(v byte) {
	u[6] = (u[6] & 0x0f) | (v << 4)
}

// SetVariant sets variant bits as described in RFC 4122.
func (u *UUID) SetVariant() {
	u[8] = (u[8] & 0xbf) | 0x80
}


// NewV3 returns UUID based on MD5 hash of namespace UUID and name.
func NewV3(ns UUID, name string) UUID {
	u := newFromHash(md5.New(), ns, name)
	u.SetVersion(3)
// 	u.SetVariant()
	u[8] = (u[8] & 0x3f) | 0x80 // Variant is 10
	
	return u
}


// The quality of the UUIDs depends directly on crypto/rand
// From Wikipedia (on UUIDs):
// "Randomly generated UUIDs have 122 random bits.  One's annual risk of being
//  hit by a meteorite is estimated to be one chance in 17 billion, that
//  means the probability is about 0.00000000006 (6 × 10−11),
//  equivalent to the odds of creating a few tens of trillions of UUIDs in a
//  year and having one duplicate."


// NewV4 returns random generated UUID.
func NewV4() UUID {
	var u UUID
	rand.Read(u[:])
 	u.SetVersion(4)
// 	u.SetVariant()
	u[8] = (u[8] & 0x3f) | 0x80 // Variant is 10
	
	return u
}


// NewRandom returns a Random (Version 4) UUID.
func NewRandom() UUID {
	var u UUID
	rand.Read(u[:])
	u.SetVersion(4)
	u[8] = (u[8] & 0x3f) | 0x80 // Variant is 10
	return u
}


// NewV5 returns UUID based on SHA-1 hash of namespace UUID and name.
func NewV5(ns UUID, name string) UUID {
	u := newFromHash(sha1.New(), ns, name)
	u.SetVersion(5)
	u.SetVariant()
	return u
}

// Returns UUID based on hashing of namespace UUID and name.
func newFromHash(h hash.Hash, ns UUID, name string) UUID {
	u := UUID{}
	h.Write(ns[:])
	h.Write([]byte(name))
	copy(u[:], h.Sum(nil))
		
	return u
}
