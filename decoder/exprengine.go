package decoder

import (
	"strconv"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/negbie/logp"
	"github.com/sipcapture/heplify-server/sipparser"
)

// ExprEngine struct
type ExprEngine struct {
	hepPkt *HEP
	prog   []*vm.Program
	env    map[string]interface{}
	v      vm.VM
}

func (e *ExprEngine) GetHEPStruct() *HEP { return e.hepPkt }

func (e *ExprEngine) GetSIPStruct() *sipparser.SipMsg { return e.hepPkt.SIP }

func (e *ExprEngine) GetHEPProtoType() uint32 { return e.hepPkt.GetProtoType() }

func (e *ExprEngine) GetHEPSrcIP() string { return e.hepPkt.GetSrcIP() }

func (e *ExprEngine) GetHEPSrcPort() uint32 { return e.hepPkt.GetSrcPort() }

func (e *ExprEngine) GetHEPDstIP() string { return e.hepPkt.GetDstIP() }

func (e *ExprEngine) GetHEPDstPort() uint32 { return e.hepPkt.GetDstPort() }

func (e *ExprEngine) GetHEPTimeSeconds() uint32 { return e.hepPkt.GetTsec() }

func (e *ExprEngine) GetHEPTimeUseconds() uint32 { return e.hepPkt.GetTmsec() }

func (e *ExprEngine) GetHEPNodeID() uint32 { return e.hepPkt.GetNodeID() }

func (e *ExprEngine) GetRawMessage() string { return e.hepPkt.GetPayload() }

func (e *ExprEngine) SetRawMessage(p string) uint8 {
	e.hepPkt.Payload = p
	return 1
}

func (e *ExprEngine) SetCustomSIPHeader(m map[string]string) uint8 {

	if e.hepPkt.SIP.CustomHeader == nil {
		e.hepPkt.SIP.CustomHeader = make(map[string]string)
	}

	for k, v := range m {
		e.hepPkt.SIP.CustomHeader[k] = v
	}
	return 1
}

func (e *ExprEngine) SetHEPField(field string, value string) uint8 {

	switch field {
	case "ProtoType":
		if i, err := strconv.Atoi(value); err == nil {
			e.hepPkt.ProtoType = uint32(i)
		}
	case "SrcIP":
		e.hepPkt.SrcIP = value
	case "SrcPort":
		if i, err := strconv.Atoi(value); err == nil {
			e.hepPkt.SrcPort = uint32(i)
		}
	case "DstIP":
		e.hepPkt.DstIP = value
	case "DstPort":
		if i, err := strconv.Atoi(value); err == nil {
			e.hepPkt.DstPort = uint32(i)
		}
	case "NodeID":
		if i, err := strconv.Atoi(value); err == nil {
			e.hepPkt.NodeID = uint32(i)
		}
	case "CID":
		e.hepPkt.CID = value
	case "SID":
		e.hepPkt.SID = value
	case "NodeName":
		e.hepPkt.NodeName = value
	case "TargetName":
		e.hepPkt.TargetName = value
	}
	return 1
}

func (e *ExprEngine) SetSIPHeader(header string, value string) uint8 {

	switch header {
	case "ViaOne":
		e.hepPkt.SIP.ViaOne = value
	case "FromUser":
		e.hepPkt.SIP.FromUser = value
	case "FromHost":
		e.hepPkt.SIP.FromHost = value
	case "FromTag":
		e.hepPkt.SIP.FromTag = value
	case "ToUser":
		e.hepPkt.SIP.ToUser = value
	case "ToHost":
		e.hepPkt.SIP.ToHost = value
	case "ToTag":
		e.hepPkt.SIP.ToTag = value
	case "CallID":
		e.hepPkt.SIP.CallID = value
	case "XCallID":
		e.hepPkt.SIP.XCallID = value
	case "ContactUser":
		e.hepPkt.SIP.ContactUser = value
	case "ContactHost":
		e.hepPkt.SIP.ContactHost = value
	case "UserAgent":
		e.hepPkt.SIP.UserAgent = value
	case "Server":
		e.hepPkt.SIP.Server = value
	case "Authorization.Username":
		e.hepPkt.SIP.Authorization.Username = value
	case "PaiUser":
		e.hepPkt.SIP.PaiUser = value
	case "PaiHost":
		e.hepPkt.SIP.PaiHost = value
	}
	return 1
}

// Close implements interface
func (e *ExprEngine) Close() {}

// NewExprEngine returns the script engine struct
func NewExprEngine() (*ExprEngine, error) {
	logp.Debug("script", "register expr engine")

	e := &ExprEngine{}
	e.env = map[string]interface{}{
		"GetHEPStruct":       e.GetHEPStruct,
		"GetSIPStruct":       e.GetSIPStruct,
		"GetHEPProtoType":    e.GetHEPProtoType,
		"GetHEPSrcIP":        e.GetHEPSrcIP,
		"GetHEPSrcPort":      e.GetHEPSrcPort,
		"GetHEPDstIP":        e.GetHEPDstIP,
		"GetHEPDstPort":      e.GetHEPDstPort,
		"GetHEPTimeSeconds":  e.GetHEPTimeSeconds,
		"GetHEPTimeUseconds": e.GetHEPTimeUseconds,
		"GetHEPNodeID":       e.GetHEPNodeID,
		"GetRawMessage":      e.GetRawMessage,
		"SetRawMessage":      e.SetRawMessage,
		"SetCustomSIPHeader": e.SetCustomSIPHeader,
		"SetHEPField":        e.SetHEPField,
		"SetSIPHeader":       e.SetSIPHeader,
		"HashTable":          HashTable,
		"HashString":         HashString,
	}

	files, _, err := scanCode()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		prog, err := expr.Compile(file, expr.Env(e.env))
		if err != nil {
			return nil, err
		}
		e.prog = append(e.prog, prog)
	}

	e.v = vm.VM{}

	return e, nil
}

// Run will execute the script
func (e *ExprEngine) Run(hep *HEP) error {
	e.hepPkt = hep

	for _, prog := range e.prog {
		_, err := e.v.Run(prog, e.env)
		if err != nil {
			return err
		}
	}

	return nil
}