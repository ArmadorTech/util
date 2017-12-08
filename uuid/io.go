/* Copyright (C) 2017 Armador Technologies
 * Author: Jose Luis Tallon <jltallon@armador.xyz>
 * * License: APL-2.0
 */

package uuid

import (
	"errors"
	"encoding/hex"
)


const (
	errPrefix = "uuid: "
	errInputShort	= "input string too short"
	errInputLong	= "input string too long" 
	urnPrefix = "urn:uuid:"
)


// String parse helpers.
var (
	byteGroups = []int{8, 4, 4, 4, 12}
)


// FromString returns UUID parsed from string input.
// Input is expected in a form accepted by UnmarshalText.
func FromString(input string) (u UUID, err error) {
	err = parseText(&u,[]byte(input))
	return
}



// MarshalText implements the encoding.TextMarshaler interface.
// The encoding is the same as returned by String.
func (u *UUID) MarshalText() (text []byte, err error) {
	return []byte(toString(u)),nil
}




// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (u *UUID) UnmarshalText(text []byte) (err error) {
	return parseText(u,text)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (u UUID) MarshalBinary() (data []byte, err error) {
	return u.Bytes(),nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// It will return error if the slice isn't 16 bytes long.
func (u *UUID) UnmarshalBinary(data []byte) (err error) {
	if len(data) != 16 {
		return errors.New(errPrefix+"invalid input length; Must be exactly 16 bytes")
	}
	copy(u[:], data)
	
	return
}


func toString(u *UUID) string {
	
	var buf [36]byte
	
	hex.Encode(buf[0:8], u[0:4])
	buf[8] = dash
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = dash
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = dash
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = dash
	hex.Encode(buf[24:], u[10:])
	
	return string(buf[:])
}


// The following formats are supported:
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
// "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
// "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
func parseText(u *UUID, text []byte) (err error) {
	if len(text) < 32 {
		return errors.New(errPrefix+errInputShort)
	}
	
	t := text[:]	// Copy into a new slice, same backing array
	braced := false
	if equalBytes(t[:9], []byte(urnPrefix), 9) {
		t = t[9:]
	} else if '{' == t[0] {
		braced = true
		t = t[1:]
	}
	
	if dash!=t[8] || dash!=t[13] || dash!=t[18] || dash!=t[23] {
		return errors.New(errPrefix+"invalid format")
	}
	
	b := u[:]
	for i, byteGroup := range byteGroups {
		if i > 0 {
			if dash != t[0] {
				return errors.New(errPrefix+"invalid string format")
			}
			t = t[1:]
		}
		
		if len(t) < byteGroup {
			return errors.New(errPrefix+errInputShort)
			return
		}
		
		if i == 4 && len(t) > byteGroup &&
			((braced && '}'!=t[byteGroup]) || len(t[byteGroup:]) > 1 || !braced) {
				err = errors.New(errPrefix+errInputLong)
				return
			}
			
			_, err = hex.Decode(b[:byteGroup/2], t[:byteGroup])
			if nil!=err {
				return
			}
			
			t = t[byteGroup:]
			b = b[byteGroup/2:]
	}
	return nil
}
