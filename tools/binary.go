package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/c3-kotatsuneko/protobuf/gen/game/resources"
	"github.com/c3-kotatsuneko/protobuf/gen/game/rpc"
	"google.golang.org/protobuf/proto"
)

// Protobufメッセージをエンコードして16進数文字列に変換する関数
func encodeRequestToHex() string {
	player := &resources.Player{
		PlayerId: "2",
		Name:     "admin1",
		Color:    "red",
		Score:    10,
		Rank:     5,
		Time:     1,
	}

	request := &rpc.GameStatusRequest{
		RoomId:  "hoge",
		Event:   resources.Event_EVENT_ENTER_ROOM,
		Mode:    resources.Mode_MODE_TIME_ATTACK,
		Players: player,
	}

	data, err := proto.Marshal(request)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	hexData := hex.EncodeToString(data)
	return hexData
}

// Protobufメッセージをエンコードして16進数文字列に変換する関数
func encodeResponseToHex() string {
	player := &resources.Player{
		PlayerId: "1",
		Name:     "admin",
		Color:    "red",
		Score:    10,
		Rank:     5,
		Time:     1,
	}

	response := &rpc.GameStatusResponse{
		RoomId:  "hoge",
		Event:   resources.Event_EVENT_ENTER_ROOM,
		Players: []*resources.Player{player},
		Time:    -1,
	}

	data, err := proto.Marshal(response)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	hexData := hex.EncodeToString(data)
	return hexData
}

// 16進数文字列をデコードしてProtobufメッセージに変換する関数
func decodeRequestFromHex() *rpc.GameStatusRequest {
	hexData := "0a04686f67651001180122150a0131120561646d696e1a03726564200a28053001"
	binaryData, err := hex.DecodeString(hexData)
	if err != nil {
		log.Fatal("decoding hex error: ", err)
	}

	var decodedRequest rpc.GameStatusRequest
	if err := proto.Unmarshal(binaryData, &decodedRequest); err != nil {
		log.Fatal("unmarshaling error: ", err)
	}

	return &decodedRequest
}

// 16進数文字列をデコードしてProtobufメッセージに変換する関数
func decodeResponseFromHex() *rpc.GameStatusResponse {
	hexData := "0A04686F676510011A150A0131120561646D696E1A03726564200A2805300120FFFFFFFFFFFFFFFFFF01"
	binaryData, err := hex.DecodeString(hexData)
	if err != nil {
		log.Fatal("decoding hex error: ", err)
	}

	var decodedResponse rpc.GameStatusResponse
	if err := proto.Unmarshal(binaryData, &decodedResponse); err != nil {
		log.Fatal("unmarshaling error: ", err)
	}

	return &decodedResponse
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("please provide a command: encode request/response or decode request/response")
	}

	command := os.Args[1]
	subCommand := os.Args[2]

	switch command {
	case "encode":
		switch subCommand {
		case "request":
			hexData := encodeRequestToHex()
			fmt.Printf("%+v\n", hexData)
		case "response":
			hexData := encodeResponseToHex()
			fmt.Printf("%+v\n", hexData)
		default:
			log.Fatal("unknown subcommand for encode: ", subCommand)
		}
	case "decode":
		switch subCommand {
		case "request":
			decodedRequest := decodeRequestFromHex()
			fmt.Printf("%+v\n", decodedRequest)
		case "response":
			decodedResponse := decodeResponseFromHex()
			fmt.Printf("%+v\n", decodedResponse)
		default:
			log.Fatal("unknown subcommand for decode: ", subCommand)
		}
	default:
		log.Fatal("unknown command: ", command)
	}
}
