package abb_b21

import (
	"abb-exporter/internal/meter"
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/simonvetter/modbus"
)

var _ meter.Meter = (*Client)(nil)

func NewClient(mc *modbus.ModbusClient, name string, slave uint8) *Client {
	return &Client{
		modbus: mc,
		name:   name,
		slave:  slave,
		tags: map[string]string{
			"name":  name,
			"slave": fmt.Sprintf("%d", slave),
		},
	}
}

type Client struct {
	slave  uint8
	modbus *modbus.ModbusClient
	name   string
	tags   map[string]string
}

func (c Client) Tags() *map[string]string {
	return &c.tags
}

func (c Client) QueryInfo() (*meter.Info, error) {
	_ = c.modbus.SetUnitId(c.slave)
	serial, err := c.modbus.ReadRegisters(35072, 2, modbus.HOLDING_REGISTER)
	if err != nil {
		return nil, err
	}

	versions, err := c.modbus.ReadRegisters(35080, 9, modbus.HOLDING_REGISTER)
	if err != nil {
		return nil, err
	}
	mmv := versions[8:][0]

	des, err := c.modbus.ReadRegisters(35168, 6, modbus.HOLDING_REGISTER)
	if err != nil {
		return nil, err
	}

	info := meter.Info{
		Serial:          convertToUint16(serial),
		FirmwareVersion: readString(convertToBinary(versions[0:8])),
		TypeDesignation: readString(convertToBinary(des)),
		ModbusMappingVersion: meter.Firmware{
			Major: mmv >> 8,
			Minor: mmv & 0xFF,
		},
	}

	return &info, nil
}

func (c Client) QueryUsageStatus() (*meter.UsageStatus, error) {
	_ = c.modbus.SetUnitId(c.slave)
	v, err := c.modbus.ReadRegisters(0x5B00, 2, modbus.HOLDING_REGISTER)
	if err != nil {
		return nil, err
	}

	current, err := c.modbus.ReadRegisters(0x5B0C, 2, modbus.HOLDING_REGISTER)
	activePower, err := c.modbus.ReadRegisters(0x5B14, 2, modbus.HOLDING_REGISTER)
	frequency, err := c.modbus.ReadRegisters(0x5B2C, 1, modbus.HOLDING_REGISTER)

	voltageRounded := math.Floor(float64(convertToUint16(v[0:2]))*0.1*100) / 100
	currentRounded := math.Floor(float64(convertToUint16(current))*0.01*100) / 100
	powerRounded := math.Floor(float64(convertToUint16(activePower))*0.01*100) / 100
	frequencyRounded := math.Floor(float64(frequency[0])*0.01*100) / 100

	reg := meter.UsageStatus{
		Timestamp:   time.Now(),
		Voltage:     voltageRounded,
		Current:     currentRounded,
		ActivePower: powerRounded,
		Frequency:   frequencyRounded,
	}

	return &reg, nil
}

func (c Client) QueryTotalActiveImport() (float64, error) {
	_ = c.modbus.SetUnitId(c.slave)
	v, err := c.modbus.ReadRegisters(0x5000, 4, modbus.HOLDING_REGISTER)
	if err != nil {
		return 0, err
	}

	b := convertToBinary(v)
	return float64(binary.BigEndian.Uint64(b)) * 0.01, nil
}

func (c Client) ReadFrequency() (float32, error) {
	_ = c.modbus.SetUnitId(c.slave)
	v, err := c.modbus.ReadRegisters(0x5B2C, 1, modbus.HOLDING_REGISTER)
	if err != nil {
		return 0, err
	}

	return float32(v[0]) / 100, nil
}

func (c Client) ActivePower() (float32, error) {
	_ = c.modbus.SetUnitId(c.slave)
	v, err := c.modbus.ReadRegisters(0x5B14, 2, modbus.HOLDING_REGISTER)
	if err != nil {
		return 0, err
	}

	watt := convertToUint16(v)
	return float32(watt) / 100, nil
}

func (c Client) Current() (float32, error) {
	_ = c.modbus.SetUnitId(c.slave)
	v, err := c.modbus.ReadRegisters(23314, 2, modbus.HOLDING_REGISTER)
	if err != nil {
		return 0, err
	}

	watt := convertToUint16(v)
	return float32(watt) / 100, nil
}

func (c Client) Name() string {
	return c.name
}

func convertToUint16(regs []uint16) uint32 {
	b := make([]byte, len(regs)*2)
	for i := 0; i < len(regs); i++ {
		b[i*2] = byte(regs[i] >> 8)
		b[i*2+1] = byte(regs[i] & 0xFF)
	}

	return binary.BigEndian.Uint32(b)
}

func convertToBinary(v []uint16) []byte {
	b := make([]byte, len(v)*2)
	for i := 0; i < len(v); i++ {
		b[i*2] = byte(v[i] >> 8)
		b[i*2+1] = byte(v[i] & 0xFF)
	}

	return b
}

func readString(v []byte) string {
	safe := make([]byte, 0)
	for _, b := range v {
		if b != 255 && b != 0 {
			safe = append(safe, b)
		}
	}

	return string(safe)
}
