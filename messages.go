package treadonme

import (
	"encoding/hex"
	"fmt"
)

type Message interface {
	MessageType() MessageType
	MarshalBinary() ([]byte, error)
	UnmarshalBinary([]byte) error
	ExpectedLength() int
	String() string
}

var ErrInvalidMessage = fmt.Errorf("invalid message")

func ParseMessage(data []byte) (Message, error) {
	if err := sanityCheckRecv(data); err != nil {
		return nil, err
	}

	// Unwrap data (trim length/start/end)
	data = data[2 : len(data)-1]

	msgType := MessageType(data[0])
	msg := msgType.Create()

	if msg == nil {
		return nil, fmt.Errorf("%w: unknown message id: %d: %s", ErrInvalidMessage, int(msgType), hex.EncodeToString(data))
	} else if err := msg.UnmarshalBinary(data); err != nil {
		return nil, err
	}

	return msg, nil
}

func EncodeMessage(msg Message) ([]byte, error) {
	msgBytes, err := msg.MarshalBinary()
	if err != nil {
		return nil, err
	}

	encoded := make([]byte, 0, len(msgBytes)+3)

	encoded = append(encoded, startOfMessage, byte(len(msgBytes)))
	encoded = append(encoded, msgBytes...)
	encoded = append(encoded, endOfMessage)

	return encoded, nil
}

func sanityCheckRecv(data []byte) error {
	if len(data) <= 3 {
		return fmt.Errorf("%w: expected >3 bytes, got: %d: %s", ErrInvalidMessage, len(data), hex.EncodeToString(data))
	} else if data[0] != startOfMessage {
		return fmt.Errorf("%w: expected message to start with %d, got %d: %s", ErrInvalidMessage, startOfMessage, data[0], hex.EncodeToString(data))
	} else if data[len(data)-1] != endOfMessage {
		return fmt.Errorf("%w: expected message to end with %d, got %d: %s", ErrInvalidMessage, endOfMessage, data[len(data)-1], hex.EncodeToString(data))
	} else if data[1] != byte(len(data)-3) {
		return fmt.Errorf("%w: expected length was %d, got %d: %s", ErrInvalidMessage, data[1], len(data)-3, hex.EncodeToString(data))
	}

	return nil
}

func messageSanityCheck(msg Message, data []byte) error {
	if len(data) != msg.ExpectedLength() {
		return fmt.Errorf("%w: expected %d bytes for %s, got: %d",
			ErrInvalidMessage,
			msg.ExpectedLength(),
			msg.MessageType(),
			len(data),
		)
	} else if MessageType(data[0]) != msg.MessageType() {
		return fmt.Errorf("%w: expected command to be %s but was %s",
			ErrInvalidMessage, msg.MessageType(), MessageType(data[0]))
	}

	return nil
}

type MessageDeviceInfo struct {
	Model       DeviceModel
	Version     byte
	Units       UnitsType
	MaxSpeed    Speed
	MinSpeed    Speed
	InclineMax  byte
	UserSegment byte
}

func (di *MessageDeviceInfo) MessageType() MessageType {
	return MessageTypeDeviceInfo
}

func (di *MessageDeviceInfo) MarshalBinary() ([]byte, error) {
	return []byte{
		byte(di.MessageType()),
	}, nil
}

func (di *MessageDeviceInfo) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(di, data); err != nil {
		return err
	}

	di.Model = DeviceModel(data[1])
	di.Version = data[2]
	di.Units = UnitsType(data[3])
	di.MaxSpeed = Speed(data[4])
	di.MinSpeed = Speed(data[5])
	di.InclineMax = data[6]
	di.UserSegment = data[7]

	return nil
}

func (di *MessageDeviceInfo) ExpectedLength() int {
	return 8
}

func (di *MessageDeviceInfo) String() string {
	return fmt.Sprintf("DeviceInfo[Model=%s,Version=%d,Units=%s,MaxSpeed=%d,MinSpeed=%d,InclineMax=%d,UserSegment=%d]",
		di.Model,
		di.Version,
		di.Units,
		di.MaxSpeed,
		di.MinSpeed,
		di.InclineMax,
		di.UserSegment,
	)
}

type MessageWorkoutMode struct {
	Mode WorkoutMode
}

func (wm *MessageWorkoutMode) MessageType() MessageType {
	return MessageTypeWorkoutMode
}

func (wm *MessageWorkoutMode) MarshalBinary() ([]byte, error) {
	return []byte{byte(wm.MessageType()), byte(wm.Mode)}, nil
}

func (wm *MessageWorkoutMode) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(wm, data); err != nil {
		return err
	}

	wm.Mode = WorkoutMode(data[1])

	return nil
}

func (wm *MessageWorkoutMode) ExpectedLength() int {
	return 2
}

func (wm *MessageWorkoutMode) String() string {
	return fmt.Sprintf("WorkoutMode[Mode=%s]", wm.Mode)
}

type MessageSetWorkoutMode struct {
	Mode WorkoutMode
}

func (wm *MessageSetWorkoutMode) MessageType() MessageType {
	return MessageTypeSetWorkoutMode
}

func (wm *MessageSetWorkoutMode) MarshalBinary() ([]byte, error) {
	return []byte{byte(wm.MessageType()), byte(wm.Mode)}, nil
}

func (wm *MessageSetWorkoutMode) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(wm, data); err != nil {
		return err
	}

	wm.Mode = WorkoutMode(data[1])

	return nil
}

func (wm *MessageSetWorkoutMode) ExpectedLength() int {
	return 2
}

func (wm *MessageSetWorkoutMode) String() string {
	return fmt.Sprintf("SetWorkoutMode[Mode=%s]", wm.Mode)
}

type MessageHeartRateType struct {
	Type1 byte
	Type2 byte
}

func (hr *MessageHeartRateType) MessageType() MessageType {
	return MessageTypeHeartRateType
}

func (hr *MessageHeartRateType) MarshalBinary() ([]byte, error) {
	return []byte{byte(hr.MessageType()), hr.Type1, hr.Type2}, nil
}

func (hr *MessageHeartRateType) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(hr, data); err != nil {
		return err
	}

	hr.Type1 = data[1]
	hr.Type2 = data[2]
	return nil
}

func (hr *MessageHeartRateType) ExpectedLength() int {
	return 3
}

func (hr *MessageHeartRateType) String() string {
	return fmt.Sprintf("HeartRateType[Type1=%d,Type2=%d]", hr.Type1, hr.Type2)
}

type MessageWorkoutData struct {
	Minute        byte
	Second        byte
	Distance      uint16
	Calories      uint16
	HeartRate     byte
	Speed         Speed
	Incline       byte
	HRType        byte
	IntervalTime  byte
	RecoveryTime  byte
	ProgramRow    byte
	ProgramColumn byte
}

func (wd *MessageWorkoutData) MessageType() MessageType {
	return MessageTypeWorkoutData
}

func (wd *MessageWorkoutData) MarshalBinary() ([]byte, error) {
	return []byte{
		byte(wd.MessageType()),
		wd.Minute,
		wd.Second,
		byte(wd.Distance >> 8),
		byte(wd.Distance),
		byte(wd.Calories >> 8),
		byte(wd.Calories),
		wd.HeartRate,
		byte(wd.Speed),
		wd.Incline,
		wd.HRType,
		wd.IntervalTime,
		wd.RecoveryTime,
		wd.ProgramRow,
		wd.ProgramColumn,
	}, nil
}

func (wd *MessageWorkoutData) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(wd, data); err != nil {
		return err
	}

	wd.Minute = data[1]
	wd.Second = data[2]
	wd.Distance = (uint16(data[3]) << 8) + uint16(data[4])
	wd.Calories = (uint16(data[5]) << 8) + uint16(data[6])
	wd.HeartRate = data[7]
	wd.Speed = Speed(data[8])
	wd.Incline = data[9]
	wd.HRType = data[10]
	wd.IntervalTime = data[11]
	wd.RecoveryTime = data[12]
	wd.ProgramRow = data[13]
	wd.ProgramColumn = data[14]

	return nil
}

func (wd *MessageWorkoutData) ExpectedLength() int {
	return 15
}

func (wd *MessageWorkoutData) String() string {
	return fmt.Sprintf("WorkoutData[Minute=%d,Second=%d,Distance=%d,Calorites=%d,HeartRate=%d,"+
		"Speed=%d,Incline=%d,HRType=%d,IntervalTime=%d,RecoveryTime=%d,ProgramRow=%d,ProgramColumn=%d]",
		wd.Minute,
		wd.Second,
		wd.Distance,
		wd.Calories,
		wd.HeartRate,
		wd.Speed,
		wd.Incline,
		wd.HRType,
		wd.IntervalTime,
		wd.RecoveryTime,
		wd.ProgramRow,
		wd.ProgramColumn,
	)
}

type MessageUserProfile struct {
	Sex    SexType
	Age    byte
	Weight Weight
	Height Height
}

func (p *MessageUserProfile) MessageType() MessageType {
	return MessageTypeUserProfile
}

func (p *MessageUserProfile) MarshalBinary() ([]byte, error) {
	return []byte{
		byte(p.MessageType()),
		byte(p.Sex),
		p.Age,
		byte(p.Weight >> 8),
		byte(p.Weight),
		byte(p.Height),
	}, nil
}

func (p *MessageUserProfile) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(p, data); err != nil {
		return err
	}

	p.Sex = SexType(data[1])
	p.Age = data[2]
	p.Weight = Weight((uint16(data[3]) << 8) + uint16(data[4]))
	p.Height = Height(data[5])

	return nil
}

func (p *MessageUserProfile) ExpectedLength() int {
	return 6
}

func (p *MessageUserProfile) String() string {
	return fmt.Sprintf("UserProfile[Sex=%s,Age=%d,Weight=%d,Height=%d]",
		p.Sex, p.Age, p.Weight, p.Height,
	)
}

type MessageACK struct {
	Acknowledged MessageType
}

func (ack *MessageACK) MessageType() MessageType {
	return MessageTypeACK
}

func (ack *MessageACK) MarshalBinary() ([]byte, error) {
	return []byte{
		byte(ack.MessageType()),
		byte(ack.Acknowledged),
		0x4f, 0x4b,
	}, nil
}

func (ack *MessageACK) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(ack, data); err != nil {
		return err
	}

	if data[2] != 0x4f && data[3] != 0x4b {
		return fmt.Errorf(
			"%w: expected ack to end with 0x4F4B, but was: 0x%s",
			ErrInvalidReadPayload, hex.EncodeToString(data[2:]),
		)
	}

	ack.Acknowledged = MessageType(data[1])

	return nil
}

func (ack *MessageACK) ExpectedLength() int {
	return 4
}

func (ack *MessageACK) String() string {
	return fmt.Sprintf("ACK[Acknowledged=%s]", ack.Acknowledged)
}

type MessageWorkoutTarget struct {
	Time     byte
	Calories uint16
}

func (wt *MessageWorkoutTarget) MessageType() MessageType {
	return MessageTypeWorkoutTarget
}

func (wt *MessageWorkoutTarget) MarshalBinary() ([]byte, error) {
	return []byte{
		byte(wt.MessageType()),
		wt.Time,
		0,
		byte(wt.Calories >> 8),
		byte(wt.Calories),
	}, nil
}

func (wt *MessageWorkoutTarget) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(wt, data); err != nil {
		return err
	}

	wt.Time = data[1]

	wt.Calories = uint16(data[3])<<8 + uint16(data[4])

	return nil
}

func (wt *MessageWorkoutTarget) ExpectedLength() int {
	return 5
}

func (wt *MessageWorkoutTarget) String() string {
	return fmt.Sprintf("WorkoutTarget[Time=%d,Calories=%d]", wt.Time, wt.Calories)
}

type MessageMaxIncline struct {
	MaxIncline byte
}

func (mi *MessageMaxIncline) MessageType() MessageType {
	return MessageTypeMaxIncline
}

func (mi *MessageMaxIncline) MarshalBinary() ([]byte, error) {
	return []byte{byte(mi.MessageType()), mi.MaxIncline}, nil
}

func (mi *MessageMaxIncline) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(mi, data); err != nil {
		return err
	}

	mi.MaxIncline = data[1]

	return nil
}

func (mi *MessageMaxIncline) ExpectedLength() int {
	return 2
}

func (mi *MessageMaxIncline) String() string {
	return fmt.Sprintf("MaxIncline[MaxIncline=%d]", mi.MaxIncline)
}

type MessageErrorCode struct {
	Code byte
}

func (e *MessageErrorCode) MessageType() MessageType {
	return MessageTypeErrorCode
}

func (e *MessageErrorCode) MarshalBinary() ([]byte, error) {
	return []byte{byte(e.MessageType()), e.Code}, nil
}

func (e *MessageErrorCode) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.Code = data[1]

	return nil
}

func (e *MessageErrorCode) ExpectedLength() int {
	return 2
}

func (e *MessageErrorCode) String() string {
	return fmt.Sprintf("ErrorCode[Code=%d]", e.Code)
}

type MessageSpeed struct {
	Speed Speed
}

func (e *MessageSpeed) MessageType() MessageType {
	return MessageTypeSpeed
}

func (e *MessageSpeed) MarshalBinary() ([]byte, error) {
	return []byte{byte(e.MessageType()), byte(e.Speed)}, nil
}

func (e *MessageSpeed) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.Speed = Speed(data[1])

	return nil
}

func (e *MessageSpeed) ExpectedLength() int {
	return 2
}

func (e *MessageSpeed) String() string {
	return fmt.Sprintf("Speed[Speed=%d]", e.Speed)
}

type MessageIncline struct {
	Incline byte
}

func (e *MessageIncline) MessageType() MessageType {
	return MessageTypeIncline
}

func (e *MessageIncline) MarshalBinary() ([]byte, error) {
	return []byte{byte(e.MessageType()), e.Incline}, nil
}

func (e *MessageIncline) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.Incline = data[1]

	return nil
}

func (e *MessageIncline) ExpectedLength() int {
	return 2
}

func (e *MessageIncline) String() string {
	return fmt.Sprintf("Incline[Incline=%d]", e.Incline)
}

type MessageLevel struct {
	Level byte
}

func (e *MessageLevel) MessageType() MessageType {
	return MessageTypeLevel
}

func (e *MessageLevel) MarshalBinary() ([]byte, error) {
	return []byte{byte(e.MessageType()), e.Level}, nil
}

func (e *MessageLevel) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.Level = data[1]

	return nil
}

func (e *MessageLevel) ExpectedLength() int {
	return 2
}

func (e *MessageLevel) String() string {
	return fmt.Sprintf("Level[Level=%d]", e.Level)
}

type MessageRPM struct {
	RPM byte
}

func (e *MessageRPM) MessageType() MessageType {
	return MessageTypeRPM
}

func (e *MessageRPM) MarshalBinary() ([]byte, error) {
	return []byte{byte(e.MessageType()), e.RPM}, nil
}

func (e *MessageRPM) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.RPM = data[1]

	return nil
}

func (e *MessageRPM) ExpectedLength() int {
	return 2
}

func (e *MessageRPM) String() string {
	return fmt.Sprintf("RPM[RPM=%d]", e.RPM)
}

type MessageHeartRate struct {
	HeartRate byte
}

func (e *MessageHeartRate) MessageType() MessageType {
	return MessageTypeHeartRate
}

func (e *MessageHeartRate) MarshalBinary() ([]byte, error) {
	return []byte{byte(e.MessageType()), e.HeartRate}, nil
}

func (e *MessageHeartRate) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.HeartRate = data[1]

	return nil
}

func (e *MessageHeartRate) ExpectedLength() int {
	return 2
}

func (e *MessageHeartRate) String() string {
	return fmt.Sprintf("HeartRate[HeartRate=%d]", e.HeartRate)
}

type MessageTargetHeartRate struct {
	HeartRate byte
}

func (e *MessageTargetHeartRate) MessageType() MessageType {
	return MessageTypeTargetHeartRate
}

func (e *MessageTargetHeartRate) MarshalBinary() ([]byte, error) {
	return []byte{byte(e.MessageType()), e.HeartRate}, nil
}

func (e *MessageTargetHeartRate) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.HeartRate = data[1]

	return nil
}

func (e *MessageTargetHeartRate) ExpectedLength() int {
	return 2
}

func (e *MessageTargetHeartRate) String() string {
	return fmt.Sprintf("TargetHeartRate[HeartRate=%d]", e.HeartRate)
}

type MessageMaxSpeed struct {
	Speed Speed
}

func (e *MessageMaxSpeed) MessageType() MessageType {
	return MessageTypeMaxSpeed
}

func (e *MessageMaxSpeed) MarshalBinary() ([]byte, error) {
	return []byte{byte(e.MessageType()), byte(e.Speed)}, nil
}

func (e *MessageMaxSpeed) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.Speed = Speed(data[1])

	return nil
}

func (e *MessageMaxSpeed) ExpectedLength() int {
	return 2
}

func (e *MessageMaxSpeed) String() string {
	return fmt.Sprintf("MaxSpeed[Speed=%d]", e.Speed)
}

type MessageMaxLevel struct {
	Level byte
}

func (e *MessageMaxLevel) MessageType() MessageType {
	return MessageTypeMaxLevel
}

func (e *MessageMaxLevel) MarshalBinary() ([]byte, error) {
	return []byte{byte(e.MessageType()), e.Level}, nil
}

func (e *MessageMaxLevel) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.Level = data[1]

	return nil
}

func (e *MessageMaxLevel) ExpectedLength() int {
	return 2
}

func (e *MessageMaxLevel) String() string {
	return fmt.Sprintf("MaxLevel[Level=%d]", e.Level)
}

type MessageUserIncline struct {
	Incline byte
}

func (e *MessageUserIncline) MessageType() MessageType {
	return MessageTypeUserIncline
}

func (e *MessageUserIncline) MarshalBinary() ([]byte, error) {
	return []byte{byte(e.MessageType()), e.Incline}, nil
}

func (e *MessageUserIncline) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.Incline = data[1]

	return nil
}

func (e *MessageUserIncline) ExpectedLength() int {
	return 2
}

func (e *MessageUserIncline) String() string {
	return fmt.Sprintf("UserIncline[Incline=%d]", e.Incline)
}

type MessageUserLevel struct {
	Level byte
}

func (e *MessageUserLevel) MessageType() MessageType {
	return MessageTypeUserLevel
}

func (e *MessageUserLevel) MarshalBinary() ([]byte, error) {
	return []byte{byte(e.MessageType()), e.Level}, nil
}

func (e *MessageUserLevel) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.Level = data[1]

	return nil
}

func (e *MessageUserLevel) ExpectedLength() int {
	return 2
}

func (e *MessageUserLevel) String() string {
	return fmt.Sprintf("UserLevel[Level=%d]", e.Level)
}

type MessageEndWorkout struct {
	Seconds   uint16
	Distance  uint16
	Calories  uint16
	Speed     Speed
	HeartRate byte
	Incline   byte
}

func (e *MessageEndWorkout) MessageType() MessageType {
	return MessageTypeEndWorkout
}

func (e *MessageEndWorkout) MarshalBinary() ([]byte, error) {
	return []byte{
		byte(e.MessageType()),
		byte(e.Seconds >> 8), byte(e.Seconds),
		byte(e.Distance >> 8), byte(e.Distance),
		byte(e.Calories >> 8), byte(e.Calories),
		byte(e.Speed),
		e.HeartRate,
		e.Incline,
	}, nil
}

func (e *MessageEndWorkout) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.Seconds = (uint16(data[1]) << 8) + uint16(data[2])
	e.Distance = (uint16(data[3]) << 8) + uint16(data[4])
	e.Calories = (uint16(data[5]) << 8) + uint16(data[6])
	e.Speed = Speed(data[7])
	e.HeartRate = data[8]
	e.Incline = data[9]

	return nil
}

func (e *MessageEndWorkout) ExpectedLength() int {
	return 10
}

func (e *MessageEndWorkout) String() string {
	return fmt.Sprintf("EndWorkout[Seconds=%d,Distance=%d,Calories=%d,Speed=%d,"+
		"HeartRate=%d,Incline=%d]",
		e.Seconds, e.Distance, e.Calories, e.Speed, e.HeartRate, e.Incline)
}

type MessageProgramGraphics struct {
	Graph [18]byte
}

func (e *MessageProgramGraphics) MessageType() MessageType {
	return MessageTypeProgramGraphics
}

func (e *MessageProgramGraphics) MarshalBinary() ([]byte, error) {
	return []byte{
		byte(e.MessageType()),
		e.Graph[0], e.Graph[1], e.Graph[2], e.Graph[3], e.Graph[4],
		e.Graph[5], e.Graph[6], e.Graph[7], e.Graph[8], e.Graph[9],
		e.Graph[10], e.Graph[11], e.Graph[12], e.Graph[13], e.Graph[14],
		e.Graph[15], e.Graph[16], e.Graph[17],
	}, nil
}

func (e *MessageProgramGraphics) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.Graph = [18]byte{
		data[1], data[2], data[3], data[4], data[5], data[6], data[7],
		data[8], data[9], data[10], data[11], data[12], data[13], data[14],
		data[15], data[16], data[17], data[18],
	}

	return nil
}

func (e *MessageProgramGraphics) ExpectedLength() int {
	return 19
}

func (e *MessageProgramGraphics) String() string {
	return fmt.Sprintf("ProgramGraphics[Graph=%d %d %d %d %d %d %d %d %d %d "+
		"%d %d %d %d %d %d %d %d]",
		e.Graph[0], e.Graph[1], e.Graph[2], e.Graph[3], e.Graph[4],
		e.Graph[5], e.Graph[6], e.Graph[7], e.Graph[8], e.Graph[9],
		e.Graph[10], e.Graph[11], e.Graph[12], e.Graph[13], e.Graph[14],
		e.Graph[15], e.Graph[16], e.Graph[17])
}

type MessageCommand struct {
	Command CommandType
}

func (e *MessageCommand) MessageType() MessageType {
	return MessageTypeCommand
}

func (e *MessageCommand) MarshalBinary() ([]byte, error) {
	return []byte{byte(e.MessageType()), byte(e.Command)}, nil
}

func (e *MessageCommand) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.Command = CommandType(data[1])

	return nil
}

func (e *MessageCommand) ExpectedLength() int {
	return 2
}

func (e *MessageCommand) String() string {
	return fmt.Sprintf("Command[Command=%s]", e.Command)
}

type MessageProgram struct {
	Program Program
}

func (e *MessageProgram) MessageType() MessageType {
	return MessageTypeProgram
}

func (e *MessageProgram) MarshalBinary() ([]byte, error) {
	return []byte{
		byte(e.MessageType()),
		byte(e.Program >> 8),
		byte(e.Program),
	}, nil
}

func (e *MessageProgram) UnmarshalBinary(data []byte) error {
	if err := messageSanityCheck(e, data); err != nil {
		return err
	}

	e.Program = Program(uint16(data[1])<<8 + uint16(data[2]))

	return nil
}

func (e *MessageProgram) ExpectedLength() int {
	return 3
}

func (e *MessageProgram) String() string {
	return fmt.Sprintf("Program[Program=%s]", e.Program)
}
