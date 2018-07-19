package main

import (
	"encoding/base64"
	"fmt"
	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/ice"
)

func main() {
	//reader := bufio.NewReader(os.Stdin)
	//rawSd, err := reader.ReadString('\n')
	//if err != nil {
	//	panic(err)
	//}

	rawSd := "dj0wCm89cGlvbi13ZWJydGMgNTU5NjM5MDY2NDI5Mjk0NDgyMiAyIElOIElQNCAwLjAuMC4wCnM9LQp0PTAgMAphPWdyb3VwOkJVTkRMRSBkYXRhCmE9bXNpZC1zZW1hbnRpYzogV01TCm09YXBwbGljYXRpb24gOSBEVExTL1NDVFAgNTAwMApjPUlOIElQNCAxMjcuMC4wLjEKYT1zZXR1cDphY3RpdmUKYT1taWQ6ZGF0YQphPWljZS11ZnJhZzpETFRiZFlVUEtqb1hwQWxxCmE9aWNlLXB3ZDpKV2tvRXJIZmxSUEdzS0VJcnFkR2Z1d0V0ZmZIb3FSVgphPWljZS1saXRlCmE9ZmluZ2VycHJpbnQ6c2hhLTI1NiBDNTo5Rjo4MDowOTo0OTozMzpCNDo1RDo2MDpCMjpDMDpBQzoyQTo2MDowMDpFQToxODoyRTo1RDoxODpGNTo3OToxRTo4QzpENzo2ODo4NDpBOTowMTo5Rjo1RDpCMQphPXNjdHBtYXA6NTAwMCB3ZWJydGMtZGF0YWNoYW5uZWwgMTAyNAphPWNhbmRpZGF0ZTp1ZHBjYW5kaWRhdGUgMSB1ZHAgNjI4NzcgMTkyLjE2OC4xLjcwIDU4MjU3IHR5cCBob3N0CmE9ZW5kLW9mLWNhbmRpZGF0ZXMK"

	fmt.Println("")
	sd, err := base64.StdEncoding.DecodeString(rawSd)
	if err != nil {
		panic(err)
	}

	/* Everything below is the pion-WebRTC API, thanks for using it! */

	// Create a new RTCPeerConnection
	peerConnection := &webrtc.RTCPeerConnection{}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange = func(connectionState ice.ConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	}

	peerConnection.Ondatachannel = func(d *webrtc.RTCDataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label, d.ID)
		d.Onmessage = func(message []byte) {
			fmt.Printf("Message from DataChannel %s '%s'\n", d.Label, string(message))
		}
	}

	// Set the remote SessionDescription
	if err := peerConnection.SetRemoteDescription(string(sd)); err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	if err := peerConnection.CreateAnswer(); err != nil {
		panic(err)
	}

	// Get the LocalDescription and take it to base64 so we can paste in browser
	localDescriptionStr := peerConnection.LocalDescription.Marshal()
	fmt.Println(base64.StdEncoding.EncodeToString([]byte(localDescriptionStr)))
	select {}
}
