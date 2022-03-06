package treadonme

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

var (
	serviceUUID = ble.MustParse("49535343-FE7D-4AE5-8FA9-9FAFD205E455")
	writeUUID   = ble.MustParse("49535343-8841-43F4-A8D4-ECBE34729BB3")
	notifyUUID  = ble.MustParse("49535343-1E4D-4BD9-BA61-23C647249616")
)

var (
	ErrMissingCharacteristic = fmt.Errorf("missing expected characteristic")
	ErrMissingDescriptor     = fmt.Errorf("missing expected descriptor")
	ErrMissingService        = fmt.Errorf("missing expected service")
	ErrInvalidReadPayload    = fmt.Errorf("unexpected data from device")
	ErrAckTimeout            = fmt.Errorf("failed to get acknowledgement from device")
)

type MessageListener func(Message, error)

type Treadmill struct {
	addr       ble.Addr
	client     ble.Client
	bleDevice  ble.Device
	notifyChr  *ble.Characteristic
	writeChr   *ble.Characteristic
	writeMutex sync.Mutex

	listeners []MessageListener

	waitMutex sync.Mutex
	waitMap   map[MessageType][]chan interface{}
}

func New(addr string) (*Treadmill, error) {
	t := &Treadmill{
		addr:    ble.NewAddr(addr),
		waitMap: make(map[MessageType][]chan interface{}),
	}

	t.AddListener(t.waitForResponseListener)

	return t, nil
}

func (t *Treadmill) Connect(ctx context.Context) error {
	if t.bleDevice == nil {
		bleDevice, err := linux.NewDevice()
		if err != nil {
			return fmt.Errorf("problem creating new linux ble device handle: %w", err)
		}

		ble.SetDefaultDevice(bleDevice)
		t.bleDevice = bleDevice
	}

	dev, err := ble.Connect(ctx, func(a ble.Advertisement) bool {
		return a.Addr().String() == t.addr.String()
	})
	if err != nil {
		return fmt.Errorf("problem connecting to treadmill: %w", err)
	}

	svcs, err := dev.DiscoverServices([]ble.UUID{serviceUUID})
	if err != nil {
		return fmt.Errorf("failed to discover services on treadmill: %w", err)
	} else if len(svcs) == 0 {
		return fmt.Errorf("%w: %s", ErrMissingService, serviceUUID.String())
	}

	chrs, err := dev.DiscoverCharacteristics([]ble.UUID{writeUUID, notifyUUID}, svcs[0])
	if err != nil {
		return fmt.Errorf("failed to discover characteristics on treadmill: %w", err)
	} else if len(chrs) != 2 {
		return fmt.Errorf("%w: expected 2, got: %d", ErrMissingCharacteristic, len(chrs))
	}

	if chrs[0].UUID.Equal(writeUUID) {
		t.writeChr, t.notifyChr = chrs[0], chrs[1]
	} else {
		t.writeChr, t.notifyChr = chrs[1], chrs[0]
	}

	desc, err := dev.DiscoverDescriptors(nil, t.notifyChr)
	if err != nil {
		return err
	} else if len(desc) == 0 {
		return fmt.Errorf("%w: %d", ErrMissingDescriptor, len(desc))
	}

	if err := dev.Subscribe(t.notifyChr, false, t.recv); err != nil {
		return fmt.Errorf("failed to subscribe to notify characteristic: %w", err)
	}

	t.client = dev

	return nil
}

func (t *Treadmill) Close() error {
	if t.client != nil {
		if err := t.client.ClearSubscriptions(); err != nil {
			return fmt.Errorf("failed to clear treadmill client subscriptions: %w", err)
		}

		if err := t.client.CancelConnection(); err != nil {
			return err
		}

		t.client = nil
		t.notifyChr = nil
		t.writeChr = nil
	}

	if t.bleDevice != nil {
		if err := t.bleDevice.Stop(); err != nil {
			return err
		}

		t.bleDevice = nil
	}

	return nil
}

func (t *Treadmill) GetDeviceInfo() (*MessageDeviceInfo, error) {
	msg, err := t.writeWithResponse(&MessageDeviceInfo{}, MessageTypeDeviceInfo)
	if err != nil {
		return nil, err
	}

	return msg.(*MessageDeviceInfo), nil
}

func (t *Treadmill) SetUserProfile(sex SexType, age byte, weight Weight, height Height) error {
	_, err := t.writeWithResponse(&MessageUserProfile{Sex: sex, Age: age, Weight: weight, Height: height}, MessageTypeACK)
	if err != nil {
		return err
	}

	return nil
}

func (t *Treadmill) SetWorkoutTime(tm time.Duration) error {
	_, err := t.writeWithResponse(&MessageWorkoutTarget{Time: byte(tm.Minutes())}, MessageTypeACK)
	if err != nil {
		return err
	}

	return nil
}

func (t *Treadmill) SetMaxIncline(maxIncline byte) error {
	_, err := t.writeWithResponse(&MessageMaxIncline{MaxIncline: maxIncline}, MessageTypeACK)
	if err != nil {
		return err
	}

	return nil
}

func (t *Treadmill) SetWorkoutMode(mode WorkoutMode) error {
	_, err := t.writeWithResponse(&MessageSetWorkoutMode{Mode: mode}, MessageTypeSetWorkoutMode)
	if err != nil {
		return err
	}

	return nil
}

func (t *Treadmill) SetProgram(program Program) error {
	_, err := t.writeWithResponse(&MessageProgram{Program: program}, MessageTypeACK)
	if err != nil {
		return err
	}

	return nil
}

func (t *Treadmill) LevelUp() error {
	_, err := t.writeWithResponse(&MessageCommand{Command: CommandTypeLevelUp}, MessageTypeACK)
	if err != nil {
		return err
	}

	return nil
}

func (t *Treadmill) Start() error {
	profile := &MessageUserProfile{Sex: SexTypeMale, Age: 30, Weight: 155, Height: 72}
	if _, err := t.writeWithResponse(profile, MessageTypeACK); err != nil {
		return err
	}

	if _, err := t.writeWithResponse(&MessageProgram{ProgramManual}, MessageTypeACK); err != nil {
		return err
	}

	if _, err := t.writeWithResponse(&MessageWorkoutTarget{}, MessageTypeACK); err != nil {
		return err
	}

	if _, err := t.writeWithResponse(&MessageSetWorkoutMode{WorkoutModeStart}, MessageTypeSetWorkoutMode); err != nil {
		return err
	}

	if err := t.Close(); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	var lastErr error

	for idx := 0; idx < 5; idx++ {
		if lastErr = t.Connect(context.Background()); lastErr == nil {
			log.Printf("done")

			_, err := t.GetDeviceInfo()

			return err
		}
	}

	return lastErr
}

func (t *Treadmill) WaitForResponse(ctx context.Context, msgType MessageType) (Message, error) {
	t.waitMutex.Lock()

	waitChan := make(chan interface{})
	defer close(waitChan)

	t.waitMap[msgType] = append(t.waitMap[msgType], waitChan)
	t.waitMutex.Unlock()

	select {
	case msg := <-waitChan:
		switch v := msg.(type) {
		case Message:
			return v, nil
		case error:
			return nil, v
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// We shouldn't usually get here.
	return nil, nil
}

func (t *Treadmill) waitForResponseListener(msg Message, err error) {
	t.waitMutex.Lock()
	defer t.waitMutex.Unlock()

	if err != nil {
		// If there's an error notify everyone and bail.
		for _, chans := range t.waitMap {
			for _, c := range chans {
				c <- err
			}
		}

		t.waitMap = make(map[MessageType][]chan interface{})
	} else {
		for _, notify := range t.waitMap[msg.MessageType()] {
			notify <- msg
		}

		t.waitMap[msg.MessageType()] = make([]chan interface{}, 0, 5)
	}
}

func (t *Treadmill) AddListener(listener MessageListener) {
	t.listeners = append(t.listeners, listener)
}

func (t *Treadmill) writeWithResponse(msg Message, expect MessageType) (Message, error) {
	var response Message
	var responseError error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		response, responseError = t.WaitForResponse(ctx, expect)
	}()

	for idx := 0; idx < 10; idx++ {
		if err := t.write(msg); err != nil {
			return nil, err
		}

		// We need to wait at least 300ms before retrying.
		time.Sleep(300 * time.Millisecond)

		if response != nil || responseError != nil {
			return response, responseError
		}
	}

	return nil, fmt.Errorf("%w: waiting on %s from command %s", ErrAckTimeout, expect, msg)
}

func (t *Treadmill) write(msg Message) error {
	data, err := EncodeMessage(msg)
	if err != nil {
		return err
	}

	t.writeMutex.Lock()
	defer t.writeMutex.Unlock()

	log.Printf("C->T: %s -- %s", hex.EncodeToString(data), msg.String())

	if err := t.client.WriteCharacteristic(t.writeChr, data, true); err != nil {
		return err
	}

	return nil
}

func (t *Treadmill) recv(data []byte) {
	msg, err := ParseMessage(data)
	if err != nil {
		for _, l := range t.listeners {
			l(nil, err)
		}

		return
	}

	log.Printf("T->C: %s -- %s", hex.EncodeToString(data), msg.String())

	switch msg.MessageType() {
	case MessageTypeACK:
		// Don't ACK the ACK's, that'd be bad.
	case MessageTypeSetWorkoutMode:
		// Special, don't do anything.
	case MessageTypeDeviceInfo:
		// Also don't do anything.
	case MessageTypeWorkoutMode:
		// These are special, we just need to echo what we heard.
		t.ackWorkoutMode(msg.(*MessageWorkoutMode))
	case MessageTypeWorkoutData:
		fallthrough
	case MessageTypeHeartRateType:
		fallthrough
	case MessageTypeErrorCode:
		fallthrough
	case MessageTypeSpeed:
		fallthrough
	case MessageTypeIncline:
		fallthrough
	case MessageTypeLevel:
		fallthrough
	case MessageTypeRPM:
		fallthrough
	case MessageTypeHeartRate:
		fallthrough
	case MessageTypeTargetHeartRate:
		fallthrough
	case MessageTypeMaxSpeed:
		fallthrough
	case MessageTypeMaxIncline:
		fallthrough
	case MessageTypeMaxLevel:
		fallthrough
	case MessageTypeEndWorkout:
		fallthrough
	case MessageTypeProgramGraphics:
		t.ackCommand(msg.MessageType())
	default:
		log.Printf("unhandled ack condition: %s", msg)
	}

	for _, l := range t.listeners {
		l(msg, nil)
	}
}

func (t *Treadmill) ackWorkoutMode(msg *MessageWorkoutMode) {
	go func() {
		if err := t.write(msg); err != nil {
			for _, l := range t.listeners {
				l(nil, err)
			}
		}
	}()
}

func (t *Treadmill) ackCommand(msg MessageType) {
	go func() {
		if err := t.write(&MessageACK{Acknowledged: msg}); err != nil {
			for _, l := range t.listeners {
				l(nil, err)
			}
		}
	}()
}
