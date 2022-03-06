package treadonme_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/swedishborgie/treadonme"
)

type MessageTestSuite struct {
	suite.Suite
}

func (s *MessageTestSuite) TestReadGetDeviceInfo() {
	msg, err := treadonme.ParseMessage(fromHex("5b08f092000178050f125d"))
	s.Require().NoError(err)
	s.Require().IsType(&treadonme.MessageDeviceInfo{}, msg)

	info := msg.(*treadonme.MessageDeviceInfo)

	s.Require().Equal(treadonme.DeviceModelF80, info.Model)
	s.Require().Equal(byte(0), info.Version)
	s.Require().Equal(treadonme.UnitsTypeImperial, info.Units)
	s.Require().Equal(treadonme.Speed(120), info.MaxSpeed)
	s.Require().Equal(treadonme.Speed(5), info.MinSpeed)
	s.Require().Equal(byte(15), info.InclineMax)
	s.Require().Equal(byte(18), info.UserSegment)

	s.T().Log(msg)
}

func (s *MessageTestSuite) TestWriteGetDeviceInfo() {
	msgBytes, err := treadonme.EncodeMessage(&treadonme.MessageDeviceInfo{})
	s.Require().NoError(err)
	s.Require().Equal(fromHex("5B01F05D"), msgBytes)
}

func (s *MessageTestSuite) TestWorkoutMode() {
	msg, err := treadonme.ParseMessage(fromHex("5b0203015d"))
	s.Require().NoError(err)
	s.Require().IsType(&treadonme.MessageWorkoutMode{}, msg)

	mode := msg.(*treadonme.MessageWorkoutMode)
	s.Require().Equal(treadonme.WorkoutModeIdle, mode.Mode)

	encoded, err := treadonme.EncodeMessage(msg)
	s.Require().NoError(err)
	s.Require().Equal(fromHex("5B0203015D"), encoded)
}

func (s *MessageTestSuite) TestSetWorkoutMode() {
	msg, err := treadonme.ParseMessage(fromHex("5B0202025D"))
	s.Require().NoError(err)
	s.Require().IsType(&treadonme.MessageSetWorkoutMode{}, msg)

	mode := msg.(*treadonme.MessageSetWorkoutMode)
	s.Require().Equal(treadonme.WorkoutModeStart, mode.Mode)

	encoded, err := treadonme.EncodeMessage(msg)
	s.Require().NoError(err)
	s.Require().Equal(fromHex("5B0202025D"), encoded)
}

func (s *MessageTestSuite) TestHeartRateType() {
	msg, err := treadonme.ParseMessage(fromHex("5b030901005d"))
	s.Require().NoError(err)
	s.Require().IsType(&treadonme.MessageHeartRateType{}, msg)

	hrType := msg.(*treadonme.MessageHeartRateType)
	s.Require().Equal(byte(1), hrType.Type1)
	s.Require().Equal(byte(0), hrType.Type2)

	encoded, err := treadonme.EncodeMessage(msg)
	s.Require().NoError(err)
	s.Require().Equal(fromHex("5b030901005d"), encoded)
}

func (s *MessageTestSuite) TestACK() {
	msg, err := treadonme.ParseMessage(fromHex("5B0400094F4B5D"))
	s.Require().NoError(err)
	s.Require().IsType(&treadonme.MessageACK{}, msg)

	ack := msg.(*treadonme.MessageACK)
	s.Require().Equal(treadonme.MessageTypeHeartRateType, ack.Acknowledged)

	encoded, err := treadonme.EncodeMessage(msg)
	s.Require().NoError(err)
	s.Require().Equal(fromHex("5B0400094F4B5D"), encoded)
}

func (s *MessageTestSuite) TestUserProfile() {
	msg, err := treadonme.ParseMessage(fromHex("5B06070123009B435D"))
	s.Require().NoError(err)
	s.Require().IsType(&treadonme.MessageUserProfile{}, msg)

	profile := msg.(*treadonme.MessageUserProfile)
	s.Require().Equal(treadonme.SexTypeMale, profile.Sex)
	s.Require().Equal(byte(35), profile.Age)
	s.Require().Equal(treadonme.Weight(155), profile.Weight)
	s.Require().Equal(treadonme.Height(67), profile.Height)

	encoded, err := treadonme.EncodeMessage(msg)
	s.Require().NoError(err)
	s.Require().Equal(fromHex("5B06070123009B435D"), encoded)
}

func (s *MessageTestSuite) TestWorkoutTarget() {
	msg, err := treadonme.ParseMessage(fromHex("5B05040A0000005D"))
	s.Require().NoError(err)
	s.Require().IsType(&treadonme.MessageWorkoutTarget{}, msg)

	target := msg.(*treadonme.MessageWorkoutTarget)
	s.Require().Equal(byte(10), target.Time)
	s.Require().Equal(uint16(0), target.Calories)

	encoded, err := treadonme.EncodeMessage(msg)
	s.Require().NoError(err)
	s.Require().Equal(fromHex("5B05040A0000005D"), encoded)
}

func (s *MessageTestSuite) TestMaxIncline() {
	msg, err := treadonme.ParseMessage(fromHex("5B0222095D"))
	s.Require().NoError(err)
	s.Require().IsType(&treadonme.MessageMaxIncline{}, msg)

	max := msg.(*treadonme.MessageMaxIncline)
	s.Require().Equal(byte(9), max.MaxIncline)

	encoded, err := treadonme.EncodeMessage(msg)
	s.Require().NoError(err)
	s.Require().Equal(fromHex("5B0222095D"), encoded)
}

func (s *MessageTestSuite) TestProgramGraphics() {
	msg, err := treadonme.ParseMessage(fromHex("5b13400101010101010101010101010101010101015d"))
	s.Require().NoError(err)
	s.Require().IsType(&treadonme.MessageProgramGraphics{}, msg)

	graph := msg.(*treadonme.MessageProgramGraphics)
	s.Require().Equal(byte(1), graph.Graph[0])

	encoded, err := treadonme.EncodeMessage(msg)
	s.Require().NoError(err)
	s.Require().Equal(fromHex("5b13400101010101010101010101010101010101015d"), encoded)
}

func (s *MessageTestSuite) TestWorkoutData() {
	msg, err := treadonme.ParseMessage(fromHex("5b0f06093b0000000000050000000000015d"))
	s.Require().NoError(err)
	s.Require().IsType(&treadonme.MessageWorkoutData{}, msg)

	wd := msg.(*treadonme.MessageWorkoutData)
	s.Require().Equal(byte(9), wd.Minute)
	s.Require().Equal(byte(0x3b), wd.Second)
	s.Require().Equal(uint16(0), wd.Distance)
	s.Require().Equal(uint16(0), wd.Calories)
	s.Require().Equal(byte(0), wd.HeartRate)
	s.Require().Equal(treadonme.Speed(5), wd.Speed)
	s.Require().Equal(byte(0), wd.Incline)
	s.Require().Equal(byte(0), wd.HRType)
	s.Require().Equal(byte(0), wd.IntervalTime)
	s.Require().Equal(byte(0), wd.RecoveryTime)
	s.Require().Equal(byte(0), wd.ProgramRow)
	s.Require().Equal(byte(1), wd.ProgramColumn)

	encoded, err := treadonme.EncodeMessage(msg)
	s.Require().NoError(err)
	s.Require().Equal(fromHex("5b0f06093b0000000000050000000000015d"), encoded)
}

func (s *MessageTestSuite) TestMessageCommand() {
	msg, err := treadonme.ParseMessage(fromHex("5B02F1025D"))
	s.Require().NoError(err)
	s.Require().IsType(&treadonme.MessageCommand{}, msg)

	wd := msg.(*treadonme.MessageCommand)
	s.Require().Equal(treadonme.CommandTypeLevelUp, wd.Command)

	encoded, err := treadonme.EncodeMessage(msg)
	s.Require().NoError(err)
	s.Require().Equal(fromHex("5B02F1025D"), encoded)
}

func (s *MessageTestSuite) TestProgram() {
	msg, err := treadonme.ParseMessage(fromHex("5b030810015d"))
	s.Require().NoError(err)
	s.Require().IsType(&treadonme.MessageProgram{}, msg)

	prog := msg.(*treadonme.MessageProgram)
	s.Require().Equal(treadonme.ProgramManual, prog.Program)

	encoded, err := treadonme.EncodeMessage(msg)
	s.Require().NoError(err)
	s.Require().Equal(fromHex("5b030810015d"), encoded)
}

func TestMessageTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, &MessageTestSuite{})
}

func fromHex(hexStr string) []byte {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}

	return data
}
