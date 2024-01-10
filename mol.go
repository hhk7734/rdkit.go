package rdkit

/*
#cgo CFLAGS: -I${SRCDIR}/third_party/rdkit
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/third_party/rdkit -lrdkitcffi_linux_amd64 -lm -lstdc++ -lfreetype

#include <stdlib.h>
#include "cffiwrapper.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

var (
	ErrDeleted = errors.New("rdkit: Mol is deleted")
)

type Mol struct {
	pkl       *C.char
	pklSize   C.size_t
	isDeleted bool
}

// NewMol creates a new Mol from input string. The input string must be a RDKit
// JSON, MolBlock(Molfile or SDF) or SMILES.
func NewMol(input string) Mol {
	smile := C.CString(input)
	defer C.free(unsafe.Pointer(smile))

	nullChar := C.CString("")
	defer C.free(unsafe.Pointer(nullChar))

	var pklSize C.size_t
	pkl := C.get_mol(smile, (*C.ulong)(unsafe.Pointer(&pklSize)), nullChar)

	return Mol{
		pkl:       pkl,
		pklSize:   pklSize,
		isDeleted: false,
	}
}

func (m *Mol) MolBlock() ([]byte, error) {
	if m.isDeleted {
		return nil, ErrDeleted
	}

	nullChar := C.CString("")
	defer C.free(unsafe.Pointer(nullChar))

	cmol := C.get_molblock(m.pkl, m.pklSize, nullChar)
	defer C.free(unsafe.Pointer(cmol))

	mol := C.GoString(cmol)

	return []byte(mol), nil
}

func (m *Mol) MolBlockV3000() ([]byte, error) {
	if m.isDeleted {
		return nil, ErrDeleted
	}

	nullChar := C.CString("")
	defer C.free(unsafe.Pointer(nullChar))

	cmol := C.get_v3kmolblock(m.pkl, m.pklSize, nullChar)
	defer C.free(unsafe.Pointer(cmol))

	mol := C.GoString(cmol)

	return []byte(mol), nil
}

func (m *Mol) SMILES() (string, error) {
	if m.isDeleted {
		return "", ErrDeleted
	}

	nullChar := C.CString("")
	defer C.free(unsafe.Pointer(nullChar))

	csmile := C.get_smiles(m.pkl, m.pklSize, nullChar)
	defer C.free(unsafe.Pointer(csmile))

	smile := C.GoString(csmile)

	return smile, nil
}

func (m Mol) MarshalJSON() ([]byte, error) {
	if m.isDeleted {
		return nil, ErrDeleted
	}

	nullChar := C.CString("")
	defer C.free(unsafe.Pointer(nullChar))

	cjson := C.get_json(m.pkl, m.pklSize, nullChar)
	defer C.free(unsafe.Pointer(cjson))

	json := C.GoString(cjson)

	return []byte(json), nil
}

func (m *Mol) UnmarshalJSON(data []byte) error {
	mol := NewMol(string(data))
	C.free(unsafe.Pointer(m.pkl))

	m.pkl = mol.pkl
	m.pklSize = mol.pklSize

	return nil
}

func (m *Mol) SVG() ([]byte, error) {
	if m.isDeleted {
		return nil, ErrDeleted
	}

	nullChar := C.CString("")
	defer C.free(unsafe.Pointer(nullChar))

	csvg := C.get_svg(m.pkl, m.pklSize, nullChar)
	defer C.free(unsafe.Pointer(csvg))

	svg := C.GoString(csvg)

	return []byte(svg), nil
}

func (m *Mol) Delete() {
	if !m.isDeleted {
		m.isDeleted = true
		C.free(unsafe.Pointer(m.pkl))
	}
}
