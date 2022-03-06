package treadonme

type MessageType byte

const (
	MessageTypeACK             MessageType = 0x00
	MessageTypeSetWorkoutMode  MessageType = 0x02
	MessageTypeWorkoutMode     MessageType = 0x03
	MessageTypeWorkoutTarget   MessageType = 0x04
	MessageTypeWorkoutData     MessageType = 0x06
	MessageTypeUserProfile     MessageType = 0x07
	MessageTypeProgram         MessageType = 0x08
	MessageTypeHeartRateType   MessageType = 0x09
	MessageTypeErrorCode       MessageType = 0x10
	MessageTypeSpeed           MessageType = 0x11
	MessageTypeIncline         MessageType = 0x12
	MessageTypeLevel           MessageType = 0x13
	MessageTypeRPM             MessageType = 0x14
	MessageTypeHeartRate       MessageType = 0x15
	MessageTypeTargetHeartRate MessageType = 0x20
	MessageTypeMaxSpeed        MessageType = 0x21
	MessageTypeMaxIncline      MessageType = 0x22
	MessageTypeMaxLevel        MessageType = 0x23
	MessageTypeUserIncline     MessageType = 0x25
	MessageTypeUserLevel       MessageType = 0x27
	MessageTypeEndWorkout      MessageType = 0x32
	MessageTypeProgramGraphics MessageType = 0x40
	MessageTypeDeviceInfo      MessageType = 0xf0
	MessageTypeCommand         MessageType = 0xf1
	MessageTypeUnknown         MessageType = 0xff
)

func (mt MessageType) Create() Message {
	switch mt {
	case MessageTypeACK:
		return &MessageACK{}
	case MessageTypeSetWorkoutMode:
		return &MessageSetWorkoutMode{}
	case MessageTypeWorkoutMode:
		return &MessageWorkoutMode{}
	case MessageTypeWorkoutTarget:
		return &MessageWorkoutTarget{}
	case MessageTypeWorkoutData:
		return &MessageWorkoutData{}
	case MessageTypeUserProfile:
		return &MessageUserProfile{}
	case MessageTypeProgram:
		return &MessageProgram{}
	case MessageTypeHeartRateType:
		return &MessageHeartRateType{}
	case MessageTypeErrorCode:
		return &MessageErrorCode{}
	case MessageTypeSpeed:
		return &MessageSpeed{}
	case MessageTypeIncline:
		return &MessageIncline{}
	case MessageTypeLevel:
		return &MessageLevel{}
	case MessageTypeRPM:
		return &MessageRPM{}
	case MessageTypeHeartRate:
		return &MessageHeartRate{}
	case MessageTypeTargetHeartRate:
		return &MessageTargetHeartRate{}
	case MessageTypeMaxSpeed:
		return &MessageMaxSpeed{}
	case MessageTypeMaxLevel:
		return &MessageMaxLevel{}
	case MessageTypeUserIncline:
		return &MessageUserLevel{}
	case MessageTypeUserLevel:
		return &MessageUserLevel{}
	case MessageTypeEndWorkout:
		return &MessageEndWorkout{}
	case MessageTypeProgramGraphics:
		return &MessageProgramGraphics{}
	case MessageTypeMaxIncline:
		return &MessageMaxIncline{}
	case MessageTypeDeviceInfo:
		return &MessageDeviceInfo{}
	case MessageTypeCommand:
		return &MessageCommand{}
	case MessageTypeUnknown:
		fallthrough
	default:
		return nil
	}
}

func (mt MessageType) String() string {
	switch mt {
	case MessageTypeACK:
		return "ACK"
	case MessageTypeSetWorkoutMode:
		return "SetWorkoutMode"
	case MessageTypeWorkoutMode:
		return "WorkoutMode"
	case MessageTypeWorkoutTarget:
		return "WorkoutTarget"
	case MessageTypeWorkoutData:
		return "WorkoutData"
	case MessageTypeUserProfile:
		return "UserProfile"
	case MessageTypeProgram:
		return "Program"
	case MessageTypeHeartRateType:
		return "HeartRateType"
	case MessageTypeErrorCode:
		return "ErrorCode"
	case MessageTypeSpeed:
		return "Speed"
	case MessageTypeIncline:
		return "Incline"
	case MessageTypeLevel:
		return "Level"
	case MessageTypeRPM:
		return "RPM"
	case MessageTypeHeartRate:
		return "HeartRate"
	case MessageTypeTargetHeartRate:
		return "TargetHeartRate"
	case MessageTypeMaxSpeed:
		return "MaxSpeed"
	case MessageTypeMaxLevel:
		return "MaxLevel"
	case MessageTypeUserIncline:
		return "UserIncline"
	case MessageTypeUserLevel:
		return "UserLevel"
	case MessageTypeEndWorkout:
		return "EndWorkout"
	case MessageTypeProgramGraphics:
		return "ProgramGraphics"
	case MessageTypeMaxIncline:
		return "MaxIncline"
	case MessageTypeDeviceInfo:
		return "DeviceInfo"
	case MessageTypeCommand:
		return "Command"
	case MessageTypeUnknown:
		fallthrough
	default:
		return "Unknown"
	}
}

type DeviceModel byte

const (
	DeviceModelF80 DeviceModel = 146
)

func (dm DeviceModel) String() string {
	switch dm {
	case DeviceModelF80:
		return "F80"
	default:
		return "Unknown"
	}
}

type UnitsType byte

const (
	UnitsTypeMetric   UnitsType = 0x0
	UnitsTypeImperial UnitsType = 0x1
)

func (ut UnitsType) String() string {
	switch ut {
	case UnitsTypeMetric:
		return "Metric"
	case UnitsTypeImperial:
		return "Imperial"
	default:
		return "Unknown"
	}
}

type WorkoutMode byte

const (
	WorkoutModeIdle    WorkoutMode = 0x01
	WorkoutModeStart   WorkoutMode = 0x02
	WorkoutModeRunning WorkoutMode = 0x04
	WorkoutModePause   WorkoutMode = 0x06
	WorkoutModeDone    WorkoutMode = 0x07
)

func (wm WorkoutMode) String() string {
	switch wm {
	case WorkoutModeIdle:
		return "Idle"
	case WorkoutModeRunning:
		return "Running"
	case WorkoutModeStart:
		return "Start"
	case WorkoutModePause:
		return "Pause"
	case WorkoutModeDone:
		return "Done"
	default:
		return "Unknown"
	}
}

type SexType byte

const (
	SexTypeMale   SexType = 0x01
	SexTypeFemale SexType = 0x02
)

func (st SexType) String() string {
	switch st {
	case SexTypeMale:
		return "Male"
	case SexTypeFemale:
		return "Female"
	default:
		return "Unknown"
	}
}

const (
	startOfMessage byte = 0x5b
	endOfMessage   byte = 0x5d
)

type CommandType byte

const (
	CommandTypeStart     CommandType = 0x01
	CommandTypeLevelUp   CommandType = 0x02
	CommandTypeLevelDown CommandType = 0x03
	CommandTypeStop      CommandType = 0x06
)

func (ct CommandType) String() string {
	switch ct {
	case CommandTypeStart:
		return "Start"
	case CommandTypeLevelUp:
		return "LevelUp"
	case CommandTypeLevelDown:
		return "LevelDown"
	case CommandTypeStop:
		return "Stop"
	default:
		return "Unknown"
	}
}

type Speed byte

type Weight uint16

type Height byte

type Program uint16

const (
	ProgramManual   Program = 0x1001 //0x10 0x01
	ProgramHill     Program = 0x2002 //0x20 0x02
	ProgramFatBurn  Program = 0x2003 //0x20 0x03
	ProgramCardio   Program = 0x2004 //0x20 0x04
	ProgramStrength Program = 0x2005 //0x20 0x05
	ProgramInterval Program = 0x2006 //0x20 0x06
	ProgramHR1      Program = 0x3009 //0x30 0x09
	ProgramHR2      Program = 0x300a //0x30 0x0a
	ProgramUser1    Program = 0x3007 //0x30 0x07
	ProgramUser2    Program = 0x3008 //0x30 0x08
	ProgramFusion   Program = 0x600c //0x60 0x0c
)

func (p Program) String() string {
	switch p {
	case ProgramManual:
		return "Manual"
	case ProgramHill:
		return "Hill"
	case ProgramFatBurn:
		return "FatBurn"
	case ProgramCardio:
		return "Cardio"
	case ProgramStrength:
		return "Strength"
	case ProgramInterval:
		return "Interval"
	case ProgramHR1:
		return "HR1"
	case ProgramHR2:
		return "HR2"
	case ProgramUser1:
		return "User1"
	case ProgramUser2:
		return "User2"
	case ProgramFusion:
		return "Fusion"
	default:
		return "Unknown"
	}
}
