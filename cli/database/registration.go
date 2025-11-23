package database

import (
	"bytes"
	"database/sql"
	"encoding/binary"

	"github.com/mattn/go-sqlite3"
	"gonum.org/v1/gonum/floats"
)

func RegisterCustomDriver() {
	sql.Register("sqlite_extended",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(sc *sqlite3.SQLiteConn) error {
				return sc.RegisterFunc("vector_distance", vectorDistance, true)
			},
		})
}

func vectorDistance(a, b []byte) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 999.0
	}

	fa, err := bytesToFloatSlice(a)
	if err != nil {
		return 999.0
	}

	fb, err := bytesToFloatSlice(b)
	if err != nil {
		return 999.0
	}

	normA := floats.Norm(fa, 2)
	normB := floats.Norm(fb, 2)

	if normA == 0.0 || normB == 0.0 {
		return 999.0
	}

	distance := 1.0 - floats.Dot(fa, fb)/(normA*normB)
	return distance
}

func bytesToFloatSlice(b []byte) ([]float64, error) {
	floatCount := len(b) / 8
	floats := make([]float64, floatCount)
	err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &floats)

	if err != nil {
		logger.Error("Binary conversion failed", "error", err)
		return nil, err
	}
	return floats, nil
}

func FloatSliceToBytes(floats []float64) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, floats)
	if err != nil {
		logger.Error("Binary write failed during conversion", "error", err)
		return nil, err
	}
	return buf.Bytes(), nil
}
