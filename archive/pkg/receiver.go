package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	events "github.com/sgnl-ai/caep.dev-receiver/pkg/ssf_events"

	"github.com/golang-jwt/jwt/v5"
)

const TransmitterConfigMetadataPath = "/.well-known/ssf-configuration"
const TransmitterPollRFC = "urn:ietf:rfc:8936"

// Initializes the SSF Receiver based on the specified configuration.
//
// Returns an error if any process of configuring the receiver, registering
// it with the transmitter, or setting up the poll interval failed
func ConfigureSsfReceiver(cfg ReceiverConfig) (SsfReceiver, error) {
	if cfg.TransmitterUrl == "" || cfg.TransmitterPollUrl == "" || len(cfg.EventsRequested) == 0 || cfg.AuthorizationToken == "" {
		return nil, errors.New("Receiver Config - missing required field")
	}

	transmitterUrl, err := url.Parse(cfg.TransmitterUrl)
	if err != nil {
		return nil, err
	}

	baseUrl := transmitterUrl.Host
	trailingPath := transmitterUrl.Path

	transmitterConfigEndpoint := "https://" + baseUrl + TransmitterConfigMetadataPath
	if trailingPath != "/" {
		transmitterConfigEndpoint += trailingPath
	}

	transmitterCfg, err := makeTransmitterConfigRequest(transmitterConfigEndpoint)
	if err != nil {
		return nil, err
	}

	if transmitterCfg.ConfigurationEndpoint == "" {
		return nil, errors.New("Given transmitter doesn't specify the configuration endpoint")
	}

	streamId, err := makeCreateStreamRequest(transmitterCfg.ConfigurationEndpoint, cfg)
	if err != nil {
		return nil, err
	}

	receiver := SsfReceiverImplementation{
		transmitterUrl:       cfg.TransmitterUrl,
		transmitterPollUrl:   cfg.TransmitterPollUrl,
		eventsRequested:      events.EventTypeArrayToEventUriArray(cfg.EventsRequested),
		authorizationToken:   cfg.AuthorizationToken,
		transmitterStatusUrl: transmitterCfg.StatusEndpoint,
		pollInterval:         300,
		streamId:             streamId,
		configurationUrl:     transmitterCfg.ConfigurationEndpoint,
	}
	if cfg.PollInterval != 0 {
		receiver.pollInterval = cfg.PollInterval
	}

	if cfg.PollCallback != nil {
		receiver.pollCallback = cfg.PollCallback
		receiver.InitPollInterval()
	}

	return &receiver, nil
}

// Makes the Transmitter Configuration Metadata request to determine
// the transmitter's configuration url for creating a stream
func makeTransmitterConfigRequest(url string) (*TransmitterConfig, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var configMetadata TransmitterConfig
	err = json.Unmarshal(body, &configMetadata)
	if err != nil {
		return nil, err
	}

	return &configMetadata, nil
}

// Makes the Create Stream Request to the transmitter
func makeCreateStreamRequest(url string, cfg ReceiverConfig) (string, error) {
	client := &http.Client{}

	delivery := SsfDelivery{Method: TransmitterPollRFC}
	createStreamRequest := CreateStreamReq{
		Delivery:        delivery,
		EventsRequested: events.EventTypeArrayToEventUriArray(cfg.EventsRequested),
	}

	requestBody, err := json.Marshal(createStreamRequest)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.AuthorizationToken)
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	type Stream struct {
		StreamId string `json:"stream_id"`
	}

	var stream Stream
	err = json.Unmarshal(body, &stream)
	if err != nil {
		return "", err
	}

	return stream.StreamId, nil
}

// Initializes the poll interval for the receiver that will intermittently
// send SSF Events to the specified callback function
func (receiver *SsfReceiverImplementation) InitPollInterval() {
	// Create a channel to listen for quit signals
	receiver.terminate = make(chan bool)

	// Start a Goroutine to run the request on a schedule
	go func() {
		for {
			select {
			case <-receiver.terminate:
				return
			default:
				println("Polling for Events")
				events, err := receiver.PollEvents()
				if err == nil {
					receiver.pollCallback(events)
				} else {
					// TODO: What to do on error?
					panic(err)
				}
				time.Sleep(time.Duration(receiver.pollInterval) * time.Second)
			}
		}
	}()
}

// TODO: Not Yet Implemented
func (receiver *SsfReceiverImplementation) ConfigureCallback(callback func(events []events.SsfEvent), pollInterval int) error {
	return nil
}

// Polls the transmitter for all available SSF Events, returning them as a list
// for use
func (receiver *SsfReceiverImplementation) PollEvents() ([]events.SsfEvent, error) {
	client := &http.Client{}
	pollRequest := PollTransmitterRequest{Acknowledgements: []string{}, MaxEvents: 10, ReturnImmediately: true}
	requestBody, err := json.Marshal(pollRequest)
	if err != nil {
		return []events.SsfEvent{}, err
	}

	req, err := http.NewRequest("POST", receiver.transmitterPollUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return []events.SsfEvent{}, err
	}

	req.Header.Set("Authorization", "Bearer "+receiver.authorizationToken)
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return []events.SsfEvent{}, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []events.SsfEvent{}, err
	}

	if response.StatusCode != 200 && response.StatusCode != 202 {
		return []events.SsfEvent{}, err
	}

	type SsfEventSets struct {
		Sets map[string]string `json:"sets"`
	}

	var ssfEventsSets SsfEventSets
	err = json.Unmarshal(body, &ssfEventsSets)
	if err != nil {
		return []events.SsfEvent{}, nil
	}

	if len(ssfEventsSets.Sets) > 0 {
		err = acknowledgeEvents(&ssfEventsSets.Sets, receiver)
		if err != nil {
			return []events.SsfEvent{}, nil
		}
	}
	events, err := parseSsfEventSets(&ssfEventsSets.Sets)
	return events, err
}

// Cleans up the resources used by the Receiver and deletes the Receiver's
// stream from the transmitter
func (receiver *SsfReceiverImplementation) DeleteReceiver() {
	receiver.terminate <- true

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", receiver.configurationUrl+"?stream_id="+receiver.streamId, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", receiver.authorizationToken)

	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
}

func (receiver *SsfReceiverImplementation) EnableStream() (StreamStatus, error) {
	if receiver.transmitterStatusUrl == "" {
		return 0, errors.New("configured receiver does not have transmitter stream url")
	}
	return receiver.sendStatusUpdateRequest(StreamEnabled)
}

func (receiver *SsfReceiverImplementation) PauseStream() (StreamStatus, error) {
	if receiver.transmitterStatusUrl == "" {
		return 0, errors.New("configured receiver does not have transmitter stream url")
	}
	return receiver.sendStatusUpdateRequest(StreamPaused)
}

func (receiver *SsfReceiverImplementation) DisableStream() (StreamStatus, error) {
	if receiver.transmitterStatusUrl == "" {
		return 0, errors.New("configured receiver does not have transmitter stream url")
	}
	return receiver.sendStatusUpdateRequest(StreamDisabled)
}

func (receiver *SsfReceiverImplementation) sendStatusUpdateRequest(streamStatus StreamStatus) (StreamStatus, error) {
	client := &http.Client{}
	updateStreamRequest := UpdateStreamRequest{StreamId: receiver.streamId, Status: EnumToStringStatusMap[streamStatus]}
	requestBody, err := json.Marshal(updateStreamRequest)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", receiver.transmitterStatusUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "Bearer "+receiver.authorizationToken)
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}
	type StatusResponse struct {
		Status string `json:"status"`
		Reason string `json:"reason,omitempty"`
	}

	var statusResponse StatusResponse
	err = json.Unmarshal(body, &statusResponse)
	if err != nil {
		return 0, err
	}
	return StatusEnumMap[statusResponse.Status], nil
}

func (receiver *SsfReceiverImplementation) GetStreamStatus() (StreamStatus, error) {
	if receiver.transmitterStatusUrl == "" {
		return 0, errors.New("transmitter does not support stream status")
	}

	client := &http.Client{}
	streamUrl := fmt.Sprintf("%s?stream_id=%s", receiver.transmitterStatusUrl, receiver.streamId)
	req, err := http.NewRequest("GET", streamUrl, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "Bearer "+receiver.authorizationToken)
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}
	type StatusResponse struct {
		Status string `json:"status"`
	}

	var statusResponse StatusResponse
	err = json.Unmarshal(body, &statusResponse)
	if err != nil {
		return 0, nil
	}

	return StatusEnumMap[statusResponse.Status], nil
}

// Method to acknowledge a list of JTI's (unique ids for each SSF Event) with the
// transmitter so the events are re-transmitted
func acknowledgeEvents(sets *map[string]string, receiver *SsfReceiverImplementation) error {
	ackList := make([]string, len(*sets))
	i := 0
	for jti := range *sets {
		ackList[i] = jti
		i++
	}

	client := &http.Client{}
	pollRequest := PollTransmitterRequest{Acknowledgements: ackList, MaxEvents: 0, ReturnImmediately: true}
	requestBody, err := json.Marshal(pollRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", receiver.transmitterPollUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+receiver.authorizationToken)
	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// Parses a list of JTI:JWT pairings, return a list of the SSF Events from the JWT's
func parseSsfEventSets(sets *map[string]string) ([]events.SsfEvent, error) {
	var ssfEventsList []events.SsfEvent

	for _, set := range *sets {
		token, err := jwt.Parse(set, func(token *jwt.Token) (interface{}, error) { return jwt.UnsafeAllowNoneSignatureType, nil })
		if err != nil {
			return []events.SsfEvent{}, err
		}

		token.Claims.GetSubject()
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return []events.SsfEvent{}, errors.New("Can't get JWT Claims")
		}

		ssfEvents := claims["events"].(map[string]interface{})
		for eventType, eventSubject := range ssfEvents {
			ssfEvent, err := events.EventStructFromEvent(eventType, eventSubject, claims)
			if err != nil {
				return []events.SsfEvent{}, err
			}

			ssfEventsList = append(ssfEventsList, ssfEvent)
		}
	}

	return ssfEventsList, nil
}
